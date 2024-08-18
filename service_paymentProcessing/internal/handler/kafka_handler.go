package handler

import (
	"context"
	"encoding/json"
	"log"
	"service_paymentProcessing/internal/domain"
	"service_paymentProcessing/internal/repository"
	"service_paymentProcessing/internal/usecase"
)

type PaymentMessageHandler struct {
	useCase     usecase.PaymentUseCase
	kafkaReader repository.KafkaReaderRepository
	kafkaWriter repository.KafkaWriterRepository
}

func NewPaymentMessageHandler(uc usecase.PaymentUseCase, kr repository.KafkaReaderRepository, kw repository.KafkaWriterRepository) *PaymentMessageHandler {
	return &PaymentMessageHandler{
		useCase:     uc,
		kafkaReader: kr,
		kafkaWriter: kw,
	}
}

func (pmh *PaymentMessageHandler) ProcessMessages(ctx context.Context) {
	for {
		msgBytes, err := pmh.kafkaReader.ReadMessage(ctx)
		if err != nil {
			continue
		}

		var msg domain.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Error parsing message: %s\n", err)
			continue
		}

		response, err := pmh.useCase.ProcessPayment(ctx, msg)
		if err != nil {
			log.Printf("Error processing payment: %s\n", err)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %s\n", err)
			continue
		}

		err = pmh.kafkaWriter.WriteMessage(ctx, []byte(msg.TransactionId), responseBytes)
		if err != nil {
			continue
		}
	}
}
