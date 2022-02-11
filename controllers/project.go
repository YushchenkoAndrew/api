package controllers

import (
	"api/config"
	"api/db"
	"api/helper"
	"api/interfaces"
	"api/logs"
	"api/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type projectController struct{}

func NewProjectController() interfaces.Default {
	return &projectController{}
}

func (*projectController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	var keys = []string{}
	client := db.DB

	if id, err := strconv.Atoi(c.DefaultQuery("id", "-1")); err == nil && id > 0 {
		keys = append(keys, fmt.Sprintf("ID=%d", id))
		client = client.Where("id = ?", id)
	}

	if name := c.DefaultQuery("name", ""); name != "" {
		keys = append(keys, fmt.Sprintf("NAME=%s", name))
		client = client.Where("name = ?", name)
	}

	if title := c.DefaultQuery("title", ""); title != "" {
		keys = append(keys, fmt.Sprintf("TITLE=%s", title))
		client = client.Where("title = ?", title)
	}

	start := c.DefaultQuery("start", "")
	end := c.DefaultQuery("end", "")

	switch helper.GetStat(start == "", end == "") {
	case 0:
		keys = append(keys, fmt.Sprintf("CREATED_AT<=%s:CREATED_AT>=%s", end, start))
		client = client.Where("created_at <= ? AND created_at >= ?", end, start)

	case 1:
		keys = append(keys, fmt.Sprintf("CREATED_AT>=%s", start))
		client = client.Where("created_at >= ?", start)

	case 2:
		keys = append(keys, fmt.Sprintf("CREATED_AT<=%s", end))
		client = client.Where("created_at <= ?", end)

	}

	if len(keys) == 0 {
		return client, ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(keys, ":")))
	return client, hex.EncodeToString(hasher.Sum(nil))
}

func (*projectController) filterFileQuery(c *gin.Context) ([]interface{}, string) {
	var keys = []string{}

	var args = []interface{}{}
	var conditions = []string{}

	if typeName := c.DefaultQuery("type", ""); typeName != "" {
		keys = append(keys, fmt.Sprintf("TYPE=%s", typeName))
		args = append(args, typeName)
		conditions = append(conditions, "type = ?")
	}

	if role := c.DefaultQuery("role", ""); role != "" {
		keys = append(keys, fmt.Sprintf("ROLE=%s", role))
		args = append(args, role)
		conditions = append(conditions, "role = ?")
	}

	if len(conditions) == 0 {
		return []interface{}{}, ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(keys, ":")))
	return append([]interface{}{strings.Join(conditions, " AND ")}, args...), hex.EncodeToString(hasher.Sum(nil))
}

func (*projectController) filterLinkQuery(c *gin.Context) ([]interface{}, string) {
	var keys = []string{}

	var args = []interface{}{}
	var conditions = []string{}

	if name := c.DefaultQuery("link_name", ""); name != "" {
		keys = append(keys, fmt.Sprintf("NAME=%s", name))
		args = append(args, name)
		conditions = append(conditions, "name = ?")
	}

	if len(conditions) == 0 {
		return []interface{}{}, ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(keys, ":")))
	return append([]interface{}{strings.Join(conditions, " AND ")}, args...), hex.EncodeToString(hasher.Sum(nil))
}

func (o *projectController) parseBody(body *models.ProjectDto, model *models.Project) bool {
	model.Name = body.Name
	model.Title = body.Title
	model.Desc = body.Desc
	model.Note = body.Note
	model.Flag = body.Flag

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

	if len(body.Links) != 0 {
		var fileMap = make(map[string]*models.Link)

		model.Links = make([]models.Link, len(body.Links))
		for i := 0; i < len(body.Links); i++ {
			if _, ok := fileMap[body.Links[i].Name]; ok {
				return false
			}

			o.parseLinkBody(&body.Links[i], &model.Links[i])
			fileMap[model.Links[i].Name] = &model.Links[i]
		}
	}
	return true
}

func (*projectController) parseFileBody(body *models.File, model *models.File) {
	model.Name = body.Name
	model.Path = body.Path
	model.Role = body.Role
	model.Type = body.Type
}

func (*projectController) parseLinkBody(body *models.Link, model *models.Link) {
	model.Name = body.Name
	model.Link = body.Link
}

