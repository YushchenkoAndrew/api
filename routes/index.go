package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Index(rg *gin.RouterGroup) {
	cIndex := c.IndexController{}

	rg.GET("/ping", cIndex.Ping)
	rg.GET("/trace/:ip", cIndex.TraceIp)
	rg.POST("/login", cIndex.Login)
	rg.POST("/refresh", cIndex.Refresh)

	rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
