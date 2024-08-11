package main

import (
	"context"
	"fmt"
	"user_service/internal/handler"
	"user_service/internal/repository"
	"user_service/internal/usecase"
)

func main() {
	// Initialize Kafka Reader and Writer Repositories
	reader := repository.NewKafkaReaderRepository([]string{"localhost:29092"}, "topic_validateUser", "my-consumer-group")
	writer := repository.NewKafkaWriterRepository([]string{"localhost:29092"}, "topic_0")

	// Create the use case
	useCase := usecase.NewMessageUseCase()

	// Create the handler
	messageHandler := handler.NewMessageHandler(useCase, reader, writer)

	defer reader.Close()
	defer writer.Close()

	fmt.Println("validateUser is waiting for messages...")

	// Start processing messages
	messageHandler.ProcessMessages(context.Background())
}
