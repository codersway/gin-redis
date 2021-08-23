package chapter2

import (
	"github.com/91go/gobase/sz/redis-in-action/conf"
	"testing"
)

func TestCache(t *testing.T) {
	conn := conf.Conn()
	client := NewClient(conn)

	t.Run("", func(t *testing.T) {
		url := "http://nstool.zhuanzfx.com/"
		result := client.CachePage(url, func(url string) string {
			return "content for:" + url
		})

		t.Log("init ctx: ", result)
		//assert.Equal(t, )
	})
}
