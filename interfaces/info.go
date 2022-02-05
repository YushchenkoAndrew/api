package interfaces

import "github.com/gin-gonic/gin"

type Info interface {
	Create(c *gin.Context)
	CreateOne(c *gin.Context)
	CreateAll(c *gin.Context)

	ReadOne(c *gin.Context)
	ReadAll(c *gin.Context)

	UpdateOne(c *gin.Context)
	UpdateAll(c *gin.Context)

	DeleteOne(c *gin.Context)
	DeleteAll(c *gin.Context)
}
