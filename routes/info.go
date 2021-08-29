package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func Info(rg *gin.RouterGroup) {
	route := rg.Group("/info")
	cInfo := c.InfoController{}

	route.POST("", cInfo.Create)
	route.POST("/list", cInfo.CreateAll)
	route.POST("/:date", cInfo.CreateOne)

	route.GET("", cInfo.ReadAll)
	route.GET("/:id", cInfo.ReadOne)

	route.PUT("", cInfo.UpdateAll)
	route.PUT("/:id", cInfo.UpdateOne)

	route.DELETE("", cInfo.DeleteAll)
	route.DELETE("/:id", cInfo.DeleteOne)
}
