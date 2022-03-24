package chapter5

import (
	"fmt"
	"gin-redis/conf"
	"math/rand"
	"strconv"
	"testing"
)

func TestIP(t *testing.T) {
	conn := conf.Conn()
	client := NewClient(conn)

	t.Run("", func(t *testing.T) {
		client.ImportIpsToRedis("xxx.csv")
		ranges := client.Conn.ZCard(ctx, "ip2cityid:").Val()
		t.Log("loaded ranges into redis:", ranges)

		client.ImportCityToRedis("xxx.csv")
		cities := client.Conn.HLen(ctx, "cityId2City:").Val()
		t.Log("Loaded city lookups into Redis:", cities)

		for i := 0; i < 5; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%s", strconv.Itoa(rand.Intn(254)+1), randString(256), randString(256), randString(256))
			t.Log(ip, client.FindCityByIp(ip))
		}
	})

	defer client.Conn.FlushDB(ctx)
}

func randString(up int) string {
	rand.Seed(rand.Int63())
	return strconv.Itoa(rand.Intn(up))
}
