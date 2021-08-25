package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func World(rg *gin.RouterGroup) {
	route := rg.Group("/world")
	cWorld := c.WorldController{}

	route.POST("", cWorld.Create)
	route.GET("", cWorld.Read)
	route.PUT("", cWorld.Update)
	route.DELETE("", cWorld.Delete)
}
