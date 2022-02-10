package db

import (
	"context"
	"fmt"

	"api/config"
	"api/interfaces"
	"api/logs"
	"api/models"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectToRedis(tables []interfaces.Table) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})

	ctx := context.Background()
	if _, err := Redis.Ping(ctx).Result(); err != nil {
		logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/redis.go",
			Message: "Bruhhh, did you even start the Redis ???",
			Desc:    err.Error(),
		})
		panic("Failed on Redis connection")
	}

	Redis.Set(ctx, "Mutex", 1, 0)
	for _, table := range tables {
		if err := table.Redis(DB, Redis); err != nil {
			panic(err)
		}
	}
}

func FlushValue(key string) {
	ctx := context.Background()
	iter := Redis.Scan(ctx, 0, key+":*", 0).Iterator()

	for iter.Next(ctx) {
		Redis.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		fmt.Println("[Redis] Error happed while setting interating through keys")
	}
}
