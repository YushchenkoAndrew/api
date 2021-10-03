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

type FileController struct{}

func (*FileController) filterQuery(c *gin.Context) (*gorm.DB, string) {
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

	if typeName := c.DefaultQuery("type", ""); typeName != "" {
		sKeys += "type"
		result = result.Where("type IN ?", strings.Split(typeName, ","))
	}

	if role := c.DefaultQuery("role", ""); role != "" {
		sKeys += "role"
		result = result.Where("role = ?", role)
	}

	return result, sKeys
}

func (*FileController) parseBody(body *models.ReqFile, model *models.File) {
	model.Name = body.Name
	model.Role = body.Role
	model.Type = body.Type
}

func (*FileController) isExist(id int, body *models.ReqFile) bool {
	var model []models.File
	res := db.DB.Where("project_id = ? AND name = ? AND role = ? AND type = ?", id, body.Name, body.Role, body.Type).Find(&model)
	return !(res.RowsAffected == 0)
}

// @Tags File
// @Summary Create file by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project primaray id"
// @Param model body models.ReqFile true "File Data"
// @Success 201 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/{id} [post]
func (o *FileController) CreateOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	var body models.ReqFile
	if err := c.ShouldBind(&body); err != nil || body.Name == "" || body.Role == "" || body.Type == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body format")
		return
	}

	var result = db.DB
	var model []models.File
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
		helper.ErrHandler(c, http.StatusBadRequest, "Such file already exist")
		return
	}

	model = make([]models.File, 1)
	model[0].ProjectID = uint32(id)

	o.parseBody(&body, &model[0])
	result = db.DB.Create(&model)

	if result.Error != nil || result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	db.Redis.Incr(ctx, "nFile")
	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	go db.Redis.Incr(ctx, "nFile")
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Create File from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Project id"
// @Param model body []models.ReqFile true "List of File Data"
// @Success 201 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/list/{id} [post]
func (o *FileController) CreateAll(c *gin.Context) {
	var err error
	var id int
	var body []models.ReqFile

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

	var model = make([]models.File, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Name == "" || body[i].Role == "" || body[i].Type == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}

		if o.isExist(id, &body[i]) {
			helper.ErrHandler(c, http.StatusBadRequest, "Such file already exist")
			return
		}

		o.parseBody(&body[i], &model[i])
		model[i].ProjectID = uint32(id)
	}

	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	go helper.RedisAdd(&ctx, "nFile", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags File
// @Summary Read File by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/{id} [get]
func (*FileController) ReadOne(c *gin.Context) {
	var id int
	var files []models.File

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	ctx := context.Background()
	key := "File:" + strconv.Itoa(id)
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &files)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("id = ?", id).Find(&files)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/file",
				File:    "/controllers/file.go",
				Message: "It's not an error Karl; It's a bug!!",
				Desc:    result.Error,
			})
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&files); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(ctx, "nFile").Int64(); err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     files,
		Items:      1,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Read File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param type query string false "Type: 'js,html,img'"
// @Param role query string false "Role: 'src,assets,styles'"
// @Param project_id query string false "ProjectID: '1'"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file [get]
func (o *FileController) ReadAll(c *gin.Context) {
	var files []models.File
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	result, sKeys := o.filterQuery(c)
	key := "File:" + c.DefaultQuery(sKeys, "-1")

	if sKeys == "id" || sKeys == "project_id" || sKeys == "role" || sKeys == "name" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &files)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)

			// Update artificially update rows Affected value
			result.RowsAffected = 1
		}
	}

	if len(files) == 0 {
		result = result.Offset(page * config.ENV.Items).Limit(limit).Find(&files)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/file",
				File:    "/controllers/file.go",
				Message: "It's not an error Karl; It's a bug!!",
				Desc:    result.Error,
			})
			return
		}

		if sKeys == "id" || sKeys == "project_id" || sKeys == "role" || sKeys == "name" {
			// Encode json to str
			if str, err := json.Marshal(&files); err == nil {
				go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
			}
		}
	}

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     files,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Update File by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.ReqFile true "File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/{id} [put]
func (o *FileController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqFile
	if err := c.ShouldBind(&body); err != nil || !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model models.File
	o.parseBody(&body, &model)
	result := db.DB.Where("id = ?", id).Updates(&model)
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	db.Redis.Del(ctx, "File:"+strconv.Itoa(id))

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Update File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param type query string false "Type: 'js,html,img'"
// @Param role query string false "Role: 'src,assets,styles'"
// @Param project_id query string false "ProjectID: '1'"
// @Param model body models.ReqFile true "File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file [put]
func (o *FileController) UpdateAll(c *gin.Context) {
	var body models.ReqFile
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var sKeys string
	var result *gorm.DB
	var model models.File

	o.parseBody(&body, &model)
	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Updates(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	var items int64
	ctx := context.Background()
	if sKeys == "id" || sKeys == "project_id" || sKeys == "role" || sKeys == "name" {
		db.Redis.Del(ctx, "File:"+c.DefaultQuery(sKeys, ""))
	}

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Delete File by :id
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
// @Router /file/{id} [delete]
func (*FileController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.File{})
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	db.Redis.Del(ctx, "File:"+strconv.Itoa(id))

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nFile")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags File
// @Summary Delete File by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Instance :id"
// @Param type query string false "Type: 'js,html,img'"
// @Param role query string false "Role: 'src,assets,styles'"
// @Param project_id query string false "ProjectID: '1'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file [delete]
func (o *FileController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.File{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/file",
			File:    "/controllers/file.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	if sKeys == "id" || sKeys == "project_id" || sKeys == "role" || sKeys == "name" {
		db.Redis.Del(ctx, "File:"+c.DefaultQuery(sKeys, ""))
	}

	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/file.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go helper.RedisSub(&ctx, "nFile", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
