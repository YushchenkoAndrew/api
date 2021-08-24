package controllers

import (
	"api/db"
	"api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InfoController struct{}

func (*InfoController) Create(c *gin.Context) {
	db.DB.Model(&models.Info{}).Create(map[string]interface{}{
		"countries": "UA",
	})
}

func (*InfoController) ReadAll(c *gin.Context) {
	var info models.Info
	db.DB.Find(&info)

	// FIXME:
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"result": info,
	})
}

func (*InfoController) ReadOne(c *gin.Context) {

}

func (*InfoController) Update(c *gin.Context) {

}

func (*InfoController) Delete(c *gin.Context) {

}
