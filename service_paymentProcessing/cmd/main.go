package main

import (
	"context"
	"fmt"
	"service_paymentProcessing/internal/handler"
	"service_paymentProcessing/internal/repository"
	"service_paymentProcessing/internal/usecase"
)

func main() {
	// Initialize Kafka Reader and Writer Repositories
	reader := repository.NewKafkaReaderRepository([]string{"localhost:29092"}, "topic_processPayment", "my-consumer-group")
	writer := repository.NewKafkaWriterRepository([]string{"localhost:29092"}, "topic_0")

	// Create the use case
	useCase := usecase.NewPaymentUseCase()

	// Create the handler
	messageHandler := handler.NewPaymentMessageHandler(useCase, reader, writer)

	defer reader.Close()
	defer writer.Close()

	fmt.Println("PaymentProcessing Service is waiting for messages...")

	// Start processing messages
	messageHandler.ProcessMessages(context.Background())
}
