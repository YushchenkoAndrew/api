package routes

import (
	index "api/controllers"

	"github.com/gin-gonic/gin"
)

var ctrIndex = new(index.Controller)

func Init(rg *gin.RouterGroup) {

	rg.GET("/ping", ctrIndex.Ping)
}
