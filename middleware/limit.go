package middleware

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := context.Background()

		rate, err := db.Redis.Get(ctx, ip).Int()
		if err != nil {
			go db.Redis.Set(ctx, ip, 1, time.Duration(config.ENV.RateTime)*time.Second)
			return
		}

		go db.Redis.Incr(ctx, ip)
		if rate >= config.ENV.RateLimit {
			helper.ErrHandler(c, http.StatusTooManyRequests, "Toggled Reqest rate limiter")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "OK",
				Name:    "API",
				Url:     "/api/refresh",
				File:    "/middleware/limit.go",
				Message: "Jeez man calm down, you've had inaff of traffic, I'm blocking you; ip=" + ip,
			})
		}
	}
}
