package routes

import (
	// conf "api/config"
	index "api/controllers"

	"github.com/gin-gonic/gin"
)

var ctrIndex = new(index.Controller)

func Init(rg *gin.RouterGroup) {

	rg.GET("/ping", ctrIndex.Ping)
}
