package controllers

import (
	"api/db"
	"api/helper"
	"api/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var info []models.Info

type InfoController struct{}

func (*InfoController) filterQuery(c *gin.Context) (*gorm.DB, bool) {
	bInUse := false
	result := db.DB
	id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
	if err == nil && id > 0 {
		bInUse = bInUse || true
		result = result.Where("id = ?", id)
	}

	createdAt := c.DefaultQuery("created_at", "")
	if createdAt != "" {
		bInUse = bInUse || true
		result = result.Where("created_at = ?", createdAt)
	}

	countries := c.DefaultQuery("countries", "")
	if countries != "" {
		bInUse = bInUse || true
		result = result.Where("countries IN ?", strings.Split(countries, ","))
	}

	return result.Find(&info), bInUse
}

func (*InfoController) parseBody(body *models.ReqInfo, model *models.Info) {
	model.Countries = body.Countries

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
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params or id parm")
		return
	}

	o.parseBody(&body, &model)
	result := db.DB.Create(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": 20 + 1,
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
		o.parseBody(&body[i], &models[i])
	}

	result := db.DB.Create(&models)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusCreated, &gin.H{
		"status":        "OK",
		"resHandlerult": models,
		"items":         result.RowsAffected,
		"totalItems":    20 + 1,
	})
}

func (o *InfoController) ReadOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	result := db.DB.Where("id = ?", id).Find(&info)
	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     info,
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) ReadAll(c *gin.Context) {
	page, limit := helper.Pagination(c)
	result, _ := o.filterQuery(c)
	result = result.Offset(page * 20).Limit(limit)

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     info,
		"page":       page,
		"limit":      limit,
		"items":      result.RowsAffected,
		"totalItems": 20,
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
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) UpdateAll(c *gin.Context) {
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var bInUse bool
	var result *gorm.DB
	var model models.Info

	o.parseBody(&body, &model)
	if result, bInUse = o.filterQuery(c); !bInUse {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Updates(&model)
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) DeleteOne(c *gin.Context) {
	var id int
	if !helper.GetID(c, &id) {
		helper.ErrHandler(c, http.StatusBadRequest, "Incorrect id params")
		return
	}

	result := db.DB.Where("id = ?", id).Delete(&models.Info{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) DeleteAll(c *gin.Context) {
	var bInUse bool
	var result *gorm.DB

	if result, bInUse = o.filterQuery(c); !bInUse {
		helper.ErrHandler(c, http.StatusBadRequest, "Query not founded")
		return
	}

	result = result.Delete(&models.Info{})
	if result.Error != nil {
		helper.ErrHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	helper.ResHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     []string{},
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}
