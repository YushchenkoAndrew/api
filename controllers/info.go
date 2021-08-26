package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InfoController struct{}

func (*InfoController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	sKeys := ""
	result := db.DB
	id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
	if err == nil && id > 0 {
		sKeys += "id"
		result = result.Where("id = ?", id)
	}

	createdAt := c.DefaultQuery("created_at", "")
	if createdAt != "" {
		sKeys += "created_at"
		result = result.Where("created_at = ?", createdAt)
	}

	countries := c.DefaultQuery("countries", "")
	if countries != "" {
		sKeys += "countries"
		result = result.Where("countries IN ?", strings.Split(countries, ","))
	}

	return result, sKeys
}

func (*InfoController) parseBody(body *models.ReqInfo, model *models.Info) {
	model.Countries = body.Countries

	if body.CreatedAt != nil {
		model.CreatedAt = *body.CreatedAt
	}

	if body.Views != nil {
		model.Views = *body.Views
	}

	if body.Clicks != nil {
		model.Clicks = *body.Clicks
	}

	if body.Media != nil {
		model.Media = *body.Media
	}

	if body.Visitors != nil {
		model.Visitors = *body.Visitors
	}
}

func (o *InfoController) CreateOne(c *gin.Context) {
	var model models.Info
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil || body.Countries == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params or id parm")
		return
	}

	o.parseBody(&body, &model)
	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()

	db.Redis.Incr(ctx, "nInfo")
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *InfoController) CreateAll(c *gin.Context) {
	var body []models.ReqInfo
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var models = make([]models.Info, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Countries == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}

		o.parseBody(&body[i], &models[i])
	}

	result := db.DB.Create(&models)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	// Make an update without stoping the response handler
	go db.Redis.Del(ctx, "Info:Sum")
	go db.RedisAdd(&ctx, "nInfo", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":        "OK",
		"resHandlerult": models,
		"items":         result.RowsAffected,
		"totalItems":    items + result.RowsAffected,
	})
}

func (o *InfoController) ReadOne(c *gin.Context) {
	var id int
	var info []models.Info

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	ctx := context.Background()
	key := "Info:" + strconv.Itoa(id)
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &info)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("id = ?", id).Find(&info)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&info); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(ctx, "nInfo").Int64(); err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     info,
		"items":      1,
		"totalItems": items,
	})
}

func (o *InfoController) ReadAll(c *gin.Context) {
	var info []models.Info
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	result, sKeys := o.filterQuery(c)
	key := "Info:" + c.DefaultQuery(sKeys, "-1")

	if sKeys == "id" || sKeys == "created_at" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &info)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)

			// Update artificially update rows Affected value
			result.RowsAffected = 1
		}
	}

	if len(info) == 0 {
		result = result.Offset(page * config.ENV.Items).Limit(limit).Find(&info)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			return
		}

		if sKeys == "id" || sKeys == "created_at" {
			// Encode json to str
			if str, err := json.Marshal(&info); err == nil {
				go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
			}
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
		"result":     info,
		"page":       page,
		"limit":      limit,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *InfoController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil || !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model models.Info
	o.parseBody(&body, &model)
	result := db.DB.Where("id = ?", id).Updates(&model)
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *InfoController) UpdateAll(c *gin.Context) {
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var sKeys string
	var result *gorm.DB
	var model models.Info

	o.parseBody(&body, &model)
	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Updates(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	var items int64
	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *InfoController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Info{})
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nInfo")
	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *InfoController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.Info{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Del(ctx, "Info:Sum")
	go db.RedisSub(&ctx, "nInfo", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": items - result.RowsAffected,
	})
}
