package info

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/models"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RangeController struct{}

func (*RangeController) Read(c *gin.Context) {
	var model []models.Info
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	orderBy := c.DefaultQuery("orderBy", "")
	desc := c.DefaultQuery("desc", "") != ""

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	if start == "" && end == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Start or End query should be set")
		return
	}

	result := db.DB.Offset(page * config.ENV.Items).Limit(limit)
	switch helper.GetStat(start == "", end == "") {
	case 0:
		result = result.Where("created_at <= ? AND created_at >= ?", end, start)

	case 1:
		result = result.Where("created_at > ?", start)

	case 2:
		result = result.Where("created_at < ?", end)

	}

	if orderBy != "" {
		if desc {
			result = result.Order(helper.ToSnakeCase(orderBy) + " DESC")
		} else {
			result = result.Order(helper.ToSnakeCase(orderBy) + " ASC")
		}
	}

	result = result.Find(&model)
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil || result == nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"page":       page,
		"limit":      limit,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}
