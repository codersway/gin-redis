package throttles

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/util/grand"
)

type Throttles struct {
	Conn *redis.Client
}

var (
	ctx                   = context.Background()
	requestTimes    int64 = 100
	ThrottleLuaKey        = "throttles_lua:"
	ThrottleStrKey        = "throttles_str:"
	ThrottleZsetKey       = "throttles_zset:"
)

// 基于setnx实现限流
// 设置有效期1s的计数器，通过decr实现1s放行1个请求
// 没有使用lua脚本实现原子操作
func (this *Throttles) ThrottlesStr() string {
	this.Conn.SetNX(ctx, ThrottleStrKey, 1, time.Second)
	left := this.Conn.Decr(ctx, ThrottleStrKey).Val()

	if left < 0 {
		return "error"
	} else {
		return "success"
	}
}

// 要求60s内最多访问100次
// todo 已失效member动态删除
func (this *Throttles) ThrottlesZset() string {
	now := time.Now().Unix()
	before := time.Now().Add(-time.Second * 60).Unix()

	random := grand.Digits(10)

	count := this.Conn.ZCount(ctx, ThrottleZsetKey, strconv.Itoa(int(before)), strconv.Itoa(int(now)))

	if count.Val() < requestTimes {
		z := redis.Z{Member: random, Score: float64(now)}
		add := this.Conn.ZAdd(ctx, ThrottleZsetKey, &z)

		if _, err := add.Result(); err != nil {
			return "error"
		}

		return "add success"
	}

	return "limited"
}

func (this *Throttles) ThrottlesLua() interface{} {
	script := redis.NewScript(`
	local max = tonumber(ARGV[1])
	local interval_milliseconds = tonumber(ARGV[2])
	local current = tonumber(redis.call('get', KEYS[1]) or 0)
	
	if (current + 1 > max) then
		return true
	else
		redis.call('incrby', KEYS[1], 1)
		if (current == 0) then
			redis.call('pexpire', KEYS[1], interval_milliseconds)
		end
		return false
	end
`)
	sha, err := script.Load(ctx, this.Conn).Result()
	if err != nil {
		log.Fatalln(err)
	}
	left := this.Conn.EvalSha(ctx, sha, []string{ThrottleLuaKey}, 1, 5)

	val := left.Val()
	return val
}
