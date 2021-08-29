package info

import (
	"api/db"
	"api/helper"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SumController struct{}

func (o *SumController) ReadAll(c *gin.Context) {
	var stat models.StatInfo
	ctx := context.Background()

	if data, err := db.Redis.Get(ctx, "Info:Sum").Result(); err == nil {
		json.Unmarshal([]byte(data), &stat)
	} else {
		db.DB.Table("info").
			Select("SUM(views) as views, SUM(clicks) AS clicks, SUM(media) as media, SUM(visitors) as visitors").
			Scan(&stat)

		// Encode json to str
		if str, err := json.Marshal(&stat); err == nil {
			go db.Redis.Set(ctx, "Info:Sum", str, 0)
		}
	}

	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []models.StatInfo{stat},
		"items":      1,
		"totalItems": items,
	})
}
