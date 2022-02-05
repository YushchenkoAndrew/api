package routes

import (
	c "api/controllers"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type linkRouter struct {
	route *gin.RouterGroup
	auth  *gin.RouterGroup
	link  interfaces.Default
}

func NewLinkRouter(rg *gin.RouterGroup) interfaces.Router {
	return &linkRouter{
		route: rg.Group(("/link")),
		auth:  rg.Group("/link", middleware.Auth()),
		link:  c.NewLinkController(),
	}
}

func (c *linkRouter) Init() {
	c.auth.POST("/list/:id", c.link.CreateAll)
	c.auth.POST("/:id", c.link.CreateOne)

	c.route.GET("/:id", c.link.ReadOne)
	c.route.GET("", c.link.ReadAll)

	c.auth.PUT("/:id", c.link.UpdateOne)
	c.auth.PUT("", c.link.UpdateAll)

	c.auth.DELETE("/:id", c.link.DeleteOne)
	c.auth.DELETE("", c.link.DeleteAll)
}
