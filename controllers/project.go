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

type ProjectController struct{}

func (*ProjectController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	sKeys := ""
	result := db.DB
	if id, err := strconv.Atoi(c.DefaultQuery("id", "-1")); err == nil && id > 0 {
		sKeys += "id"
		result = result.Where("id = ?", id)
	}

	if name := c.DefaultQuery("name", ""); name != "" {
		sKeys += "name"
		result = result.Where("name = ?", name)
	}

	if dir := c.DefaultQuery("dir", ""); dir != "" {
		sKeys += "dir"
		result = result.Where("name = ?", dir)
	}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	switch helper.GetStat(start == "", end == "") {
	case 0:
		result = result.Where("created_at <= ? AND created_at >= ?", end, start)
		sKeys = sKeys + "startend"

	case 1:
		result = result.Where("created_at >= ?", start)
		sKeys = sKeys + "start"

	case 2:
		result = result.Where("created_at <= ?", end)
		sKeys = sKeys + "end"

	}

	return result, sKeys
}

func (o *ProjectController) parseBody(body *models.ReqProject, model *models.Project) bool {
	model.Name = body.Name
	model.Dir = body.Dir
	model.Title = body.Title
	model.Desc = body.Desc

	if len(body.Files) != 0 {
		var fileMap = make(map[string]*models.File)

		model.Files = make([]models.File, len(body.Files))
		for i := 0; i < len(body.Files); i++ {
			if ptr, ok := fileMap[body.Files[i].Name]; ok && (*ptr).Role == body.Files[i].Role && (*ptr).Type == body.Files[i].Type {
				return false
			}

			o.parseFileBody(&body.Files[i], &model.Files[i])
			fileMap[model.Files[i].Name] = &model.Files[i]
		}
	}
	return true
}

func (*ProjectController) parseFileBody(body *models.File, model *models.File) {
	model.Name = body.Name
	model.Role = body.Role
	model.Type = body.Type
}

