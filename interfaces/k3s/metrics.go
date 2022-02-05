package k3s

import "github.com/gin-gonic/gin"

type Metrics interface {
	ReadAll(c *gin.Context)
	ReadOne(c *gin.Context)
}
