package helper

import "github.com/gin-gonic/gin"

func ResHandler(c *gin.Context, stat int, data *gin.H) {
	switch c.GetHeader("Accept") {
	case "application/xml":
		c.XML(stat, *data)

	default:
		c.JSON(stat, *data)
	}
}

func ErrHandler(c *gin.Context, stat int, message string) {
	ResHandler(c, stat, &gin.H{
		"status":  "ERR",
		"result":  []string{},
		"message": message,
	})
}
