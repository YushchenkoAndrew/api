package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Info(rg *gin.RouterGroup) {
	route := rg.Group("/info")
	auth := rg.Group("/info", middleware.Auth())
	cInfo := c.InfoController{}

	auth.POST("", cInfo.Create)
	auth.POST("/list", cInfo.CreateAll)
	auth.POST("/:date", cInfo.CreateOne)

	route.GET("", cInfo.ReadAll)
	route.GET("/:id", cInfo.ReadOne)

	route.PUT("", cInfo.UpdateAll)
	route.PUT("/:id", cInfo.UpdateOne)

	route.DELETE("", cInfo.DeleteAll)
	route.DELETE("/:id", cInfo.DeleteOne)
}
