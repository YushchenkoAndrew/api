package info

import (
	c "api/controllers/info"
	"api/interfaces"
	"api/interfaces/info"

	"github.com/gin-gonic/gin"
)

type rangeRouter struct {
	route  *gin.RouterGroup
	cRange info.Default
}

func NewRangeRouterFactory() func(*gin.RouterGroup) interfaces.Router {
	return func(rg *gin.RouterGroup) interfaces.Router {
		return &rangeRouter{route: rg.Group("/range"), cRange: c.NewRangeController()}
	}
}

func (c *rangeRouter) Init() {
	c.route.GET("", c.cRange.Read)
}
