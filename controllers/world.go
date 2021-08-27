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
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorldController struct{}

func (*WorldController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	sKeys := ""
	result := db.DB
	id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
	if err == nil && id > 0 {
		sKeys += "id"
		result = result.Where("id = ?", id)
	}

	updatedAll := c.DefaultQuery("updated_at", "")
	if updatedAll != "" {
		sKeys += "update_all"
		result = result.Where("updated_at = ?", updatedAll)
	}

	country := c.DefaultQuery("country", "")
	if country != "" {
		sKeys += "country"
		result = result.Where("country = ?", country)
	}

	return result, sKeys
}

func (*WorldController) parseBody(body *models.ReqWorld, model *models.World) {
	if body.Country != "" {
		model.Country = body.Country
	}

	if body.Visitors != nil {
		model.Visitors = *body.Visitors
	}
}

func (o *WorldController) CreateOne(c *gin.Context) {
	var model models.World
	var body models.ReqWorld
	if err := c.ShouldBind(&body); err != nil || body.Country == "" {
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

	db.Redis.Incr(ctx, "nWorld")
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *WorldController) CreateAll(c *gin.Context) {
	var body []models.ReqWorld
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var models = make([]models.World, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Country == "" {
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
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	// Make an update without stoping the response handler
	go helper.RedisAdd(&ctx, "nWorld", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":        "OK",
		"resHandlerult": models,
		"items":         result.RowsAffected,
		"totalItems":    items + result.RowsAffected,
	})
}

func (*WorldController) ReadOne(c *gin.Context) {
	var id int
	var model []models.World

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	ctx := context.Background()
	key := "World:" + strconv.Itoa(id)
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &model)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("id = ?", id).Find(&model)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&model); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(ctx, "nWorld").Int64(); err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      1,
		"totalItems": items,
	})
}

func (o *WorldController) ReadAll(c *gin.Context) {
	var model []models.World
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	result, sKeys := o.filterQuery(c)
	key := "World:" + c.DefaultQuery(sKeys, "-1")

	// NOTE: Maybe add this feature at some point
	// orderBy := c.DefaultQuery("orderBy", "")
	// desc := c.DefaultQuery("desc", "") != ""

	if sKeys == "id" || sKeys == "updated_at" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &model)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)

			// Update artificially update rows Affected value
			result.RowsAffected = 1
		}
	}

	if len(model) == 0 {
		result = result.Offset(page * config.ENV.Items).Limit(limit).Find(&model)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			return
		}

		if sKeys == "id" || sKeys == "updated_at" {
			// Encode json to str
			if str, err := json.Marshal(&model); err == nil {
				go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
			}
		}
	}

	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
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

func (o *WorldController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqWorld
	if err := c.ShouldBind(&body); err != nil || !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model models.World
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
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *WorldController) UpdateAll(c *gin.Context) {
	var body models.ReqWorld
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var sKeys string
	var result *gorm.DB
	var model models.World

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
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (*WorldController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.World{})
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nWorld")
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": items,
	})
}

func (o *WorldController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.World{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nWorld").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		fmt.Println("Something wrong with Caching!!!")
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go helper.RedisSub(&ctx, "nWorld", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": items - result.RowsAffected,
	})
}
