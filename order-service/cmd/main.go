package main

import (
	"fmt"
	"log"
	"order-service/internal/handler"
	"order-service/internal/middleware"
	"order-service/internal/provider/db"
	"order-service/internal/repository"
	"order-service/internal/routes"
	"order-service/internal/usecase"
	"sync"

	"github.com/gin-gonic/gin"
)

func main() {
	// Mengatur koneksi ke database
	database, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()

	var wg sync.WaitGroup

	// Inisialisasi repository dan use case
	kafkaWriter := repository.NewKafkaWriter([]string{"localhost:29092"})
	orderRepo := repository.NewOrderRepository(database)
	orderUseCase := usecase.NewOrderUseCase(kafkaWriter, orderRepo)
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// Inisialisasi Gin dan route
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	routes.InitializeRoutes(router, orderHandler)

	// Middleware
	router.Use(middleware.LoggingMiddleware(), middleware.AuthMiddleware())

	fmt.Printf("Server running on :8080\n")
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := router.Run(":8080")
		if err != nil {
			log.Fatal("Server failed to start: ", err)
		}
	}()

	wg.Wait()
}
