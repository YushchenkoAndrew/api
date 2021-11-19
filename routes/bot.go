package routes

import (
	c "api/controllers"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Bot(rg *gin.RouterGroup) {
	cBot := c.BotController{}
	auth := rg.Group("/bot", middleware.Auth())

	auth.POST("/redis", cBot.Redis)
}
