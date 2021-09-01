package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func World(rg *gin.RouterGroup) {
	route := rg.Group("/world")
	auth := rg.Group("/world", middleware.Auth())
	cWorld := c.WorldController{}

	auth.POST("", cWorld.CreateOne)
	auth.POST("/list", cWorld.CreateAll)

	route.GET("/:id", cWorld.ReadOne)
	route.GET("", cWorld.ReadAll)

	auth.PUT("/:id", cWorld.UpdateOne)
	auth.PUT("", cWorld.UpdateAll)

	auth.DELETE("/:id", cWorld.DeleteOne)
	auth.DELETE("", cWorld.DeleteAll)
}
