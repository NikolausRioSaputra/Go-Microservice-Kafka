package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"service-orchestration/m/internal/handler"
	"service-orchestration/m/internal/middleware"
	"service-orchestration/m/internal/provider/db"
	"service-orchestration/m/internal/repository"
	"service-orchestration/m/internal/routes"
	"service-orchestration/m/internal/usecase"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(2) // Mengatur jumlah prosesor yang digunakan
	database, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close() // Menutup koneksi database setelah selesai

	var wg sync.WaitGroup

	// Initialize repositories
	// Initialize Kafka repositories
	//-> ini nge kosume pesan dari topic dangan consumer group
	kafkaReader := repository.NewKafkaReader([]string{"localhost:29092"}, "topic_0", "my-consumer-group")
	kafkaWriter := repository.NewKafkaWriter([]string{"localhost:29092"}, "topic_validateUser")
	kafkaActivateWriter := repository.NewKafkaWriter([]string{"localhost:29092"}, "topic_activatePackage")
	kafkaUseCase := usecase.NewKafkaUseCase(kafkaReader, kafkaWriter, kafkaActivateWriter)

	// 1.Initialize Order repository
	// 2. Initialize use cases
	// 3. menagani logica order dan pakai kafka writer untuk publis pesan
	// 4. bertangung jawab mengkonsumsi pesan kafka
	orderRepo := repository.NewOrderRepository(database)
	orderUseCase := usecase.NewOrderUseCase(kafkaWriter, orderRepo)
	orderHandler := handler.NewOrderHandler(orderUseCase) //-> mengatur rute terkait penguna


	// Initialize user repository
	// Initialize user use case
	// Initialize handlers
	userRepo := repository.NewUserRepository(database)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase) // -> mengatur rute penguna

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Initialize routes
	routes.InitializeRoutes(router, orderHandler, userHandler)

	// Middleware Setup
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

	// Start Kafka consumer
	go func() {
		time.Sleep(2 * time.Second)
		kafkaUseCase.ConsumeMessages(context.Background())
	}()

	wg.Wait()
}
