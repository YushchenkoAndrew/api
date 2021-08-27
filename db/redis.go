package db

import (
	"context"
	"fmt"

	"api/config"
	"api/models"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectToRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})
}

func RedisInitDefault() {
	var nInfo int64
	var nWorld int64

	DB.Model(&models.Info{}).Count(&nInfo)
	DB.Model(&models.World{}).Count(&nWorld)

	// FIXME: ERROR log format
	ctx := context.Background()
	var SetVar = func(ctx *context.Context, param string, value interface{}) {
		if err := Redis.Set(*ctx, param, value, 0).Err(); err != nil {
			fmt.Println("ERROR:")
		}
	}

	SetVar(&ctx, "nInfo", nInfo)
	SetVar(&ctx, "nWorld", nWorld)

	// // Set Default settings configuration
	// SetVar(&ctx, "nItems", conf.Items)
	// SetVar(&ctx, "nLimit", conf.Limit)
	// SetVar(&ctx, "nLiveTime", conf.LiveTime)
	// SetVar(&ctx, "nUsersReq", conf.UsersReq)
}
