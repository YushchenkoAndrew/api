package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func World(rg *gin.RouterGroup) {
	route := rg.Group("/world")
	cWorld := c.WorldController{}

	route.POST("", cWorld.CreateOne)
	route.POST("/list", cWorld.CreateAll)

	route.GET("/:id", cWorld.ReadOne)
	route.GET("", cWorld.ReadAll)

	route.PUT("/:id", cWorld.UpdateOne)
	route.PUT("", cWorld.UpdateAll)

	route.DELETE("/:id", cWorld.DeleteOne)
	route.DELETE("", cWorld.DeleteAll)
}
