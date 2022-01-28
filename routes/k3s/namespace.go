package k3s

import (
	c "api/controllers/k3s"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Namespace(rg *gin.RouterGroup) {
	auth := rg.Group("/namespace", middleware.Auth())
	cNamespace := c.NamespaceController{}

	auth.POST("", cNamespace.Create)

	auth.GET("", cNamespace.ReadAll)
	auth.GET("/:name", cNamespace.ReadOne)
}
