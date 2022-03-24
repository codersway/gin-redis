package lock

import (
	"fmt"
	"gin-redis/conf"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Locker struct {
	key    string
	expire time.Duration
	unlock bool // 为true则解锁
	script *redis.Script
}

const incrLua = `
if redis.call('get', KEYS[1]) == ARGV[1] then
  return redis.call('expire', KEYS[1],ARGV[2]) 				
 else
   return '0' 					
end`

// 默认30s过期时间
func NewLocker(key string) *Locker {
	return &Locker{key: key, expire: 30 * time.Second, script: redis.NewScript(incrLua)}
}

// 需要设置过期时间
func NewLockerWithTTL(key string, expire time.Duration) *Locker {
	if expire.Seconds() <= 0 {
		panic("error expire")
	}
	return &Locker{key: key, expire: expire, script: redis.NewScript(incrLua)}
}

func (l *Locker) Lock() *Locker {
	nx := conf.RedisClient.SetNX(conf.Ctx, l.key, "1", l.expire)
	if ok, err := nx.Result(); err != nil || !ok {
		panic(fmt.Sprintf("lock error with key:%s", l.key))
	}
	l.ExpandLockTime()
	return l
}

// 锁自动续期
func (l *Locker) ExpandLockTime() {
	sleepTime := l.expire.Seconds() * 2 / 3
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(sleepTime))
			if l.unlock {
				break
			}
			l.resetExpire()
		}
	}()
}

// 重新设置过期时间
func (l *Locker) resetExpire() {
	// cmd := conf.RedisClient.Expire(conf.Ctx, this.key, this.expire)
	//
	// fmt.Println(cmd)
	cmd := l.script.Run(conf.Ctx, conf.RedisClient, []string{l.key}, 1, l.expire.Seconds())
	v, err := cmd.Result()
	log.Printf("key=%s ,续期结果:%v,%v\n", l.key, err, v)
}

func (l *Locker) Unlock() {
	l.unlock = true
	conf.RedisClient.Del(conf.Ctx, l.key)
}
