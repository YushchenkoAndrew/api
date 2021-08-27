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

func (o *SumController) Read(c *gin.Context) {
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

// TODO: Test this!!
func (o *SumController) Update(c *gin.Context) {
	var stat models.StatInfo
	var body models.ReqInfo
	var date = c.DefaultQuery("date", "")
	if err := c.ShouldBind(&body); err != nil || date == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params or date is not setted")
		return
	}

	var model models.Info
	ctx := context.Background()
	// if data, err := db.Redis.Get(ctx, "Info:Sum").Result(); err == nil {
	// 	json.Unmarshal([]byte(data), &stat)
	// } else {
	// 	db.DB.Table("info").
	// 		Select("SUM(views) as views, SUM(clicks) AS clicks, SUM(media) as media, SUM(visitors) as visitors").
	// 		Scan(&stat)
	// }

	result := db.DB.Where("created_at = ?", date).Find(&model)

	// FIXME: Country
	model.Countries = model.Countries + body.Countries
	model.Clicks = model.Clicks + *body.Clicks - stat.Clicks
	model.Views = model.Views + *body.Views - stat.Views
	model.Media = model.Media + *body.Media - stat.Media
	model.Visitors = model.Visitors + *body.Visitors - stat.Visitors

	if result.RowsAffected == 0 {
		result = db.DB.Create(&model)
	} else {
		result = db.DB.Where("created_at = ?", date).Updates(&model)
	}

	if result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		return
	}

	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []models.Info{model},
		"items":      1,
		"totalItems": items,
	})
}
