package helper

import (
	"api/db"
	"api/logs"
	"api/models"
	"context"
)

func RedisAdd(ctx *context.Context, sParam string, nSize int64) {
	var mutex int64
	var err error

	for {
		if mutex, err = db.Redis.Get(*ctx, "Mutex").Int64(); err != nil {
			go db.RedisInitDefault()
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				File:    "/helper/redis.go",
				Message: "Ohh nooo Cache is broken; Anyway...",
				Desc:    err.Error(),
			})
			break
		}

		if mutex%2 != 0 {
			db.Redis.Incr(*ctx, "Mutex")
			if num, err := db.Redis.Get(*ctx, sParam).Int64(); err == nil {
				db.Redis.Set(*ctx, sParam, num+nSize, 0)
			}

			db.Redis.Incr(*ctx, "Mutex")
			break
		}
	}
}

func RedisSub(ctx *context.Context, sParam string, nSize int64) {
	RedisAdd(ctx, sParam, -nSize)
}
