package routes

import (
	c "api/controllers"
	"api/interfaces"
	"api/middleware"

	"github.com/gin-gonic/gin"
)

type subscriptionRouter struct {
	auth         *gin.RouterGroup
	subscription interfaces.Default
}

func NewSubscribeRouter(rg *gin.RouterGroup) interfaces.Router {
	return &subscriptionRouter{
		auth:         rg.Group("/subscription", middleware.Auth()),
		subscription: c.NewSubscriptionController(),
	}
}

func (c *subscriptionRouter) Init() {
	c.auth.POST("", c.subscription.CreateOne)

	c.auth.GET("/:id", c.subscription.ReadOne)
	c.auth.GET("", c.subscription.ReadAll)

	c.auth.PUT("/:id", c.subscription.UpdateOne)
	c.auth.PUT("", c.subscription.UpdateAll)

	c.auth.DELETE("/:id", c.subscription.DeleteOne)
	c.auth.DELETE("", c.subscription.DeleteAll)
}