// @Tags Project
// @Summary Create file by project id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param model body models.ProjectDto true "Project Data"
// @Success 201 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [post]
func (o *projectController) CreateOne(c *gin.Context) {
	var body models.ProjectDto
	if err := c.ShouldBind(&body); err != nil || body.Name == "" || body.Title == "" {
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
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	ctx := context.Background()
	db.Redis.Incr(ctx, "nPROJECT")
	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	go db.FlushValue("LINK")
	go db.FlushValue("FILE")
	go db.FlushValue("METRICS")
	go db.FlushValue("SUBSCRIPTION")
	go db.FlushValue("PROJECT")

	go db.Redis.Incr(ctx, "nPROJECT")
	helper.ResHandler(c, http.StatusCreated, &models.Success{
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
// @Param model body []models.ProjectDto true "List of Project Data"
// @Success 201 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/list/{id} [post]
func (o *projectController) CreateAll(c *gin.Context) {
	var body []models.ProjectDto
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	// TODO: Add query param to disable or enable loading Files or not

	var model = make([]models.Project, len(body))
	for i := 0; i < len(body); i++ {
		if body[i].Name == "" || body[i].Title == "" {
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
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	go db.FlushValue("LINK")
	go db.FlushValue("FILE")
	go db.FlushValue("METRICS")
	go db.FlushValue("SUBSCRIPTION")
	go db.FlushValue("PROJECT")

	go helper.RedisAdd(&ctx, "nPROJECT", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, &models.Success{
		Status:     "OK",
		Result:     model,
		Items:      result.RowsAffected,
		TotalItems: items + result.RowsAffected,
	})
}

// @Tags Project
// @Summary Read Project by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Param name path string true "Project Name"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/{name} [get]
func (*projectController) ReadOne(c *gin.Context) {
	var name string
	var model []models.Project

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("NAME=%s", name)))
	if err := helper.PrecacheResult(fmt.Sprintf("PROJECT:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where("name = ?", name).Preload("Files").Preload("Links").Preload("Metrics").Preload("Subscription"), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/project.go", err.Error())
		return
	}

	var items int64
	var err error
	if items, err = db.Redis.Get(context.Background(), "nPROJECT").Int64(); err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     model,
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
// @Param name query string false "Name: 'CodeRain'"
// @Param title query string false "Title: 'Code Rain'"
// @Param start query string false "CreatedAt date >= start"
// @Param end query string false "CreatedAt date <= end"
// @Param page query int false "Page: '0'"
// @Param limit query int false "Limit: '1'"
// @Param type query string false "Files Type: 'js,html,img'"
// @Param role query string false "Files Role: 'src,assets,styles'"
// @Param link_name query string false "Name: 'main'"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 429 {object} models.Error
// @failure 400 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [get]
func (o *projectController) ReadAll(c *gin.Context) {
	var model []models.Project
	ctx := context.Background()

	page, limit := helper.Pagination(c)

	result, suffix := o.filterQuery(c)
	fileCondition, fileSuffix := o.filterFileQuery(c)
	linkCondition, linkSufix := o.filterLinkQuery(c)

	if err := helper.PrecacheResult(fmt.Sprintf("PROJECT:%s:%s:%s:%d:%d", suffix, fileSuffix, linkSufix, page, limit), result.Offset(page*config.ENV.Items).Limit(limit).Preload("Files", fileCondition...).Preload("Links", linkCondition...), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/project.go", err.Error())
		return
	}

	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     model,
		Page:       page,
		Limit:      limit,
		Items:      int64(len(model)),
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Update Project by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Project Name"
// @Param model body models.ProjectDto true "Project without File Data"
// @Success 200 {object} models.Success{result=[]models.Project}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 416 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/{name} [put]
func (o *projectController) UpdateOne(c *gin.Context) {
	var body models.ProjectDto

	var name = c.Param("name")
	if err := c.ShouldBind(&body); err != nil || name == "" || len(body.Files) != 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	// TODO: Maybe I need to check if new dir is already exits
	// because right now it can be achieved by getting an error
	// from db
	var model models.Project
	o.parseBody(&body, &model)
	result := db.DB.Where("name = ?", name).Updates(&model)
	if result.RowsAffected != 1 {
		helper.ErrHandler(c, http.StatusRequestedRangeNotSatisfiable, "Such record doesn't exist within db")
		return
	}

	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	go db.FlushValue("LINK")
	go db.FlushValue("FILE")
	go db.FlushValue("METRICS")
	go db.FlushValue("SUBSCRIPTION")
	go db.FlushValue("PROJECT")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     []models.Project{model},
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
// @Param name query string false "Type: 'CodeRain'"
// @Param title query string false "Type: 'Code Rain'"
// @Param model body models.ProjectDto true "Project without File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [put]
func (o *projectController) UpdateAll(c *gin.Context) {
	var body models.ProjectDto
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
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	go db.FlushValue("LINK")
	go db.FlushValue("FILE")
	go db.FlushValue("METRICS")
	go db.FlushValue("SUBSCRIPTION")
	go db.FlushValue("PROJECT")

	var items int64
	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     []models.Project{model},
		Items:      result.RowsAffected,
		TotalItems: items,
	})
}

// @Tags Project
// @Summary Delete Project and Files by it's name
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param name path string true "Project Name"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 416 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project/{name} [delete]
func (*projectController) DeleteOne(c *gin.Context) {
	var name string
	var project []models.Project

	if name = c.Param("name"); name == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	key := "Project:" + name
	ctx := context.Background()
	// Check if cache have requested data
	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), &project)
		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		result := db.DB.Where("name = ?", name).Find(&project)
		if result.Error != nil {
			helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
			go logs.DefaultLog("/controllers/project.go", result.Error)
			return
		}

		// Encode json to str
		if str, err := json.Marshal(&project); err == nil {
			go db.Redis.Set(ctx, key, str, time.Duration(config.ENV.LiveTime)*time.Second)
		}
	}

	if len(project) == 0 {
		helper.ErrHandler(c, http.StatusRequestedRangeNotSatisfiable, "Such record doesn't exist within db")
		return
	}

	// Delete in both place Files & Project
	if result := db.DB.Where("project_id = ?", project[0].ID).Delete(&models.File{}); result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	result := db.DB.Where("name = ?", name).Delete(&models.Project{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/project.go", result.Error)
		return
	}

	db.Redis.Del(ctx, key)
	items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	if err != nil {
		items = -1
		go (&models.Project{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/project.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.FlushValue("LINK")
	go db.FlushValue("FILE")
	go db.FlushValue("METRICS")
	go db.FlushValue("SUBSCRIPTION")
	go db.FlushValue("PROJECT")

	go db.Redis.Decr(ctx, "nPROJECT")
	helper.ResHandler(c, http.StatusOK, &models.Success{
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
// @Param name query string false "Type: 'CodeRain'"
// @Param title query string false "Type: 'Code Rain'"
// @Success 200 {object} models.Success{result=[]string{}}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 416 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /project [delete]
func (o *projectController) DeleteAll(c *gin.Context) {
	// FIXME::
	// var sKeys string
	// var result *gorm.DB
	// var project []models.Project

	// if result, sKeys = o.filterQuery(c); sKeys == "" {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
	// 	return
	// }

	// ctx := context.Background()
	// key := "Project:" + c.DefaultQuery(sKeys, "-1")
	// if sKeys == "id" || sKeys == "name" || sKeys == "start" || sKeys == "end" {
	// 	// Check if cache have requested data
	// 	if data, err := db.Redis.Get(ctx, key).Result(); err == nil {
	// 		json.Unmarshal([]byte(data), &project)
	// 		go db.Redis.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	// 	}
	// }

	// if len(project) == 0 {
	// 	if result := result.Find(&project); result.Error != nil || len(project) == 0 {
	// 		helper.ErrHandler(c, http.StatusRequestedRangeNotSatisfiable, "Such record doesn't exist within db")
	// 		go logs.DefaultLog("/controllers/project.go", result.Error)
	// 		return
	// 	}
	// }

	// // Delete in both place Files & Project
	// if result := db.DB.Where("project_id = ?", project[0].ID).Delete(&models.File{}); result.Error != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
	// 	go logs.DefaultLog("/controllers/project.go", result.Error)
	// 	return
	// }

	// result = result.Delete(&models.Project{})
	// if result.Error != nil {
	// 	helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
	// 	go logs.DefaultLog("/controllers/project.go", result.Error)
	// 	return
	// }

	// if sKeys == "id" || sKeys == "name" || sKeys == "title" {
	// 	db.Redis.Del(ctx, "Project:"+c.DefaultQuery(sKeys, ""))
	// }

	// items, err := db.Redis.Get(ctx, "nPROJECT").Int64()
	// if err != nil {
	// 	items = -1
	// 	go (&models.Project{}).Redis(db.DB, db.Redis)
	// 	go logs.DefaultLog("/controllers/project.go", err.Error())
	// }

	// if items == 0 {
	// 	helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
	// 	return
	// }

	// go db.FlushValue("File")
	// go db.FlushValue("Project")

	// go helper.RedisSub(&ctx, "nProject", result.RowsAffected)
	// helper.ResHandler(c, http.StatusOK, &models.Success{
	// 	Status:     "OK",
	// 	Result:     []string{},
	// 	Items:      result.RowsAffected,
	// 	TotalItems: items - result.RowsAffected,
	// })
}
