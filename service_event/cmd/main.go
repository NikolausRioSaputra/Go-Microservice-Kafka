package main

import (
	"context"
	"fmt"
	"service-package/internal/handler"
	"service-package/internal/repository"
	"service-package/internal/usecase"
)

func main() {
	// Initialize Kafka Reader and Writer Repositories
	reader := repository.NewKafkaReaderRepository([]string{"localhost:29092"}, "topic_validate_event", "my-consumer-group")
	writer := repository.NewKafkaWriterRepository([]string{"localhost:29092"}, "topic_0")

	// Create the use case
	useCase := usecase.NewEventUseCase()

	// Create the handler
	eventHandler := handler.NewEventHandler(useCase, reader, writer)

	defer reader.Close()
	defer writer.Close()

	fmt.Println("Event service is waiting for messages...")

	// Start processing messages
	eventHandler.ProcessMessages(context.Background())
}
