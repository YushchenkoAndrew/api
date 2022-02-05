package k3s

import "github.com/gin-gonic/gin"

type Pods interface {
	ReadOne(c *gin.Context)
	ReadAll(c *gin.Context)

	Exec(c *gin.Context)
}
