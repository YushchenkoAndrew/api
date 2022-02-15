package routes

import (
	c "api/controllers"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type fileRouter struct {
	route *gin.RouterGroup
	auth  *gin.RouterGroup
	file  interfaces.Default
}

func NewFileRouter(rg *gin.RouterGroup) interfaces.Router {
	return &fileRouter{
		route: rg.Group("/file"),
		auth:  rg.Group("/file", middleware.Auth()),
		file:  c.NewFileController(),
	}
}

func (c *fileRouter) Init() {
	c.auth.POST("/list/:id", c.file.CreateAll)
	c.auth.POST("/:id", c.file.CreateOne)

	c.route.GET("/:id", c.file.ReadOne)
	c.route.GET("", c.file.ReadAll)

	c.auth.PUT("/:id", c.file.UpdateOne)
	c.auth.PUT("", c.file.UpdateAll)

	c.auth.DELETE("/:id", c.file.DeleteOne)
	c.auth.DELETE("", c.file.DeleteAll)
}
