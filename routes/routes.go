package routes

import (
	"api/config"
	info "api/routes/info"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

// @host mortis-grimreaper.ddns.net:31337
// @BasePath /api
func Init(rg *gin.Engine) {
	route := rg.Group(config.ENV.BasePath)

	Index(route)
	Info(route)
	World(route)

	// Init SubRoutes
	info.Init(route)

	// Init Route by hand
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
