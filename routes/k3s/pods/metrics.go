package pods

import (
	c "api/controllers/k3s/pods"
	"api/interfaces"
	"api/interfaces/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type metricsRouter struct {
	auth    *gin.RouterGroup
	metrics k3s.Metrics
}

func NewMetricsRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &metricsRouter{
			auth:    rg.Group("/metrics", middleware.Auth()),
			metrics: c.NewMetricsController(),
		}
	}
}

func (c *metricsRouter) Init() {
	c.auth.GET("", c.metrics.ReadAll)
	c.auth.GET("/:name", c.metrics.ReadOne)
}
