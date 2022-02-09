package routes

import (
	c "api/controllers"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type infoRouter struct {
	route     *gin.RouterGroup
	auth      *gin.RouterGroup
	info      interfaces.Info
	subRoutes []interfaces.Router
}

func NewInfoRouter(rg *gin.RouterGroup, handlers []func(*gin.RouterGroup) interfaces.Router) interfaces.Router {
	var subRoutes []interfaces.Router
	for _, handler := range handlers {
		subRoutes = append(subRoutes, handler(rg))
	}

	return &infoRouter{
		route:     rg.Group(("/info")),
		auth:      rg.Group("/info", middleware.Auth()),
		info:      c.NewInfoController(),
		subRoutes: subRoutes,
	}
}

func (c *infoRouter) Init() {
	c.auth.POST("", c.info.Create)
	c.auth.POST("/list", c.info.CreateAll)
	c.auth.POST("/:date", c.info.CreateOne)

	c.route.GET("", c.info.ReadAll)
	c.route.GET("/:id", c.info.ReadOne)

	c.auth.PUT("", c.info.UpdateAll)
	c.auth.PUT("/:id", c.info.UpdateOne)

	c.auth.DELETE("", c.info.DeleteAll)
	c.auth.DELETE("/:id", c.info.DeleteOne)

	for _, route := range c.subRoutes {
		route.Init()
	}
}
