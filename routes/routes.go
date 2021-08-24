package routes

import (
	"api/config"

	"github.com/gin-gonic/gin"
)

type Routes struct{}

func (r *Routes) Init(rg *gin.Engine) {
	route := rg.Group(config.ENV.BasePath)

	r.index(route)
	r.info(route)
	r.world(route)
}
