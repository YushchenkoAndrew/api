package interfaces

import "github.com/gin-gonic/gin"

type Index interface {
	Ping(c *gin.Context)
	TraceIp(c *gin.Context)

	Login(c *gin.Context)
	Refresh(c *gin.Context)
}
