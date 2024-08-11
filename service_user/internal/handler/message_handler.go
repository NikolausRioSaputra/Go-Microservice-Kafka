package handler

import (
	"context"
	"encoding/json"
	"log"
	"user_service/internal/domain"
	"user_service/internal/repository"
	"user_service/internal/usecase"
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

		response, err := mh.useCase.ValidateUser(ctx, msg)
		if err != nil {
			log.Printf("Error validating user: %s\n", err)
			continue
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
