package k3s

import (
	c "api/controllers/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Ingress(rg *gin.RouterGroup) {
	auth := rg.Group("/ingress", middleware.Auth())
	cIngress := c.IngressController{}

	auth.POST("/:namespace", cIngress.Create)

	auth.GET("/", cIngress.ReadAll)
	auth.GET("/:name", cIngress.ReadOne)
}
