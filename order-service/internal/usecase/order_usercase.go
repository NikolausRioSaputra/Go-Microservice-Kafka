package usecase

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/repository"

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

	// Kirim pesan ke topik Kafka untuk validasi user
	message := kafka.Message{
		Key: []byte(order.TransactionID),
		Value: []byte(`{
			"orderType": "Buy Item",
			"transactionId": "` + order.TransactionID + `",
			"userId": "` + order.UserId + `",
			"itemId": "` + order.ItemId + `",
			"orderService": "start"}`),
	}

	return uc.kafkaWriter.WriteMessage(ctx, "topic_0", message)
}
