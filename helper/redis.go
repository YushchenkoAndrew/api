package helper

import (
	"api/db"
	"context"
)

func RedisAdd(ctx *context.Context, sParam string, nSize int64) {
	for i := 0; i < int(nSize); i++ {
		db.Redis.Incr(*ctx, sParam)
	}
}

func RedisSub(ctx *context.Context, sParam string, nSize int64) {
	for i := 0; i < int(nSize); i++ {
		db.Redis.Decr(*ctx, sParam)
	}
}
