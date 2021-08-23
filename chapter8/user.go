package chapter8

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"log"
	"math"
	"reflect"
	"strconv"
	"time"
)

type Client struct {
	Conn *redis.Client
}

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

var (
	ctx                    = context.Background()
	Postperpass      int64 = 1000
	Hometimelinesize int64 = 1000
)

func (this *Client) CreateUser(login, name string) string {
	lock := this.AcquireLockWithTimeout("user:"+login, 10, 10)
	defer this.ReleaseLock("user:"+login, lock)

	if lock == "" {
		return ""
	}
	if this.Conn.HGet(ctx, "users:", login).Val() != "" {
		return ""
	}
	id := this.Conn.Incr(ctx, "user:id:").Val()

	pipeline := this.Conn.TxPipeline()
	pipeline.HSet(ctx, "users:", login, id)
	// todo
	pipeline.HMSet(ctx, fmt.Sprintf("user:%s", strconv.Itoa(int(id))), "login", login, "id", id, "name", name, "followers", 0, "following", 0, "posts", 0, "signup", time.Now().UnixNano())

	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in CreateUser: ", err)
		return ""
	}
	return strconv.Itoa(int(id))
}

func (this *Client) CreateStatus(uid, message string, data map[string]interface{}) string {
	pipeline := this.Conn.TxPipeline()
	pipeline.HGet(ctx, fmt.Sprintf("user:%s", uid), "login")
	pipeline.Incr(ctx, "status:id:")
	res, err := pipeline.Exec(ctx)

	if err != nil {
		log.Println("pipeline err in CreateStatus: ", err)
		return ""
	}
	login, id := res[0].(*redis.StringCmd).Val(), res[1].(*redis.IntCmd).Val()
	if login == "" {
		return ""
	}
	data["message"] = message
	data["posted"] = time.Now().UnixNano()
	data["id"] = id
	data["uid"] = uid
	data["login"] = login

	pipeline.HMSet(ctx, fmt.Sprintf("status:%s", strconv.Itoa(int(id))), data)
	pipeline.HIncrBy(ctx, fmt.Sprintf("user:%s", uid), "posts", 1)
	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in CreateStatus: ", err)
		return ""
	}
	return strconv.Itoa(int(id))
}

func (this *Client) PostStatus(uid, message string, data map[string]interface{}) string {
	id := this.CreateStatus(uid, message, data)
	if id == "" {
		return ""
	}
	posted, err := this.Conn.HGet(ctx, fmt.Sprintf("status:%s", id), "posted").Float64()
	if err != nil {
		log.Printf("hget from status: %s, err: %s", id, err.Error())
		return ""
	}
	if posted == 0 {
		return ""
	}
	post := redis.Z{Member: id, Score: posted}
	this.Conn.ZAdd(ctx, fmt.Sprintf("profile:%s", uid), &post)

	this.SyndicateStatus(uid, post, 0)
	return id
}

func (this *Client) SyndicateStatus(uid string, post redis.Z, start int) {
	followers := this.Conn.ZRangeByScoreWithScores(ctx, fmt.Sprintf("followers:%s", uid), &redis.ZRangeBy{
		Min: "0", Max: "inf", Offset: int64(start), Count: Postperpass,
	}).Val()

	pipeline := this.Conn.TxPipeline()
	for i, z := range followers {
		follower := z.Member.(string)
		start = i + 1
		pipeline.ZAdd(ctx, fmt.Sprintf("home:%s", follower), &post)
		pipeline.ZRemRangeByRank(ctx, fmt.Sprintf("home:%s", follower), 0, -Hometimelinesize-1)
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in syndicatesStatus: ", err)
		return
	}

	if len(followers) >= int(Postperpass) {
		this.executeLater("default", "SyndicateStatus", uid, post, start)
	}
}

func (this *Client) executeLater(queue, name string, args ...interface{}) {
	go func() {
		methodValue := reflect.ValueOf(this).MethodByName(name)
		methodArgs := make([]reflect.Value, 0, len(args))
		for _, v := range args {
			value := reflect.ValueOf(v).Interface()
			methodArgs = append(methodArgs, reflect.ValueOf(value))
		}
		methodValue.Call(methodArgs)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (this *Client) AcquireLockWithTimeout(lockName string, acquireTimeout, lockTimeout float64) string {
	identifier := uuid.NewV4().String()
	lockName = "lock:" + lockName
	finalLockTimeout := math.Ceil(lockTimeout)

	end := time.Now().UnixNano() + int64(acquireTimeout*1e9)
	for time.Now().UnixNano() < end {
		if this.Conn.SetNX(ctx, lockName, identifier, 0).Val() {
			this.Conn.Expire(ctx, lockName, time.Duration(finalLockTimeout)*time.Second)
			return identifier
		} else if this.Conn.TTL(ctx, lockName).Val() < 0 {
			this.Conn.Expire(ctx, lockName, time.Duration(finalLockTimeout)*time.Second)
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ""
}

func (this *Client) ReleaseLock(lockName, identifier string) bool {
	lockName = "lock:" + lockName
	var flag = true
	for flag {
		err := this.Conn.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()
			if tx.Get(ctx, lockName).Val() == identifier {
				pipe.Del(ctx, lockName)
				if _, err := pipe.Exec(ctx); err != nil {
					return err
				}
				flag = true
				return nil
			}

			tx.Unwatch(ctx)
			flag = false
			return nil
		})

		if err != nil {
			log.Println("watch failed in ReleaseLock, err is: ", err)
			return false
		}
	}
	return true
}
