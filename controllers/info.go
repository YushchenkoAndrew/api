package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/logs"
	"api/models"
	"context"
	"encoding/json"
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

	// FIXME: If need it
	// if body.CreatedAt != nil {
	// 	model.CreatedAt = *body.CreatedAt
	// }

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

// @Tags Info
// @Summary Create one instace of Info
// @Description 'CreatedAt' setted automatically
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body models.ReqInfo true "Info Data"
// @Success 201 {object} models.Success{result=[]models.Info}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info [post]
func (o *InfoController) Create(c *gin.Context) {
	var model = make([]models.Info, 1)
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil || body.Countries == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params or id parm")
		return
	}

	o.parseBody(&body, &model[0])
	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	go db.FlushValue("Info")

	ctx := context.Background()
	db.Redis.Incr(ctx, "nInfo")
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Create/Update Info
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param date path string true "Created at instance"
// @Param model body models.ReqInfo true "Info Data"
// @Success 201 {object} models.Success{result=[]models.Info}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/{date} [post]
func (o *InfoController) CreateOne(c *gin.Context) {
	var date = c.Param("date")
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil || body.Countries == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body is not setted")
		return
	}

	created, err := time.Parse("2006-01-02", date)
	if err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect date param")
		return
	}

	var result = db.DB
	var model []models.Info
	ctx := context.Background()

	key := "Info:" + date
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &model)
		result.RowsAffected = int64(len(model))
	} else if result = db.DB.Where("created_at = ?", date).Find(&model); result.RowsAffected == 0 {
		model = make([]models.Info, 1)
	}

	model[0].CreatedAt = created
	o.parseBody(&body, &model[0])
	if result.RowsAffected == 0 {
		result = db.DB.Create(&model)
		db.Redis.Incr(ctx, "nInfo")
	} else {
		result = db.DB.Where("created_at = ?", date).Updates(&model[0])
	}

	if result.Error != nil || result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	// Encode json to str
	if str, err := json.Marshal(&model); err == nil {
		go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
	}

	go db.FlushValue("Info")

	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	// Make an update without stoping the response handler
	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Create Info from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body []models.ReqInfo true "List of Info Data"
// @Success 201 {object} models.Success{result=[]models.Info}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/list [post]
func (o *InfoController) CreateAll(c *gin.Context) {
	var body []models.ReqInfo
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model = make([]models.Info, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Countries == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}

		o.parseBody(&body[i], &model[i])
	}

	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	// Make an update without stoping the response handler
	go db.Redis.Del(ctx, "Info:Sum")
	go helper.RedisAdd(&ctx, "nInfo", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags Info
// @Summary Read Info by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]models.Info}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/{id} [get]
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
			go logs.DefaultLog("/controllers/info.go", result.Error)
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
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     info,
		Items:      1,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Read All Info
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Instance :id"
// @Param created_at query string false "CreatedAt date"
// @Param countries query string false "Countries: 'UK,US'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} models.Success{result=[]models.Info}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info [get]
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
			go logs.DefaultLog("/controllers/info.go", result.Error)
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
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     info,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Update Info by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.ReqInfo true "Info Data"
// @Success 200 {object} models.Success{result=[]models.Info}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/{id} [put]
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
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	go db.FlushValue("Info")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     [1]models.Info{model},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Update Info by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param created_at query string false "CreatedAt date"
// @Param countries query string false "Countries: 'UK,US'"
// @Param model body models.ReqInfo true "Info Data"
// @Success 200 {object} models.Success{result=[]models.Info}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info [put]
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
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	go db.FlushValue("Info")

	var items int64
	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     [1]models.Info{model},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Delete Info by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info/{id} [delete]
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
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	go db.FlushValue("Info")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nInfo")
	go db.Redis.Del(ctx, "Info:Sum")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Info
// @Summary Delete Info by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param created_at query string false "CreatedAt date"
// @Param countries query string false "Countries: 'UK,US'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /info [delete]
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
		go logs.DefaultLog("/controllers/info.go", result.Error)
		return
	}

	go db.FlushValue("Info")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nInfo").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/info.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Del(ctx, "Info:Sum")
	go helper.RedisSub(&ctx, "nInfo", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
