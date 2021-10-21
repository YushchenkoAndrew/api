package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Project(rg *gin.RouterGroup) {
	route := rg.Group("/project")
	auth := rg.Group("/project", middleware.Auth())
	cProject := c.ProjectController{}

	auth.POST("", cProject.CreateOne)
	auth.POST("/list", cProject.CreateAll)

	route.GET("/:name", cProject.ReadOne)
	route.GET("", cProject.ReadAll)

	auth.PUT("/:name", cProject.UpdateOne)
	auth.PUT("", cProject.UpdateAll)

	auth.DELETE("/:name", cProject.DeleteOne)
	auth.DELETE("", cProject.DeleteAll)
}
