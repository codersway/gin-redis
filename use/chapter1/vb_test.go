package chapter1

import (
	"fmt"
	"gin-redis/conf"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/kr/pretty"
)

var (
	conn   = conf.Conn()
	client = NewWeibo(conn)
)

func TestWeibo(t *testing.T) {
	t.Run("准备数据", func(t *testing.T) {
		addWeibo()
	})

	t.Run("写微博", func(t *testing.T) {
		for i := 0; i <= 5; i++ {
			weiboId := client.WriteWeibo("17172", "redis实现自动补全", "https://blog.csdn.net/m0_37556444/article/details/82709764")
			t.Log(weiboId)
			time.Sleep(2 * time.Second)
		}
	})

	t.Run("某条微博详情", func(t *testing.T) {
		detail := client.WeiboDetail("1")
		t.Log(detail)
	})

	t.Run("所有微博列表-feeds", func(t *testing.T) {
		feeds := client.Feeds(2)
		pretty.Log(feeds)
	})

	client.Conn.FlushDB(ctx)
}

func addWeibo() {
	for i := 0; i <= 10; i++ {
		weiboKey := fmt.Sprintf("weibo:%s", grand.Digits(5))
		posterKey := fmt.Sprintf("user:%s", grand.Digits(5))

		now := time.Now().Unix()

		s := map[string]interface{}{"title": "this is a title", "link": "https://learnku.com", "poster": posterKey, "time": now, "votes": 0}

		client.Conn.HSet(ctx, weiboKey, s)

		client.Conn.ZAdd(ctx, SendTimeKey, &redis.Z{Score: float64(now), Member: weiboKey})
	}
}

func TestThumb(t *testing.T) {
	// weiboId := gofunc.RandNum(5)
	weiboId := "03569"
	userId := grand.Digits(5)

	client.ThumbWeibo(weiboId, userId)
}

func TestThumbList(t *testing.T) {
	weiboId := "03561"

	list, err := client.ThumbList(weiboId)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(list)
}
