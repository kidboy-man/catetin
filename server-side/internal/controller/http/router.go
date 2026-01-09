package http

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ingunawandra/catetin/internal/controller/http/v1"
	"github.com/ingunawandra/catetin/internal/controller/http/middleware"
)

// RouterConfig holds the configuration for setting up routes
type RouterConfig struct {
	AuthHandler *v1.AuthHandler
	// Add more handlers here as needed
}

// SetupRouter sets up the HTTP router with all routes
func SetupRouter(config *RouterConfig) *gin.Engine {
	// Create Gin router
	router := gin.Default()

	// Apply error handler middleware globally
	router.Use(middleware.ErrorHandler())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "catetin-api",
		})
	})

	// API v1 routes
	v1Group := router.Group("/api/v1")
	{
		// Authentication routes
		authGroup := v1Group.Group("/authentications")
		{
			authGroup.POST("/register", config.AuthHandler.Register)
			authGroup.POST("/login", config.AuthHandler.Login)
		}

		// Future routes
		// userGroup := v1Group.Group("/users")
		// expenseGroup := v1Group.Group("/expenses")
		// webhookGroup := v1Group.Group("/webhook")
	}

	return router
}
