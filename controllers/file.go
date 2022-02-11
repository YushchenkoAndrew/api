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

type fileController struct{}

func NewFileController() interfaces.Default {
	return &fileController{}
}

func (*fileController) filterQuery(c *gin.Context) (*gorm.DB, string) {
	var keys = []string{}
	client := db.DB

	if id, err := strconv.Atoi(c.DefaultQuery("id", "-1")); err == nil && id > 0 {
		keys = append(keys, fmt.Sprintf("ID=%d", id))
		client = client.Where("id = ?", id)
	}

	if projectId, err := strconv.Atoi(c.DefaultQuery("project_id", "-1")); err == nil && projectId > 0 {
		keys = append(keys, fmt.Sprintf("PROJECT_ID =%d", projectId))
		client = client.Where("project_id = ?", projectId)
	}

	if name := c.DefaultQuery("name", ""); name != "" {
		keys = append(keys, fmt.Sprintf("NAME=%s", name))
		client = client.Where("name = ?", name)
	}

	if typeName := c.DefaultQuery("type", ""); typeName != "" {
		keys = append(keys, fmt.Sprintf("TYPE=%s", typeName))
		client = client.Where("type IN ?", strings.Split(typeName, ","))
	}

	if role := c.DefaultQuery("role", ""); role != "" {
		keys = append(keys, fmt.Sprintf("ROLE=%s", role))
		client = client.Where("role = ?", role)
	}

	if len(keys) == 0 {
		return client, ""
	}

	hasher := md5.New()
	hasher.Write([]byte(strings.Join(keys, ":")))
	return client, hex.EncodeToString(hasher.Sum(nil))
}

func (*fileController) parseBody(body *models.FileDto, model *models.File) {
	model.Name = body.Name
	model.Path = body.Path
	model.Role = body.Role
	model.Type = body.Type
}

func (*fileController) isExist(id int, body *models.FileDto) bool {
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
// @Param model body models.FileDto true "File Data"
// @Success 201 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/{id} [post]
func (o *fileController) CreateOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	var body models.FileDto
	if err := c.ShouldBind(&body); err != nil || body.Name == "" || body.Type == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body format")
		return
	}

	var result = db.DB
	var model []models.File

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%d", id)))
	var key = fmt.Sprintf("PROJECT:%s", hex.EncodeToString(hasher.Sum(nil)))

	// Check if such project is exist
	ctx := context.Background()
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
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	db.Redis.Incr(ctx, "nFILE")
	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	go db.Redis.Incr(ctx, "nFIlE")
	helper.ResHandler(c, http.StatusCreated, &models.Success{
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
// @Param model body []models.FileDto true "List of File Data"
// @Success 201 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/list/{id} [post]
func (o *fileController) CreateAll(c *gin.Context) {
	var err error
	var id int
	var body []models.FileDto

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect project id param")
		return
	}

	if err = c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%d", id)))
	var key = fmt.Sprintf("PROJECT:%s", hex.EncodeToString(hasher.Sum(nil)))

	ctx := context.Background()
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
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	go helper.RedisAdd(&ctx, "nFILE", result.RowsAffected)
	helper.ResHandler(c, http.StatusCreated, &models.Success{
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
func (*fileController) ReadOne(c *gin.Context) {
	var id int
	var model []models.File

	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("ID=%d", id)))
	if err := helper.PrecacheResult(fmt.Sprintf("FILE:%s", hex.EncodeToString(hasher.Sum(nil))), db.DB.Where("id = ?", id), &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/file.go", err.Error())
		return
	}

	var err error
	var items int64
	if items, err = db.Redis.Get(context.Background(), "nFILE").Int64(); err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     model,
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
// @Param name query string false "Type: 'index.js'"
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
func (o *fileController) ReadAll(c *gin.Context) {
	var model []models.File
	page, limit := helper.Pagination(c)

	client, suffix := o.filterQuery(c)
	if err := helper.PrecacheResult(fmt.Sprintf("FILE:%s", suffix), client, &model); err != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
		go logs.DefaultLog("/controllers/file.go", err.Error())
		return
	}

	items, err := db.Redis.Get(context.Background(), "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
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

// @Tags File
// @Summary Update File by :id
// @Accept json
// @Produce application/json
// @Produce application/xml
// @Security BearerAuth
// @Param id path int true "Instance id"
// @Param model body models.FileDto true "File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file/{id} [put]
func (o *fileController) UpdateOne(c *gin.Context) {
	var id int
	var body models.FileDto
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
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     []models.File{model},
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
// @Param model body models.FileDto true "File Data"
// @Success 200 {object} models.Success{result=[]models.File}
// @failure 400 {object} models.Error
// @failure 401 {object} models.Error
// @failure 422 {object} models.Error
// @failure 429 {object} models.Error
// @failure 500 {object} models.Error
// @Router /file [put]
func (o *fileController) UpdateAll(c *gin.Context) {
	var body models.FileDto
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
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	var items int64
	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     []models.File{model},
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
func (*fileController) DeleteOne(c *gin.Context) {
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
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go db.Redis.Decr(ctx, "nFILE")
	helper.ResHandler(c, http.StatusOK, &models.Success{
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
func (o *fileController) DeleteAll(c *gin.Context) {
	var sKeys string
	var result *gorm.DB

	if result, sKeys = o.filterQuery(c); sKeys == "" {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.File{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		go logs.DefaultLog("/controllers/file.go", result.Error)
		return
	}

	go db.FlushValue("FILE")

	ctx := context.Background()
	items, err := db.Redis.Get(ctx, "nFILE").Int64()
	if err != nil {
		items = -1
		go (&models.File{}).Redis(db.DB, db.Redis)
		go logs.DefaultLog("/controllers/file.go", err.Error())
	}

	if items == 0 {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect request")
		return
	}

	go helper.RedisSub(&ctx, "nFILE", result.RowsAffected)
	helper.ResHandler(c, http.StatusOK, &models.Success{
		Status:     "OK",
		Result:     []string{},
		Items:      result.RowsAffected,
		TotalItems: items - result.RowsAffected,
	})
}
