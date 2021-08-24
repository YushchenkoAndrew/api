package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{}

func (*Controller) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
