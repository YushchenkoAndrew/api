package sum

import (
	"github.com/gin-gonic/gin"
)

type Routes struct{}

func (r *Routes) Init(rg *gin.RouterGroup) {
	route := rg.Group("/info")

	r.sum(route)
}
