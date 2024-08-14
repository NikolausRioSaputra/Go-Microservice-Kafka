package main

import (
	"context"
	"service-orchestration/m/internal/provider/db"
	"service-orchestration/m/internal/repository"
	"service-orchestration/m/internal/usecase"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	database, err := db.GetConnection()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	// Initialize Kafka repositories
	kafkaReader := repository.NewKafkaReader([]string{"localhost:29092"}, "topic_0", "my-consumer-group")
	kafkaWriter := repository.NewKafkaWriter([]string{"localhost:29092"}) // Menghilangkan topik dari inisialisasi
	ocresrepo := repository.NewViewTopicRepository(database)
	kafkaUseCase := usecase.NewKafkaUseCase(kafkaReader, kafkaWriter, kafkaWriter, kafkaWriter, ocresrepo)

	wg.Add(1)
	// Start Kafka consumer
	go func() {
		time.Sleep(2 * time.Second)
		kafkaUseCase.ConsumeMessages(context.Background())
	}()

	wg.Wait()
}
