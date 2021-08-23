package chapter2

import (
	"context"
	"crypto"
	"encoding/hex"
	"github.com/go-redis/redis/v8"
	"net/url"
	"strings"
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

func (this *Client) CachePage(url string, callback func(string) string) string {
	if this.CanCache(url) {
		return callback(url)
	}

	pageKey := "cache:" + hashRequest(url)
	content := this.Conn.Get(ctx, pageKey).Val()

	if content == "" {
		content = callback(url)
		this.Conn.Set(ctx, pageKey, content, 300*time.Second)
	}
	return content
}

func (this *Client) CanCache(url string) bool {
	itemId := extractItemId(url)
	if itemId == "" || isDynamic(url) {
		return false
	}
	rank := this.Conn.ZRank(ctx, "viewed:", itemId).Val()
	return rank != 0 && rank < 10000
}

func extractItemId(request string) string {
	parsed, _ := url.Parse(request)
	queryValue, _ := url.ParseQuery(parsed.RawQuery)
	query := queryValue.Get("item")

	return query
}

func isDynamic(request string) bool {
	parsed, _ := url.Parse(request)
	queryValue, _ := url.ParseQuery(parsed.RawQuery)
	for _, v := range queryValue {
		for _, j := range v {
			if strings.Contains(j, "_") {
				return false
			}
		}
	}

	return true
}

func hashRequest(url string) string {
	// todo
	hash := crypto.MD5.New()
	hash.Write([]byte(url))
	res := hash.Sum(nil)
	return hex.EncodeToString(res)
}
