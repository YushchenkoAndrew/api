package k3s

import (
	c "api/controllers/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Service(rg *gin.RouterGroup) {
	auth := rg.Group("/service", middleware.Auth())
	cService := c.ServiceController{}

	auth.POST("/:namespace", cService.Create)

	auth.GET("/", cService.ReadAll)
	auth.GET("/:name", cService.ReadOne)
}
