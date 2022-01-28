package k3s

import (
	c "api/controllers/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Deployment(rg *gin.RouterGroup) {
	auth := rg.Group("/deployment", middleware.Auth())
	cDeployment := c.DeploymentController{}

	auth.POST("/:namespace", cDeployment.Create)

	auth.GET("", cDeployment.ReadAll)
	auth.GET("/:name", cDeployment.ReadOne)
}
