package conf

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

var (
	Addr     = "localhost:6379"
	Password = "root"
	DB       = 1
)

var ctx = context.Background()

func Conn() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password,
		DB:       DB,
	})

	if _, err := conn.Ping(ctx).Result(); err != nil {
		log.Fatalf("connect to redis client failed, err: %v \n", err)
	}

	return conn
}
