package sum

import (
	c "api/controllers/info"

	"github.com/gin-gonic/gin"
)

func (*Routes) sum(rg *gin.RouterGroup) {
	route := rg.Group("/sum")
	cSum := c.SumController{}

	route.GET("", cSum.Read)
}
