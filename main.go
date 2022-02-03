package main

import (
	"api/config"
	"api/db"
	"api/routes"

	"github.com/gin-gonic/gin"
)

// @title Gin API
// @version 1.0
// @description Remake of my previous attampted on creating API with Node.js

// @contact.name API Author
// @contact.url https://mortis-grimreaper.ddns.net/projects
// @contact.email AndrewYushchenko@gmail.com

// @license.name MIT
// @license.url https://github.com/YushchenkoAndrew/API_Server/blob/master/LICENSE

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host mortis-grimreaper.ddns.net:31337
// @BasePath /api
func main() {
	config.Init()

	db.ConnectToDB()
	db.MigrateTables()

	db.ConnectToRedis()
	db.RedisInitDefault()

	r := gin.Default()
	routes.Init(r)

	r.Run(config.ENV.Host + ":" + config.ENV.Port)
}
