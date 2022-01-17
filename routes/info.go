package routes

import (
	c "api/controllers"
	"api/middleware"
	sub "api/routes/info"

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

	auth.PUT("", cInfo.UpdateAll)
	auth.PUT("/:id", cInfo.UpdateOne)

	auth.DELETE("", cInfo.DeleteAll)
	auth.DELETE("/:id", cInfo.DeleteOne)

	// Sub routes
	sub.Sum(route)
	sub.Range(route)
}
