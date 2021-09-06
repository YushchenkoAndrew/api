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
		sKeys += "updated_at"
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

// @Tags World
// @Summary Create/Update World Data by Country
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body models.ReqWorld true "World Data"
// @Success 201 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world [post]
func (o *WorldController) CreateOne(c *gin.Context) {
	var model = make([]models.World, 1)
	var body models.ReqWorld
	if err := c.ShouldBind(&body); err != nil || body.Country == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params or id parm")
		return
	}

	var result *gorm.DB
	if result = db.DB.Where("country = ?", body.Country).Find(&model); result.RowsAffected == 0 {
		model = make([]models.World, 1)
	}

	o.parseBody(&body, &model[0])
	if result.RowsAffected == 0 {
		result = db.DB.Create(&model)
	} else {
		result = db.DB.Where("country = ?", body.Country).Updates(&model[0])
	}

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

	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Create/Update World from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body []models.ReqWorld true "List of World Data"
// @Success 201 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world/list [post]
func (o *WorldController) CreateAll(c *gin.Context) {
	var body []models.ReqWorld
	if err := c.ShouldBind(&body); err != nil || len(body) == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model []models.World
	var countries = []string{}
	for i := 0; i < len(body); i++ {
		if body[i].Country == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}
		countries = append(countries, body[i].Country)
	}

	var result *gorm.DB
	result = db.DB.Where("country IN ?", countries).Find(&model)
	var modelToCreate = make([]models.World, len(body)-int(result.RowsAffected))

	// Init hash for simplify intersection detection
	var hash = map[string]bool{}
	for i := 0; i < int(result.RowsAffected); i++ {
		hash[model[i].Country] = true
	}

	var j = 0
	for i := 0; i < len(body); i++ {
		var m models.World
		o.parseBody(&body[i], &m)

		if _, ok := hash[body[i].Country]; !ok {
			modelToCreate[j] = m
			j++
		} else if result = db.DB.Where("country = ?", body[i].Country).Updates(&m); result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong2")
			return
		}
	}

	if result = db.DB.Create(&modelToCreate); result.Error != nil {
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
	go helper.RedisAdd(&ctx, "nWorld", int64(len(body)))
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags World
// @Summary Read World by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world/{id} [get]
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

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      1,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Read All World
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Instance :id"
// @Param updated_at query string false "UpdatedAt date"
// @Param country query string false "Country: 'UK'"
// @Success 200 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world [get]
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
		if page != -1 {
			result = result.Offset(page * config.ENV.Items).Limit(limit).Find(&model)
		} else {
			result = result.Find(&model)
		}

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

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Update World by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.ReqWorld true "World Data"
// @Success 200 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world/{id} [put]
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

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Update World by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param updated_at query string false "UpdatedAt date"
// @Param country query string false "Country: 'UK'"
// @Param model body models.ReqWorld true "World Data"
// @Success 200 {object} models.Success{result=[]models.World}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world [put]
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

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Delete World by :id
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
// @Router /world/{id} [delete]
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
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags World
// @Summary Delete World by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param updated_at query string false "UpdatedAt date"
// @Param country query string false "Country: 'UK'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /world [delete]
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
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
