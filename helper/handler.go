package helper

import (
	"api/models"

	"github.com/gin-gonic/gin"
)

func ResHandler(c *gin.Context, stat int, data interface{}) {
	switch c.GetHeader("Accept") {
	case "application/xml":
		c.XML(stat, data)

	default:
		c.JSON(stat, data)
	}
}

func ErrHandler(c *gin.Context, stat int, message string) {
	ResHandler(c, stat, models.Error{
		Status:  "ERR",
		Message: message,
		Result:  []string{},
	})
}
