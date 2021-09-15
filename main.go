package main

import (
	"api/config"
	"api/db"
	"api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv("./")

	db.ConnectToDB()
	db.MigrateTables(config.ENV.ForceMigrate)

	db.ConnectToRedis()
	db.RedisInitDefault()

	r := gin.Default()
	routes.Init(r)

	r.Run(config.ENV.Host + ":" + config.ENV.Port)
}
