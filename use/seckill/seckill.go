package seckill

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"time"
)

type Client struct {
	Conn *redis.Client
}

var (
	ctx = context.Background()
)

func NewClient(conn *redis.Client) *Client {
	return &Client{
		Conn: conn,
	}
}

func (this *Client) SecKill(total int) bool {
	actualTotal, _ := this.Conn.Get(ctx, "sec-total").Int()

	if actualTotal < total {
		err := this.Conn.Watch(ctx, func(tx *redis.Tx) error {
			pipe := tx.TxPipeline()

			pipe.HSet(ctx, "sec-list", fmt.Sprintf("user_id:%d", rand.Int()), time.Now().String())
			pipe.Incr(ctx, "sec-total")
			if _, err := pipe.Exec(ctx); err != nil {
				return err
			}
			tx.Unwatch(ctx)
			return nil
		}, "sec-total")

		if err != nil {
			log.Println("watch failed in ReleaseLock, err is: ", err)
			return false
		}
	}

	return true
}

// 查看商品是否超发
func (this *Client) CheckIsOverIssued() (keysCount int, sanitizeCount int, subtract int) {
	// list := this.Conn.HGetAll(ctx, "sec-list")
	keys := this.Conn.HKeys(ctx, "sec-list").Val()
	sanitizeKeys := removeDuplicates(keys)

	keysCount = len(keys)
	sanitizeCount = len(sanitizeKeys)
	subtract = len(keys) - len(sanitizeKeys)
	return keysCount, sanitizeCount, subtract
}

// 缓存秒杀商品
// redis-in-action的chapter2已实现

// 基于队列实现
func (this *Client) SecKillQueue() {

}

func removeDuplicates(elements []string) []string { // change string to int here if required
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{} // change string to int here if required
	result := []string{}             // change string to int here if required

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
