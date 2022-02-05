package info

import (
	c "api/controllers/info"
	"api/interfaces"
	"api/interfaces/info"

	"github.com/gin-gonic/gin"
)

type sumRouter struct {
	route *gin.RouterGroup
	sum   info.Default
}

func NewSumRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &sumRouter{route: rg.Group("/sum"), sum: c.NewSumController()}
	}
}

func (c *sumRouter) Init() {
	c.route.GET("", c.sum.Read)
}
