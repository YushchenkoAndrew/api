package db

import (
	"context"
	"fmt"

	"api/config"
	"api/logs"
	"api/models"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectToRedis() {
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
}

func RedisInitDefault() {
	var nInfo int64
	var nWorld int64
	var nFile int64
	var nLink int64
	var nProject int64

	DB.Model(&models.Info{}).Count(&nInfo)
	DB.Model(&models.World{}).Count(&nWorld)
	DB.Model(&models.File{}).Count(&nFile)
	DB.Model(&models.Link{}).Count(&nLink)
	DB.Model(&models.Project{}).Count(&nProject)

	// FIXME: ERROR log format
	ctx := context.Background()
	var SetVar = func(ctx *context.Context, param string, value interface{}) {
		if err := Redis.Set(*ctx, param, value, 0).Err(); err != nil {
			fmt.Println("[Redis] Error happed while setting value to Cache")
		}
	}

	SetVar(&ctx, "nInfo", nInfo)
	SetVar(&ctx, "nWorld", nWorld)
	SetVar(&ctx, "nFile", nFile)
	SetVar(&ctx, "nLink", nLink)
	SetVar(&ctx, "nProject", nFile)
	SetVar(&ctx, "Mutex", 1)
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
