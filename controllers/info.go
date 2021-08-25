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

type InfoController struct {
	*Controller
}

func (*InfoController) filterQuery(c *gin.Context) *gorm.DB {
	resHandlerult := db.DB
	id, err := strconv.Atoi(c.DefaultQuery("id", "-1"))
	if err == nil && id > 0 {
		resHandlerult = resHandlerult.Where("id = ?", id)
	}

	createdAt := c.DefaultQuery("created_at", "")
	if createdAt != "" {
		resHandlerult = resHandlerult.Where("created_at = ?", createdAt)
	}

	countries := c.DefaultQuery("country", "")
	if countries != "" {
		resHandlerult = resHandlerult.Where("countries IN ?", strings.Split(countries, ","))
	}

	return resHandlerult.Find(&info)
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
		o.errHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	o.parseBody(&body, &model)
	result := db.DB.Create(&model)
	if result.Error != nil {
		o.errHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	o.resHandler(c, http.StatusCreated, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": 20 + 1,
	})
}

func (o *InfoController) CreateAll(c *gin.Context) {
	var body []models.ReqInfo
	if err := c.ShouldBind(&body); err != nil {
		o.errHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var models = make([]models.Info, len(body))
	for i := 0; i < len(body); i++ {
		o.parseBody(&body[i], &models[i])
	}

	resHandlerult := db.DB.Create(&models)
	if resHandlerult.Error != nil {
		o.errHandler(c, http.StatusInternalServerError, "Server side error: Something went wrong")
		return
	}

	o.resHandler(c, http.StatusCreated, &gin.H{
		"status":        "OK",
		"resHandlerult": models,
		"items":         resHandlerult.RowsAffected,
		"totalItems":    20 + 1,
	})
}

func (o *InfoController) ReadOne(c *gin.Context) {
	var id int
	if !o.getID(c, &id) {
		o.errHandler(c, http.StatusBadRequest, "Incorrect id value")
		return
	}

	resHandlerult := db.DB.Where("id = ?", id).Find(&info)
	o.resHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     info,
		"items":      resHandlerult.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) ReadAll(c *gin.Context) {
	page, limit := helper.Pagination(c)
	resHandlerult := o.filterQuery(c).Offset(page * 20).Limit(limit)

	o.resHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     info,
		"page":       page,
		"limit":      limit,
		"items":      resHandlerult.RowsAffected,
		"totalItems": 20,
	})
}

func (o *InfoController) UpdateOne(c *gin.Context) {
	var id int
	var body models.ReqInfo
	if err := c.ShouldBind(&body); err != nil || !o.getID(c, &id) {
		o.errHandler(c, http.StatusBadRequest, "Incorrect body params")
		return
	}

	var model models.Info
	o.parseBody(&body, &model)
	result := db.DB.Where("id = ?", id).Updates(&model)

	o.resHandler(c, http.StatusOK, &gin.H{
		"status":     "OK",
		"result":     model,
		"items":      result.RowsAffected,
		"totalItems": 20,
	})
}

func (*InfoController) UpdateAll(c *gin.Context) {

}

func (*InfoController) Delete(c *gin.Context) {

}
