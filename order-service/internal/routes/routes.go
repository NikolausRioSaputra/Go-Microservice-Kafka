package routes

import (
	"net/http"
	"order-service/internal/handler"
	"order-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, orderHandler *handler.OrderHandler) {
	// Token route
	router.GET("/token", func(c *gin.Context) {
		token, err := middleware.GenerateToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Order Routes
	orderRoutes := router.Group("/orders")
	orderRoutes.Use(middleware.AuthMiddleware(), middleware.LoggingMiddleware())
	{
		orderRoutes.POST("/create", orderHandler.CreateOrder)
		// Tambahkan rute terkait order lainnya di sini
	}

	// Tambahkan group route lain jika diperlukan
}
