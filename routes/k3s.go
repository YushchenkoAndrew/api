package routes

import (
	c "api/controllers"
	"api/middleware"
	sub "api/routes/k3s"

	"github.com/gin-gonic/gin"
)

func K3s(rg *gin.RouterGroup) {
	route := rg.Group("/k3s")
	auth := rg.Group("/k3s", middleware.Auth())
	cK3s := c.K3sController{}

	auth.POST("/subscribe/:id", cK3s.Subscribe)
	auth.DELETE("/subscribe/:id", cK3s.Unsubscribe)

	sub.Deployment(route)
}
