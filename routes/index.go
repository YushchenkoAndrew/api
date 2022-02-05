package routes

import (
	c "api/controllers"
	"api/interfaces"

	"github.com/gin-gonic/gin"
)

type indexRouter struct {
	route     *gin.RouterGroup
	index     interfaces.Index
	subRoutes *[]interfaces.Router
}

func NewIndexRouter(route *gin.RouterGroup, subRoutes *[]interfaces.Router) interfaces.Router {
	return &indexRouter{route: route, index: c.NewIndexController(), subRoutes: subRoutes}
}

func (c *indexRouter) Init() {
	c.route.GET("/ping", c.index.Ping)
	c.route.GET("/trace/:ip", c.index.TraceIp)
	c.route.POST("/login", c.index.Login)
	c.route.POST("/refresh", c.index.Refresh)

	for _, route := range *c.subRoutes {
		route.Init()
	}
}
