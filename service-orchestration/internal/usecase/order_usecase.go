package usecase

import (
	"context"
	"service-orchestration/m/internal/domain"
	"service-orchestration/m/internal/repository"

	"github.com/segmentio/kafka-go"
)
type OrderUseCase interface {
	ProcessOrder(ctx context.Context, order domain.OrderRequest) error
}

type orderUseCase struct {
	kafkaWriter repository.KafkaWriterRepository
	orderRepo   repository.OrderRepository
}

func NewOrderUseCase(kafkaWriter repository.KafkaWriterRepository, orderRepo repository.OrderRepository) OrderUseCase {
	return &orderUseCase{
		kafkaWriter: kafkaWriter,
		orderRepo:   orderRepo,
	}
}

func (uc *orderUseCase) ProcessOrder(ctx context.Context, order domain.OrderRequest) error {
	if err := uc.orderRepo.SaveOrder(ctx, order); err != nil {
		return err
	}

	// Tentukan topik Kafka berdasarkan tipe order atau kondisi lainnya
	var topic string
	switch order.OrderType {
	case "Buy Package":
		topic = "topic_validateUser"
	default:
		topic = "default_topic" // Topik default jika tipe order tidak dikenali
	}

	// Membuat pesan Kafka
	message := kafka.Message{
		Key:   []byte(order.TransactionID),
		Value: []byte(`{"orderType": "` + order.OrderType + `", "transactionId": "` + order.TransactionID + `", "userId": "` + order.UserId + `", "packageId": "` + order.PackageId + `"}`),
	}

	// Menulis pesan ke Kafka dengan topik yang sesuai
	return uc.kafkaWriter.WriteMessage(ctx, topic, message)
}
