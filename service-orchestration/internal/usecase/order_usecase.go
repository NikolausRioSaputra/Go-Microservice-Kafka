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

	message := kafka.Message{
		Key:   []byte(order.TransactionID),
		Value: []byte(`{"orderType": "` + order.OrderType + `", "transactionId": "` + order.TransactionID + `", "userId": "` + order.UserId + `", "packageId": "` + order.PackageId + `"}`),
	}
	return uc.kafkaWriter.WriteMessage(ctx, message)
}
