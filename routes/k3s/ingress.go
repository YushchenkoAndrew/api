package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type ingressRouter struct {
	auth    *gin.RouterGroup
	ingress interfaces.Default
}

func NewIngressRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &ingressRouter{
			auth:    rg.Group("/ingress", middleware.Auth()),
			ingress: c.NewIngressController(),
		}
	}
}

func (c *ingressRouter) Init() {
	c.auth.POST("/:namespace", c.ingress.CreateOne)

	c.auth.GET("", c.ingress.ReadAll)
	c.auth.GET("/:name", c.ingress.ReadOne)
}
