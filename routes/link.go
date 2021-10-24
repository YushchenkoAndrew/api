package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Link(rg *gin.RouterGroup) {
	route := rg.Group("/link")
	auth := rg.Group("/link", middleware.Auth())
	cLink := c.LinkController{}

	auth.POST("/list/:id", cLink.CreateAll)
	auth.POST("/:id", cLink.CreateOne)

	route.GET("/:id", cLink.ReadOne)
	route.GET("", cLink.ReadAll)

	auth.PUT("/:id", cLink.UpdateOne)
	auth.PUT("", cLink.UpdateAll)

	auth.DELETE("/:id", cLink.DeleteOne)
	auth.DELETE("", cLink.DeleteAll)
}
