package routes

import (
	c "api/controllers"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type k3sRouter struct {
	auth      *gin.RouterGroup
	k3s       interfaces.K3s
	subRoutes []interfaces.Router
}

func NewK3sRouter(rg *gin.RouterGroup, handlers []func(*gin.RouterGroup) interfaces.Router) interfaces.Router {
	route := rg.Group("/k3s")
	var subRoutes []interfaces.Router
	for _, handler := range handlers {
		subRoutes = append(subRoutes, handler(route))
	}

	return &k3sRouter{
		auth:      rg.Group("/k3s", middleware.Auth()),
		k3s:       c.NewK3sController(),
		subRoutes: subRoutes,
	}
}

func (c *k3sRouter) Init() {
	c.auth.POST("/subscribe/:id", c.k3s.Subscribe)
	c.auth.DELETE("/subscribe/:id", c.k3s.Unsubscribe)

	for _, route := range c.subRoutes {
		route.Init()
	}
}
