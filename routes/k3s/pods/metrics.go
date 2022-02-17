package pods

import (
	c "api/controllers/k3s/pods"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type metricsRouter struct {
	auth      *gin.RouterGroup
	authToken *gin.RouterGroup
	metrics   interfaces.Default
}

func NewMetricsRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &metricsRouter{
			auth:      rg.Group("/metrics", middleware.Auth()),
			authToken: rg.Group("/metrics", middleware.AuthToken()),
			metrics:   c.NewMetricsController(),
		}
	}
}

func (c *metricsRouter) Init() {
	c.auth.GET("", c.metrics.ReadAll)
	c.auth.GET("/:id", c.metrics.ReadOne)
	// c.auth.GET("/:id/:namespace/:name", c.metrics.ReadOne)

	c.authToken.POST("/:id/:namespace", c.metrics.CreateAll)
	c.authToken.POST("/:id/:namespace/:name", c.metrics.CreateOne)
}
