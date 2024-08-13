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
	reader := repository.NewKafkaReaderRepository([]string{"localhost:29092"}, "topic_validatePackage", "my-consumer-group") // -> membaca dari topic melalui broker di localhost, membaca topic topic_activatepackage
	writer := repository.NewKafkaWriterRepository([]string{"localhost:29092"}, "topic_0") // -> inisialisai menulis kafka yang sama dan menulis ke topic 0

	// Create the use case
	useCase := usecase.NewMessageUseCase() // -> membuat instance dari logika bisnis yang di terapkan pada pesan kafka yang di terima

	// Create the handler
	messageHandler := handler.NewMessageHandler(useCase, reader, writer) // -> instance dari handler yang menggabungkan logika bisnis (useCase) dengan kemampuan membaca (reader) dan menulis (writer) Kafka.

	//  memastikan bahwa resource reader dan writer akan ditutup ketika fungsi main selesai dieksekusi, meskipun terjadi error sebelumnya.
	defer reader.Close()
	defer writer.Close()

	fmt.Println("activatePackage is waiting for messages...")

	// Start processing messages
	messageHandler.ProcessMessages(context.Background()) // -> mulai memproses pesan dari Kafka menggunakan konteks context.Background() sebagai root context. Ini mungkin mencakup operasi asinkronus untuk menerima dan menulis pesan.
}
