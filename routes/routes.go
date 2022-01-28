package routes

import (
	"api/config"
	_ "api/docs"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

func Init(rg *gin.Engine) {
	route := rg.Group(config.ENV.BasePath, middleware.Limit())

	Index(route)
	Info(route)
	World(route)
	Project(route)
	File(route)
	Link(route)
	Bot(route)
	K3s(route)
}
