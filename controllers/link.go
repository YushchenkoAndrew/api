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
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LinkController struct{}

func (*LinkController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	sKeys := ""
	result := db.DB
	if id, err := strconv.Atoi(c.DefaultQuery("id", "-1")); err == nil && id > 0 {
		sKeys += "id"
		result = result.Where("id = ?", id)
	}

	if projectId, err := strconv.Atoi(c.DefaultQuery("project_id", "-1")); err == nil && projectId > 0 {
		sKeys += "project_id"
		result = result.Where("project_id = ?", projectId)
	}

	if name := c.DefaultQuery("name", ""); name != "" {
		sKeys += "name"
		result = result.Where("name = ?", name)
	}

	return result, sKeys
}

func (*LinkController) parseBody(body *models.ReqLink, model *models.Link) {
	model.Name = body.Name
	model.Link = body.Link
}

func (*LinkController) isExist(id int, body *models.ReqLink) bool {
	var model []models.Link
	res := db.DB.Where("project_id = ? AND name = ? AND Link = ?", id, body.Name, body.Link).Find(&model)
	return !(res.RowsAffected == 0)
}

// @Tags Link
// @Summary Create link by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body models.ReqLink true "Link info"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [post]
func (o *LinkController) CreateOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	var body models.ReqLink
	if err := c.ShouldBind(&body); err != nil || body.Name == "" || body.Link == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body format")
		return
	}

	var result = db.DB
	var model []models.Link
	ctx := context.Background()

	// Check if such project is exist
	var key = "Project:" + strconv.Itoa(id)
	if _, err := db.Redis.Get(ctx, key).Result(); err != nil {
		var project []models.Project
		if res := db.DB.Where("id = ?", id).Find(&project); res.RowsAffected == 0 {
			helper.ErrHandler(c, http.StatusBadRequest, "Unknown project id")
			return
		} else if str, err := json.Marshal(&project); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	if o.isExist(id, &body) {
		helper.ErrHandler(c, http.StatusBadRequest, "Such link already exist")
		return
	}

	model = make([]models.Link, 1)
	model[0].ProjectID = uint32(id)

	o.parseBody(&body, &model[0])
	result = db.DB.Create(&model)

	if result.Error != nil || result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	db.Redis.Incr(ctx, "nLink")
	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	go db.Redis.Incr(ctx, "nLink")
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     &model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Create File from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project id"
// @Param model body []models.ReqLink true "List of Links info"
// @Success 201 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/list/{id} [post]
func (o *LinkController) CreateAll(c *gin.Context) {
	var err error
	var id int
	var body []models.ReqLink

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	if err = c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	ctx := context.Background()
	var key = "Project:" + strconv.Itoa(id)
	if _, err := db.Redis.Get(ctx, key).Result(); err != nil {
		var project []models.Project
		if res := db.DB.Where("id = ?", id).Find(&project); res.RowsAffected == 0 {
			helper.ErrHandler(c, http.StatusBadRequest, "Unknown project id")
			return
		} else if str, err := json.Marshal(&project); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var model = make([]models.Link, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Name == "" || body[i].Link == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}

		if o.isExist(id, &body[i]) {
			helper.ErrHandler(c, http.StatusBadRequest, "Such link already exist")
			return
		}

		o.parseBody(&body[i], &model[i])
		model[i].ProjectID = uint32(id)
	}

	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	go helper.RedisAdd(&ctx, "nLink", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     &model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags Link
// @Summary Read Link by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]models.Link}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [get]
func (*LinkController) ReadOne(c *gin.Context) {
	var id int
	var links []models.Link

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	ctx := context.Background()
	key := "Link:" + strconv.Itoa(id)
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &links)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("id = ?", id).Find(&links)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.DefaultLog("/controllers/link.go", result.Error)
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&links); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(ctx, "nLink").Int64(); err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &links,
		Items:      1,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Read Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} models.Success{result=[]models.Link}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link [get]
func (o *LinkController) ReadAll(c *gin.Context) {
	var links []models.Link
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	result, sKeys := o.filterQuery(c)
	key := "Link:" + c.DefaultQuery(sKeys, "-1")

	if sKeys == "id" || sKeys == "project_id" || sKeys == "name" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &links)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)

			// Update artificially update rows Affected value
			result.RowsAffected = 1
		}
	}

	if len(links) == 0 {
		result = result.Offset(page * config.ENV.Items).Limit(limit).Find(&links)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.DefaultLog("/controllers/link.go", result.Error)
			return
		}

		if sKeys == "id" || sKeys == "project_id" || sKeys == "name" {
			// Encode json to str
			if str, err := json.Marshal(&links); err == nil {
				go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
			}
		}
	}

	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &links,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Update Link by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.ReqLink true "Link Data"
// @Success 200 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link/{id} [put]
func (o *LinkController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqLink
	if err := c.ShouldBind(&body); err != nil || !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model models.Link
	o.parseBody(&body, &model)
	result := db.DB.Where("id = ?", id).Updates(&model)
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &[1]models.Link{model},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Update Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Param model body models.ReqLink true "Link Data"
// @Success 200 {object} models.Success{result=[]models.Link}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link [put]
func (o *LinkController) UpdateAll(c *gin.Context) {
	var body models.ReqLink
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var sKeys string
	var result *gorm.DB
	var model models.Link

	o.parseBody(&body, &model)
	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Updates(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	var items int64
	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &[1]models.Link{model},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Delete Link by :id
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
// @Router /link/{id} [delete]
func (*LinkController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Link{})
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nLink")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &[]string{},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Link
// @Summary Delete Link by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param name query string false "Type: 'Name: 'main'"
// @Param project_id query string false "ProjectID: '1'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /link [delete]
func (o *LinkController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.Link{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/link.go", result.Error)
		return
	}

	go db.FlushValue("Link")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nLink").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.DefaultLog("/controllers/link.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go helper.RedisSub(&ctx, "nFile", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     &[]string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
