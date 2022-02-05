package info

import "github.com/gin-gonic/gin"

type Default interface {
	Read(c *gin.Context)
}
