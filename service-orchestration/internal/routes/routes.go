package routes

import (
	"net/http"
	"service-orchestration/m/internal/handler"
	"service-orchestration/m/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, orderHandler *handler.OrderHandler, userHandler *handler.UserHandler) {
	// Token route
	router.GET("/token", func(c *gin.Context) {
		token, err := middleware.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/create", userHandler.StoreNewUser)
		// userRoutes.GET("/:id", userHandler.GetUser)
	}

	// Order Routes
	orderRoutes := router.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware(), middleware.LoggingMiddleware())
	{
		orderRoutes.POST("/create", orderHandler.CreateOrder)
		// Tambahkan rute terkait order lainnya di sini
	}

	// Tambahkan group route lain jika diperlukan
}
