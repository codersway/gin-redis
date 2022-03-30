package lua

import (
	"context"
	"gin-redis/conf"
	"log"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/go-redis/redis/v8"
)

type Lua struct {
	Conn *redis.Client
}

var (
	ctx             = context.Background()
	requestTimes    = 100
	beforeSeconds   = 60
	ThrottleStrKey  = "throttles_str:"
	ThrottleZsetKey = "throttles_zset:"
	conn            = conf.Conn()
	this            = Lua{
		Conn: conn,
	}
)

func TestLua1(t *testing.T) {
	sc := redis.NewScript(`
return {KEYS[1],KEYS[2],ARGV[1],ARGV[2]}
`)
	sha, err := sc.Load(ctx, this.Conn).Result()
	if err != nil {
		log.Fatalln(err)
	}

	left := this.Conn.EvalSha(ctx, sha, []string{"key1", "key2"}, "val1", "val2")

	t.Log(left)
	text, err := left.Text()
	if err != nil {
		return
	}
	assert.Equal(t, "key1 key2 val1 val2", text)
}

// todo
func TestLua12(t *testing.T) {
	sc := redis.NewScript(`

	local num = redis.call('SET', KEYS[1], ARGV[1]);  
	if not num then
		return 0;
	else
		local res = tonumber(num) * 3;
		redis.call('SET',KEYS[1], res); 
		return res;
	end
`)
	sha, err := sc.Load(ctx, this.Conn).Result()
	if err != nil {
		log.Fatalln(err)
	}

	left := this.Conn.EvalSha(ctx, sha, []string{"lua12"}, 33, 2)

	t.Log(left)
}

// 某个ip在短时间内频繁访问页面，需要记录并检测出来，就可以通过Lua脚本高效的实现
func TestLua3(t *testing.T) {
	sc := redis.NewScript(`
	local times = redis.call('incr',KEYS[1])
	
	if times == 1 then
		redis.call('expire',KEYS[1], ARGV[1])
	end
	
	if times > tonumber(ARGV[2]) then
		return 0
	end
	return 1
`)
	sha, err := sc.Load(ctx, this.Conn).Result()
	if err != nil {
		log.Fatalln(err)
	}

	left := this.Conn.EvalSha(ctx, sha, []string{"lua12"}, 33, 99)

	t.Log(left)
	t.Log(left.Result())
	res, _ := left.Result()
	assert.Equal(t, res.(int64), int64(1))
}

func TestLua4(t *testing.T) {
	sc := redis.NewScript(`
	local key=KEYS[1]
	local v1=ARGV[1]
	local v2=ARGV[2]
	redis.call("DEL", key)
	redis.call("LPUSH", key, v1, v2)
	local list=redis.call("LRANGE", key, 0, -1)
	return list
`)

	sha, err := sc.Load(ctx, this.Conn).Result()
	if err != nil {
		log.Fatalln(err)
	}
	left := this.Conn.EvalSha(ctx, sha, []string{"lua4"}, "jack", "peter")
	t.Log(left.Result())
	t.Log(left.Val())

	assert.Equal(t, []string{"jack", "peter"}, left)
}
