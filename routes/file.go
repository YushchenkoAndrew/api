package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func File(rg *gin.RouterGroup) {
	route := rg.Group("/file")
	auth := rg.Group("/file", middleware.Auth())
	cFile := c.FileController{}

	auth.POST("/list/:id", cFile.CreateAll)
	auth.POST("/:id", cFile.CreateOne)

	route.GET("/:id", cFile.ReadOne)
	route.GET("", cFile.ReadAll)

	auth.PUT("/:id", cFile.UpdateOne)
	auth.PUT("", cFile.UpdateAll)

	auth.DELETE("/:id", cFile.DeleteOne)
	auth.DELETE("", cFile.DeleteAll)
}
