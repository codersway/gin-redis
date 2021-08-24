package throttles

import (
	"github.com/91go/gobase/sz/redis-in-action/conf"
	"testing"
	"time"
)

var (
	conn   = conf.Conn()
	client = Throttles{
		Conn: conn,
	}
)

func TestThrottles_ThrottlesStr(t *testing.T) {
	for i := 0; i <= 150; i++ {
		echo := client.ThrottlesStr()
		time.Sleep(time.Second)
		t.Log(echo)
	}
}

func TestThrottles_ThrottlesLua(t *testing.T) {
	for i := 0; i <= 150; i++ {
		echo := client.ThrottlesLua()
		time.Sleep(time.Millisecond)
		t.Log(echo)
	}
}

func TestThrottles_ThrottlesZset(t *testing.T) {
	for i := 0; i <= 150; i++ {
		echo := client.ThrottlesZset()

		t.Log(echo)
	}
}
