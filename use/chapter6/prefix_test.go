package chapter6

import (
	"context"
	"github.com/go-redis/redis/v8"
	"redis-in-action/conf"
	"testing"
)

var (
	conn   = conf.Conn()
	client = PrefixStruct{
		Conn: conn,
	}
	ctx = context.Background()
)

func TestPrefixAutoComplete(t *testing.T) {
	t.Run("setup测试数据", func(t *testing.T) {
		add()
	})
	t.Run("", func(t *testing.T) {
		client.findPrefixRange("")
	})

}

func add() {
	s := []*redis.Z{
		{Score: 12, Member: "abc"},
		{Score: 13, Member: "abcd"},
		{Score: 14, Member: "abce"},
		{Score: 15, Member: "abck"},
	}
	conn.ZAdd(ctx, "aa", s...)
}
