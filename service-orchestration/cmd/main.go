package main

import (
	"context"
	"log"
	"service-orchestration/m/internal/handler"
	"service-orchestration/m/internal/provider/db"
	"service-orchestration/m/internal/repository"
	"service-orchestration/m/internal/usecase"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	var wg sync.WaitGroup

	// Mendapatkan koneksi ke database
	database, err := db.GetConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v\n", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing the database connection: %v\n", err)
		}
	}()

	// Inisialisasi repository Kafka dan database
	kafkaReader := repository.NewKafkaReader([]string{"localhost:29092"}, "topic_0", "my-consumer-group")
	kafkaWriter := repository.NewKafkaWriter([]string{"localhost:29092"})
	ocresRepo := repository.NewOcresRepository(database)
	transactionRepo := repository.NewOcresRepository(database)

	// Inisialisasi use case Kafka
	kafkaUseCase := usecase.NewKafkaUseCase(kafkaReader, kafkaWriter, kafkaWriter, kafkaWriter, ocresRepo, transactionRepo)

	// Inisialisasi handler untuk transaksi
	transactionHandler := handler.NewTransactionHandler(ocresRepo, kafkaWriter)

	// Menambahkan 1 ke wait group untuk Kafka consumer
	wg.Add(1)

	// Memulai Kafka consumer dalam goroutine
	go func() {
		defer wg.Done()             // Pastikan Done dipanggil setelah fungsi selesai
		time.Sleep(2 * time.Second) // Delay untuk memberikan waktu Kafka siap

		log.Println("Starting Kafka message consumption...")
		kafkaUseCase.ConsumeMessages(context.Background())
	}()

	// Memulai HTTP server untuk menangani permintaan HTTP
	router := gin.Default()
	router.GET("/transactions", transactionHandler.GetAllTransactions) // Menggunakan GetAllTransactions handler
	router.PUT("/transactions/:transactionId/resend", transactionHandler.UpdateItemIdAndResend)

	// Jalankan HTTP server dalam goroutine
	go func() {
		if err := router.Run(":8181"); err != nil {
			log.Fatalf("Failed to run HTTP server: %v\n", err)
		}
	}()

	// Menunggu semua goroutine selesai
	wg.Wait()
}
