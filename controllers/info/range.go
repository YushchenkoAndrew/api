package info

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RangeController struct{}

// @Tags Info
// @Summary Get Info data by date Range
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param start query string true "CreatedAt date >= start"
// @Param end query string false "CreatedAt date <= end"
// @Param orderBy query string false "Column which final result will be sorted by"
// @Param desc query int false "Sort by direction"
// @Success 200 {object} models.Success{result=[]models.Info}
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/range [get]
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
		result = result.Where("created_at >= ?", start)

	case 2:
		result = result.Where("created_at <= ?", end)

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
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/info/range.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}
