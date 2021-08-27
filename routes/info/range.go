package sum

import (
	c "api/controllers/info"

	"github.com/gin-gonic/gin"
)

func Range(rg *gin.RouterGroup) {
	route := rg.Group("/range")
	cRange := c.RangeController{}

	route.GET("", cRange.Read)
}
