package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type namespaceRouter struct {
	auth      *gin.RouterGroup
	namespace interfaces.Default
}

func NewNamespaceRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &namespaceRouter{
			auth:      rg.Group("/namespace", middleware.Auth()),
			namespace: c.NewNamespaceController(),
		}
	}
}

func (c *namespaceRouter) Init() {
	c.auth.POST("", c.namespace.CreateOne)

	c.auth.GET("", c.namespace.ReadAll)
	c.auth.GET("/:name", c.namespace.ReadOne)
}
