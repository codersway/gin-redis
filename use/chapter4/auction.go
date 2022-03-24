package chapter4

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// 实现游戏物品拍卖行
type Client struct {
	Conn *redis.Client
}

var ctx = context.Background()

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

func (this *Client) ListItem(itemId, sellerId string, price float64) bool {
	inventory := fmt.Sprintf("inventory:%s", sellerId)
	item := fmt.Sprintf("%s.%s", itemId, sellerId)
	end := time.Now().Unix() + 5

	for time.Now().Unix() < end {
		err := this.Conn.Watch(ctx, func(tx *redis.Tx) error {
			if _, err := tx.TxPipelined(ctx, func(pipeline redis.Pipeliner) error {
				if !tx.SIsMember(ctx, inventory, itemId).Val() {
					tx.Unwatch(ctx, inventory)
					return nil
				}
				pipeline.ZAdd(ctx, "market:", &redis.Z{Member: item, Score: price})
				pipeline.SRem(ctx, inventory, itemId)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}, inventory)
		if err != nil {
			log.Println("watch err:", err)
			return false
		}

		return true
	}

	return false
}

func (this Client) PurchaseItem(buyerId, itemId, sellerId string, lPrice int64) bool {
	buyer := fmt.Sprintf("users:%s", buyerId)
	seller := fmt.Sprintf("users:%s", sellerId)
	item := fmt.Sprintf("%s.%s", itemId, sellerId)
	inventory := fmt.Sprintf("inventory:%s", buyerId)
	end := time.Now().Unix() + 10

	for time.Now().Unix() < end {
		err := this.Conn.Watch(ctx, func(tx *redis.Tx) error {
			if _, err := tx.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
				price := int64(pipeliner.ZScore(ctx, "market:", item).Val())
				funds, _ := tx.HGet(ctx, buyer, "funds").Int64()
				if price != lPrice || price > funds {
					tx.Unwatch(ctx)
				}

				pipeliner.HIncrBy(ctx, seller, "funds", price)
				pipeliner.HIncrBy(ctx, buyer, "funds", -price)
				pipeliner.SAdd(ctx, inventory, itemId)
				pipeliner.ZRem(ctx, "market:", item)
				return nil
			}); err != nil {
				return err
			}
			return nil
		}, "market:", buyer)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false
}
