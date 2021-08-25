package sum

import (
	"github.com/gin-gonic/gin"
)

func Init(rg *gin.RouterGroup) {
	route := rg.Group("/info")

	Sum(route)
}
