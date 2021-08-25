package routes

import (
	"api/config"
	info "api/routes/info"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "api/docs"

	"github.com/gin-gonic/gin"
)

type SubRoutes struct {
	info info.Routes
}
type Routes struct {
	sub SubRoutes
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /api

// @Router /info [get]
func (r *Routes) Init(rg *gin.Engine) {
	route := rg.Group(config.ENV.BasePath)

	r.index(route)
	r.info(route)
	r.world(route)

	// Init SubRoutes
	r.sub.info.Init(route)

	// url := ginSwagger.URL("http://localhost:31337/api/swagger/doc.json")
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