// @Tags Project
// @Summary Create file by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body models.ReqProject true "Project Data"
// @Success 201 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [post]
func (o *ProjectController) CreateOne(c *gin.Context) {
	var body models.ReqProject
	if err := c.ShouldBind(&body); err != nil || body.Name == "" || body.Dir == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body format")
		return
	}

	var result = db.DB
	var model = make([]models.Project, 1)

	if !o.parseBody(&body, &model[0]) {
		helper.ErrHandler(c, http.StatusBadRequest, "Files are repeated")
		return
	}

	result = db.DB.Create(&model)
	if result.Error != nil || result.RowsAffected == 0 {
		helper.ErrHandler(c, http.StatusInternalServerError, "Something unexpected happend")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	db.Redis.Incr(ctx, "nProject")
	items, err := db.Redis.Get(ctx, "nFile").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	go db.Redis.Incr(ctx, "nProject")
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Create Project from list of objects
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body []models.ReqProject true "List of Project Data"
// @Success 201 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/list/{id} [post]
func (o *ProjectController) CreateAll(c *gin.Context) {
	var body []models.ReqProject
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model = make([]models.Project, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Name == "" || body[i].Dir == "" {
			helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
			return
		}

		if !o.parseBody(&body[i], &model[i]) {
			helper.ErrHandler(c, http.StatusBadRequest, "Files are repeated")
			return
		}
	}

	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	go helper.RedisAdd(&ctx, "nProject", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags Project
// @Summary Read Project by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id path int true "Instance id"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/{id} [get]
func (*ProjectController) ReadOne(c *gin.Context) {
	var id int
	var project []models.Project

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	ctx := context.Background()
	key := "Project:" + strconv.Itoa(id)
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &project)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("id = ?", id).Preload("Files").Find(&project)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/project",
				File:    "/controllers/project.go",
				Message: "It's not an error Karl; It's a bug!!",
				Desc:    result.Error,
			})
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&project); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(ctx, "nProject").Int64(); err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     project,
		Items:      1,
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Read Project by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Code Rain'"
// @Param dir query string false "Type: 'CodeRain'"
// @Param start query string false "CreatedAt date >= start"
// @Param end query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [get]
func (o *ProjectController) ReadAll(c *gin.Context) {
	var project []models.Project
	ctx := context.Background()

	page, limit := helper.Pagination(c)
	result, sKeys := o.filterQuery(c)
	key := "Project:" + c.DefaultQuery(sKeys, "-1")

	if sKeys == "id" || sKeys == "name" || sKeys == "start" || sKeys == "end" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &project)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)

			// Update artificially update rows Affected value
			result.RowsAffected = 1
		}
	}

	if len(project) == 0 {
		result = result.Offset(page * config.ENV.Items).Limit(limit).Preload("Files").Find(&project)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/project",
				File:    "/controllers/project.go",
				Message: "It's not an error Karl; It's a bug!!",
				Desc:    result.Error,
			})
			return
		}

		if sKeys == "id" || sKeys == "name" || sKeys == "dir" || sKeys == "start" || sKeys == "end" {
			// Encode json to str
			if str, err := json.Marshal(&project); err == nil {
				go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
			}
		}
	}

	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     project,
		Page:       page,
		Limit:      limit,
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Update Project by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.ReqProject true "Project without File Data"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/{id} [put]
func (o *ProjectController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqProject
	if err := c.ShouldBind(&body); err != nil || !helper.GetID(c, &id) || len(body.Files) != 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	// TODO: Maybe I need to check if new dir is already exits
	// because right now it can be achieved by getting an error
	// from db
	var model models.Project
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
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	db.Redis.Del(ctx, "Project:"+strconv.Itoa(id))

	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
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

// @Tags Project
// @Summary Update Project by Query
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Code Rain'"
// @Param dir query string false "Type: 'CodeRain'"
// @Param model body models.ReqProject true "Project without File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [put]
func (o *ProjectController) UpdateAll(c *gin.Context) {
	var body models.ReqProject
	if err := c.ShouldBind(&body); err != nil || len(body.Files) != 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var sKeys string
	var result *gorm.DB
	var model models.Project

	// TODO: Maybe I need to check if new dir is already exits
	// because right now it can be achieved by getting an error
	// from db
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
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	var items int64
	ctx := context.Background()
	if sKeys == "id" || sKeys == "name" || sKeys == "dir" {
		db.Redis.Del(ctx, "Project:"+c.DefaultQuery(sKeys, ""))
	}

	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
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

// @Tags Project
// @Summary Delete Project and Files by :id
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
// @Router /project/{id} [delete]
func (*ProjectController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	// Delete in both place Files & Project
	if result := db.DB.Where("project_id = ?", id).Delete(&models.File{}); result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Project{})
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id param")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	ctx := context.Background()
	db.Redis.Del(ctx, "Project:"+strconv.Itoa(id))

	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nProject")
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Delete Project by Query and Files with the same project_id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id query int false "Type: '1'"
// @Param name query string false "Type: 'Code Rain'"
// @Param dir query string false "Type: 'CodeRain'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [delete]
func (o *ProjectController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB
	var project []models.Project

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	ctx := context.Background()
	key := "Project:" + c.DefaultQuery(sKeys, "-1")
	if sKeys == "id" || sKeys == "name" || sKeys == "start" || sKeys == "end" {
		// Check if cache have requested data
		if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
			json.Unmarshal([]byte(data), &project)
			go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	if len(project) == 0 {
		if result := result.Find(&project); result.Error != nil || len(project) == 0 {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.SendLogs(&models.LogMessage{
				Stat:    "ERR",
				Name:    "API",
				Url:     "/api/project",
				File:    "/controllers/project.go",
				Message: "It's not an error Karl; It's a bug!!",
				Desc:    result.Error,
			})
			return
		}
	}

	// Delete in both place Files & Project
	if result := db.DB.Where("project_id = ?", project[0].ID).Delete(&models.File{}); result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	result = result.Delete(&models.Project{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			Url:     "/api/project",
			File:    "/controllers/project.go",
			Message: "It's not an error Karl; It's a bug!!",
			Desc:    result.Error,
		})
		return
	}

	if sKeys == "id" || sKeys == "name" || sKeys == "dir" {
		db.Redis.Del(ctx, "Project:"+c.DefaultQuery(sKeys, ""))
	}

	items, err := db.Redis.Get(ctx, "nProject").Int64()
	if err != nil {
		items = -1
		go db.RedisInitDefault()
		go logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/controllers/project.go",
			Message: "Ohh nooo Cache is broken; Anyway...",
			Desc:    err.Error(),
		})
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go helper.RedisSub(&ctx, "nProject", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
