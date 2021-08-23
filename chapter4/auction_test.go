package chapter4

import (
	"github.com/91go/gobase/sz/redis-in-action/conf"
	"testing"
)

func TestAuction(t *testing.T) {
	conn := conf.Conn()
	client := NewClient(conn)

	t.Run("", func(t *testing.T) {
		client.ListItem("itemX", "userX", 10)
		buyer := "userY"
		client.Conn.HSet(ctx, "users:userY", "funds", 125)
		r := client.Conn.HGetAll(ctx, "users:userY").Val()
		t.Log("this user's money: ", r)

		p := client.PurchaseItem("userY", "itemX", "userX", 10)
		t.Log("purchasing an item succeeded?", p)

		r = client.Conn.HGetAll(ctx, "users:userY").Val()
		t.Log("their money is now:", r)

		i := client.Conn.SMembers(ctx, "inventory:"+buyer).Val()
		t.Log("their inventory is now:", i)
	})

	defer client.Conn.FlushDB(ctx)
}
