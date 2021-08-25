package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func (*Routes) info(rg *gin.RouterGroup) {
	route := rg.Group("/info")
	cInfo := c.InfoController{}

	route.POST("", cInfo.CreateOne)
	route.POST("/list", cInfo.CreateAll)

	route.GET("", cInfo.ReadAll)
	route.GET("/:id", cInfo.ReadOne)

	route.PUT("/:id", cInfo.UpdateOne)
	route.PUT("", cInfo.UpdateAll)

	route.DELETE("", cInfo.Delete)
}
