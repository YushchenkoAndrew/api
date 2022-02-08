package k3s

import (
	c "api/controllers/k3s"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type deploymentRouter struct {
	auth       *gin.RouterGroup
	deployment interfaces.Default
}

func NewDeploymentRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &deploymentRouter{
			auth:       rg.Group("/deployment", middleware.Auth()),
			deployment: c.NewDeploymentController(),
		}
	}
}

func (c *deploymentRouter) Init() {
	c.auth.POST("/:namespace", c.deployment.CreateOne)

	c.auth.GET("", c.deployment.ReadAll)
	c.auth.GET("/:namespace", c.deployment.ReadAll)
	c.auth.GET("/:namespace/:name", c.deployment.ReadOne)
}
