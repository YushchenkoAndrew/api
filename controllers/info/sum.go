package info

import (
	"github.com/gin-gonic/gin"
)

type SumController struct{}

func (o *SumController) Read(c *gin.Context) {
	// result, _ := o.filterQuery(c)
	// result = result.Offset(page * 20).Limit(limit)

	// helper.ResHandler(c, http.StatusOK, &gin.H{
	// 	"status":     "OK",
	// 	"result":     info,
	// 	"items":      result.RowsAffected,
	// 	"totalItems": 20,
	// })
}
