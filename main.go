package main

import (
	utils "api/config"
	"api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	env := new(utils.Config).Init(".")

	db := new(utils.Database)
	db.Init(&env)
	db.Migrate(true)

	redis := new(utils.Redis)
	redis.Init(&env)
	redis.Test()

	r := gin.Default()
	api := r.Group(env.BasePath)

	routes.Init(api)

	r.Run(env.Host + ":" + env.Port)
}
