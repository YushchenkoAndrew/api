package routes

import (
	_ "api/docs"
	"api/interfaces"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type swaggerRouter struct {
	route *gin.RouterGroup
}

func NewSwaggerRouter(route *gin.RouterGroup) interfaces.Router {
	return &swaggerRouter{route: route.Group("/swagger")}
}

func (c *swaggerRouter) Init() {
	c.route.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
