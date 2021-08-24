package main

import (
	"api/config"
	"api/db"
	r "api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv(".")

	db.ConnectToDB()
	db.MigrateTables(true)

	db.ConnectToRedis()
	db.TestRedis()

	routes := new(r.Routes)

	api := gin.Default()
	routes.Init(api)

	api.Run(config.ENV.Host + ":" + config.ENV.Port)
}
