package routes

import (
	"api/config"
	"api/middleware"

	_ "api/docs"

	"github.com/gin-gonic/gin"
)

// @title Gin API
// @version 1.0
// @description Remake of my previous attampted on creating API with Node.js

// @contact.name API Support
// @contact.url https://mortis-grimreaper.ddns.net/projects
// @contact.email AndrewYushchenko@gmail.com

// @license.name MIT
// @license.url https://github.com/YushchenkoAndrew/API_Server/blob/master/LICENSE

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

//  FIXME: DEBUG OPTION
// host mortis-grimreaper.ddns.net:31337
// @host 127.0.0.1:31337
// @BasePath /api
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
