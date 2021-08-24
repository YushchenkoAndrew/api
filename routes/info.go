package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func (*Routes) info(rg *gin.RouterGroup) {
	route := rg.Group("/info")
	cInfo := new(c.InfoController)

	route.POST("", cInfo.Create)

	route.GET("", cInfo.ReadAll)
	route.GET("/:id", cInfo.ReadOne)

	route.PUT("", cInfo.Update)
	route.DELETE("", cInfo.Delete)
}
