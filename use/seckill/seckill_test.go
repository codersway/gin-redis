package seckill

import (
	"gin-redis/conf"
	"testing"
)

func BenchmarkSecKill(b *testing.B) {
	conn := conf.Conn()
	client := NewClient(conn)

	// 358799 ns/op
	// 3k左右的qps
	b.Run("", func(b2 *testing.B) {
		for i := 0; i <= b2.N; i++ {
			client.SecKill(i)
		}

		keysCount, sanitizeCount, subtract := client.CheckIsOverIssued()
		b2.Logf("共%d个key, 未重复%d个key，重复key共%d个", keysCount, sanitizeCount, subtract)
	})

	client.Conn.FlushDB(ctx)
}
