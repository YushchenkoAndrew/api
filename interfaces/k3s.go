package interfaces

import "github.com/gin-gonic/gin"

type K3s interface {
	Subscribe(c *gin.Context)
	Unsubscribe(c *gin.Context)
}
