package handler

import (
	"context"
	"encoding/json"
	"log"
	"service-package/internal/domain"
	"service-package/internal/repository"
	"service-package/internal/usecase"
)

type MessageHandler struct {
	useCase     usecase.MessageUseCase
	kafkaReader repository.KafkaReaderRepository
	kafkaWriter repository.KafkaWriterRepository
}

func NewMessageHandler(uc usecase.MessageUseCase, kr repository.KafkaReaderRepository, kw repository.KafkaWriterRepository) *MessageHandler {
	return &MessageHandler{
		useCase:     uc,
		kafkaReader: kr,
		kafkaWriter: kw,
	}
}

func (mh *MessageHandler) ProcessMessages(ctx context.Context) {
	for {
		msgBytes, err := mh.kafkaReader.ReadMessage(ctx)
		if err != nil {
			continue
		}

		var msg domain.Message

		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Error parsing message: %s\n", err)
			continue
		}

		response, err := mh.useCase.CheckItem(ctx, msg)

		if err != nil {
			log.Printf("Error activating package: %s\n", err)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %s\n", err)
			continue
		}

		err = mh.kafkaWriter.WriteMessage(ctx, []byte(msg.TransactionId), responseBytes)

		if err != nil {
			continue
		}
	}
}
