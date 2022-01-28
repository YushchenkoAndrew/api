package k3s

import (
	c "api/controllers/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Pods(rg *gin.RouterGroup) {
	auth := rg.Group("/pods", middleware.Auth())
	cPods := c.PodsController{}

	auth.POST("/:name", cPods.Exec)

	auth.GET("", cPods.ReadAll)
	auth.GET("/:name", cPods.ReadOne)

	// Subroutes 'metrics'
	auth.GET("/metrics", cPods.ReadMetricsAll)
	auth.GET("/metrics/:name", cPods.ReadMetricsOne)
}
