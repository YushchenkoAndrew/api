package main

import (
	"api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api := r.Group("/api")

	routes.Init(api)

	r.Run(":31337")
}
