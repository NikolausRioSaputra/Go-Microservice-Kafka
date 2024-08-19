package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/domain"
	"order-service/internal/repository"

	"github.com/segmentio/kafka-go"
)

type OrderUseCase interface {
	ProcessOrder(ctx context.Context, order domain.OrderRequest) error
	ListenForFailedOrders(ctx context.Context)
	ProcessEventRegistration(ctx context.Context, registration domain.EventRegistrationRequest) error
}

type orderUseCase struct {
	kafkaReader repository.KafkaReaderRepository
	kafkaWriter repository.KafkaWriterRepository
	orderRepo   repository.OrderRepository
}

func NewOrderUseCase(kafkaWriter repository.KafkaWriterRepository, orderRepo repository.OrderRepository, kafkaReader repository.KafkaReaderRepository) OrderUseCase {
	return &orderUseCase{
		kafkaReader: kafkaReader,
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
func (uc *orderUseCase) ProcessEventRegistration(ctx context.Context, registration domain.EventRegistrationRequest) error {
	registration, err := uc.orderRepo.SaveEventRegistration(ctx, registration)

	if err != nil {
		return err
	}

	messageSend := domain.Message{
		OrderType:     "Register Event",
		OrderService:  "start",
		OderID:        registration.OrderID,
		TransactionId: registration.TransactionID,
		Amount:        registration.Amount,
		UserId:        registration.UserID,
		EventName:     registration.EventName,
		PaymentMethod: registration.PaymentMethod,
		RespCode:      200,
		RespStatus:    "success",
		RespMessage:   "success register event",
	}

	bytes, err := json.Marshal(messageSend)

	if err != nil {
		return err
	}

	// Kirim pesan ke topik Kafka untuk validasi user
	message := kafka.Message{
		Key:   []byte(registration.TransactionID),
		Value: bytes,
	}

	return uc.kafkaWriter.WriteMessage(ctx, "topic_0", message)
}

func (uc *orderUseCase) ListenForFailedOrders(ctx context.Context) {
	for {
		message, err := uc.kafkaReader.ReadMessage(ctx)
		fmt.Println(string(message.Value))
		if err != nil {
			log.Printf("Error while reading message: %v\n", err)
			continue
		}

		var receivedMessage domain.Message
		err = json.Unmarshal(message.Value, &receivedMessage)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		if receivedMessage.RespStatus == "Failed" {
			log.Printf("Order failed for TransactionID: %s, OrderID: %s", receivedMessage.TransactionId, receivedMessage.OderID)

			// Update the order status to "cancelled" in the database
			err = uc.orderRepo.UpdateOrderStatus(ctx, receivedMessage.OderID, "cancelled")
			if err != nil {
				log.Printf("Failed to update order status: %v\n", err)
			} else {
				log.Printf("Order status updated to cancelled for OrderID: %s", receivedMessage.OderID)
			}
		}

		if receivedMessage.RespStatus == "Success" {
			log.Printf("Order success for TransactionID: %s, OrderID: %s", receivedMessage.TransactionId, receivedMessage.OderID)

			err = uc.orderRepo.UpdateOrderStatus(ctx, receivedMessage.OderID, "success")
			if err != nil {
				log.Printf("Failed to update order status: %v\n", err)
			} else {
				log.Printf("Order status updated to success for OrderID: %s", receivedMessage.OderID)
			}
		}
	}
}
