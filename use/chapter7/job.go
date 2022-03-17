package chapter7

import (
	"context"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"log"
	"reflect"
	"strings"
	"time"
)

// 职位搜索引擎
type Client struct {
	Conn *redis.Client
}

var (
	ctx = context.Background()
)

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

func (this *Client) AddJob(jobId string, requireSkills []string) {
	this.Conn.SAdd(ctx, "job:"+jobId, requireSkills)
}

func (this *Client) IsQualified(jobId string, candidateSkills []string) []string {
	var res *redis.StringSliceCmd
	temp := uuid.NewV4().String()
	pipeline := this.Conn.TxPipeline()
	pipeline.SAdd(ctx, temp, candidateSkills)
	pipeline.Expire(ctx, temp, 5*time.Second)
	res = pipeline.SDiff(ctx, "job:"+jobId, temp)

	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in RecordClick: ", err)
		return nil
	}
	return res.Val()
}

func (this *Client) IndexJob(jobId string, skills []string) {
	countSkill := Set{}
	pipeline := this.Conn.TxPipeline()
	for _, skill := range skills {
		pipeline.SAdd(ctx, "idx:skill:"+skill, jobId)
		countSkill.Add(skill)
	}
	pipeline.ZAdd(ctx, "idx:jobs:req", &redis.Z{Member: jobId, Score: float64(len(countSkill))})
	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in RecordClick: ", err)
		return
	}
}

func (this *Client) FindJobs(candidateSkills []string) []string {
	skills := map[string]float64{}
	for _, skill := range candidateSkills {
		skills["skill:"+skill] = 1
	}
	jobScores := this.ZUnion(skills, 30, "")
	finalResult := this.Zintersect(map[string]float64{jobScores: -1, "jobs:req": 1}, 30, "")
	return this.Conn.ZRangeByScore(ctx, "idx:"+finalResult, &redis.ZRangeBy{Max: "0", Min: "0"}).Val()
}

// todo 封装zset的ZInterStore和ZUnionStore操作
func (this *Client) zsetCommon(method string, scores map[string]float64, ttl int, Aggre string) string {
	id := uuid.NewV4().String()
	pipeline := this.Conn.TxPipeline()

	zstore := redis.ZStore{}
	for key := range scores {
		zstore.Keys = append(zstore.Keys, "idx:"+key)
		zstore.Weights = append(zstore.Weights, scores[key])
	}
	switch strings.ToLower(Aggre) {
	case "max":
		zstore.Aggregate = "MAX"
	case "min":
		zstore.Aggregate = "MIN"
	default:
		zstore.Aggregate = "SUM"
	}
	methodValue := reflect.ValueOf(pipeline).MethodByName(method)
	args := []reflect.Value{reflect.ValueOf("idx:" + id), reflect.ValueOf(&zstore)}
	methodValue.Call(args)
	pipeline.Expire(ctx, "idx:"+id, time.Duration(ttl)*time.Second)
	if _, err := pipeline.Exec(ctx); err != nil {
		log.Println("pipeline err in zsetCommon: ", err)
		return ""
	}
	return id
}

func (this *Client) Zintersect(items map[string]float64, ttl int, Aggre string) string {
	return this.zsetCommon("ZInterStore", items, ttl, Aggre)
}

func (this *Client) ZUnion(items map[string]float64, ttl int, Aggre string) string {
	return this.zsetCommon("ZUnionStore", items, ttl, Aggre)
}
