package main

import (
	"api/config"
	"api/db"
	"api/routes"

	"github.com/gin-gonic/gin"
)

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
