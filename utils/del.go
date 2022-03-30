package functions

import (
	"context"
	"fmt"
	"gin-redis/conf"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/util/grand"
)

type Conn struct {
	Conn *redis.Client
}

var (
	conn   = conf.Conn()
	client = Conn{
		Conn: conn,
	}
	ctx = context.Background()
)

func BatchInsert() {
	random := grand.Digits(10)

	conn.Set(ctx, random, random, time.Minute*10)
}

func BatchInsert2(values []string) {
	conn.MSet(ctx, values)
}

func BatchInsertSet(key string, values []string) {
	conn.SAdd(ctx, key, values)
}

// scan
func DelKeys(pattern string) bool {
	foundedRecordCount := 0
	iter := conn.Scan(ctx, 0, pattern, 0).Iterator()
	// fmt.Printf("YOUR SEARCH PATTERN= %s\n", pattern)

	for iter.Next(ctx) {
		fmt.Printf("Deleted= %s\n", iter.Val())
		conn.Del(ctx, iter.Val())
		foundedRecordCount++
	}
	if err := iter.Err(); err != nil {
		panic(err)
		return false
	}

	fmt.Printf("Deleted Count %d\n", foundedRecordCount)
	// defer conn.Close()

	return true
}

// sscan+srem
func DelSet(key, pattern string) bool {
	foundedRecordCount := 0
	iter := conn.SScan(ctx, key, 0, pattern, 0).Iterator()
	fmt.Printf("YOUR SEARCH PATTERN= %s\n", key)

	for iter.Next(ctx) {
		fmt.Printf("Deleted= %s\n", iter.Val())
		// conn.Del(ctx, iter.Val())
		conn.SRem(ctx, key, iter.Val())
		foundedRecordCount++
	}
	if err := iter.Err(); err != nil {
		panic(err)
		return false
	}

	fmt.Printf("Deleted Count %d\n", foundedRecordCount)

	return true
}

// hscan+hdel
func DelHash(key, pattern string) bool {
	foundedRecordCount := 0
	iter := conn.HScan(ctx, key, 0, pattern, 0).Iterator()
	fmt.Printf("YOUR SEARCH PATTERN= %s\n", key)

	for iter.Next(ctx) {
		fmt.Printf("Deleted= %s\n", iter.Val())
		// conn.Del(ctx, iter.Val())
		conn.HDel(ctx, key, iter.Val())
		foundedRecordCount++
	}
	if err := iter.Err(); err != nil {
		panic(err)
		return false
	}

	fmt.Printf("Deleted Count %d\n", foundedRecordCount)

	return true
}

// zscan+zrem
func DelZSet(key, pattern string) bool {
	foundedRecordCount := 0
	iter := conn.ZScan(ctx, key, 0, pattern, 0).Iterator()
	fmt.Printf("YOUR SEARCH PATTERN= %s\n", key)

	for iter.Next(ctx) {
		fmt.Printf("Deleted= %s\n", iter.Val())
		// conn.Del(ctx, iter.Val())
		conn.ZRem(ctx, key, iter.Val())
		foundedRecordCount++
	}
	if err := iter.Err(); err != nil {
		panic(err)
		return false
	}

	fmt.Printf("Deleted Count %d\n", foundedRecordCount)

	return true
}

// todo [聊聊go-redis的一些高级用法 – 萌叔](http://vearne.cc/archives/1113)
