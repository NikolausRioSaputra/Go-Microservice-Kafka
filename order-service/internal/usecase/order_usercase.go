package usecase

import (
	"context"
	"encoding/json"
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

	order, err := uc.orderRepo.SaveOrder(ctx, order)

	if err != nil {
		return err
	}

	MessageSend := domain.Message{
		OrderType:     "Buy Item",
		OrderService:  "start",
		OderID:        order.OrderID,
		TransactionId: order.TransactionID,
		UserId:        order.UserId,
		ItemId:        order.ItemId,
		OrderAmount:   order.OrderAmount,
		PaymentMethod: order.PaymentMethod,
		RespCode:      200,
		RespStatus:    "success",
		RespMessage:   "succes create order",
	}

	bytes, err := json.Marshal(MessageSend)

	if err != nil {
		return err
	}

	// Kirim pesan ke topik Kafka untuk validasi user
	message := kafka.Message{
		Key:   []byte(order.TransactionID),
		Value: bytes,
	}

	return uc.kafkaWriter.WriteMessage(ctx, "topic_0", message)
}
