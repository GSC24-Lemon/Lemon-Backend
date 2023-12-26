// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"lemon_be/internal/usecase"

	// Swagger docs.
	_ "lemon_be/docs"
	"lemon_be/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title      Lemon_Be
// @description asdasd
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, a usecase.Auth, ws usecase.Websocket, cg usecase.Caregiver) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("/v1")
	{
		newAuthRoutes(h, a, l)
		NewWebsocketRoutes(h, ws, l)
		newCaregiverRoutes(h, cg, l)
	}
}
