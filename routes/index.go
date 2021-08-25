package routes

import (
	c "api/controllers"

	"github.com/gin-gonic/gin"
)

func Index(rg *gin.RouterGroup) {
	cIndex := c.IndexController{}

	rg.GET("", cIndex.Navigation)
	rg.GET("/ping", cIndex.Ping)
	rg.GET("/tables", cIndex.TableList)
}
