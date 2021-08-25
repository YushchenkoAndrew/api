package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (*Controller) getID(c *gin.Context, id *int) bool {
	var err error
	if *id, err = strconv.Atoi(c.Param("id")); err != nil || *id <= 0 {
		return false
	}
	return true
}

func (*Controller) resHandler(c *gin.Context, stat int, data *gin.H) {
	switch c.GetHeader("Accept") {
	case "application/xml":
		c.XML(stat, *data)

	default:
		c.JSON(stat, *data)
	}
}

func (o *Controller) errHandler(c *gin.Context, stat int, message string) {
	o.resHandler(c, stat, &gin.H{
		"status":  "ERR",
		"result":  []string{},
		"message": message,
	})
}
