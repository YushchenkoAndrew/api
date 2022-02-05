package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type serviceRouter struct {
	auth    *gin.RouterGroup
	service interfaces.Default
}

func NewServiceRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &serviceRouter{
			auth:    rg.Group("/service", middleware.Auth()),
			service: c.NewServiceController(),
		}
	}
}

func (c *serviceRouter) Init() {
	c.auth.POST("/:namespace", c.service.CreateOne)

	c.auth.GET("", c.service.ReadAll)
	c.auth.GET("/:name", c.service.ReadOne)
}
