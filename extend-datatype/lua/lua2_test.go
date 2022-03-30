package lua

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/go-redis/redis/v8"
)

var wg = sync.WaitGroup{}

func TestLua2(t *testing.T) {
	var wg sync.WaitGroup

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       1,        // use default DB
	})
	for _, v := range []string{"5824742984", "5824742984", "5824742983", "5824742983", "5824742982", "5824742980"} {
		wg.Add(1)
		go evalScript(client, v, &wg)
	}
	wg.Wait()
}

func createScript() *redis.Script {
	script := redis.NewScript(`
		local goodsSurplus
		local flag
		local existUserIds    = tostring(KEYS[1])
		local memberUid       = tonumber(ARGV[1])
		local goodsSurplusKey = tostring(KEYS[2])
		local hasBuy = redis.call("sIsMember", existUserIds, memberUid)

		if hasBuy ~= 0 then
		  return 0
		end
		

		goodsSurplus =  redis.call("GET", goodsSurplusKey)
		if goodsSurplus == false then
		  return 0
		end
		
		-- 没有剩余可抢购物品
		goodsSurplus = tonumber(goodsSurplus)
		if goodsSurplus <= 0 then
		  return 0
		end
		
		flag = redis.call("SADD", existUserIds, memberUid)
		flag = redis.call("DECR", goodsSurplusKey)
		
		return 1
	`)
	return script
}

func evalScript(client *redis.Client, userId string, wg *sync.WaitGroup) {
	defer wg.Done()
	script := createScript()
	sha, err := script.Load(client.Context(), client).Result()
	if err != nil {
		log.Fatalln(err)
	}
	ret := client.EvalSha(client.Context(), sha, []string{
		"hadBuyUids",
		"goodsSurplus",
	}, userId)
	if result, err := ret.Result(); err != nil {
		log.Fatalf("Execute Redis fail: %v", err.Error())
	} else {
		fmt.Println("")
		fmt.Printf("userid: %s, result: %d", userId, result)
	}
}
