package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func Info(rg *gin.RouterGroup) {
	route := rg.Group("/info")
	cInfo := c.InfoController{}

	// @Success 200 {string} string	"ok"
	// @failure 400 {string} string	"error"
	// @response default {string} string	"other error"
	// @Header 200 {string} Location "/entity/1"
	// @Header 200,400,default {string} Token "token"
	// @Header all {string} Token2 "token2"
	route.POST("", cInfo.CreateOne)

	route.POST("/list", cInfo.CreateAll)

	route.GET("", cInfo.ReadAll)
	route.GET("/:id", cInfo.ReadOne)

	route.PUT("", cInfo.UpdateAll)
	route.PUT("/:id", cInfo.UpdateOne)

	route.DELETE("", cInfo.DeleteAll)
	route.DELETE("/:id", cInfo.DeleteOne)
}
