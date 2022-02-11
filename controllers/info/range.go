package info

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces/info"
	"api/logs"
	"api/models"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type rangeController struct{}

func NewRangeController() info.Default {
	return &rangeController{}
}

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
func (*rangeController) Read(c *gin.Context) {
	var model []models.Info

	page, limit := helper.Pagination(c)
	orderBy := c.DefaultQuery("orderBy", "")
	desc := c.DefaultQuery("desc", "") != ""

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	if start == "" && end == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Start or End query should be set")
		return
	}

	var keys = []string{}
	client := db.DB.Offset(page * config.ENV.Items).Limit(limit)
	switch helper.GetStat(start == "", end == "") {
	case 0:
		keys = append(keys, fmt.Sprintf("CREATED_AT<=%s:CREATED_AT=>%s", end, start))
		client = client.Where("created_at <= ? AND created_at >= ?", end, start)

	case 1:
		keys = append(keys, fmt.Sprintf("CREATED_AT>=%s", start))
		client = client.Where("created_at >= ?", start)

	case 2:
		keys = append(keys, fmt.Sprintf("CREATED_AT<=%s", end))
		client = client.Where("created_at <= ?", end)

	}

	if orderBy != "" {
		if desc {
			keys = append(keys, "DESC")
			client = client.Order(helper.ToSnakeCase(orderBy) + " DESC")
		} else {
			keys = append(keys, "ASC")
			client = client.Order(helper.ToSnakeCase(orderBy) + " ASC")
		}
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(keys, ":")))
	if err := helper.PrecacheResult(fmt.Sprintf("INFO:RANGE:%s", hex.EncodeToString(hasher.Sum(nil))), client, model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/info/range.go", err.Error())
		return
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status: "OK",
		Result: model,
		Page:   page,
		Limit:  limit,
		Items:  int64(len(model)),
		// TotalItems: items,
	})
}
