package db

import (
	"context"
	"fmt"

	"api/config"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectToRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})
}

func TestRedis() {
	ctx := context.Background()
	fmt.Println("Test Redis")
	err := Redis.Set(ctx, "test", "HELLO WORLD", 0).Err()
	if err != nil {
		panic(err)
	}

	val := Redis.Get(ctx, "test").String()
	fmt.Println("Get 'test' = ", val)
}
