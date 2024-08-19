package handler

import (
	"context"
	"encoding/json"
	"log"
	"service-package/internal/domain"
	"service-package/internal/repository"
	"service-package/internal/usecase"
)

type EventHandler struct {
	useCase     usecase.EventUseCase
	kafkaReader repository.KafkaReaderRepository
	kafkaWriter repository.KafkaWriterRepository
}

func NewEventHandler(uc usecase.EventUseCase, kr repository.KafkaReaderRepository, kw repository.KafkaWriterRepository) *EventHandler {
	return &EventHandler{
		useCase:     uc,
		kafkaReader: kr,
		kafkaWriter: kw,
	}
}

func (eh *EventHandler) ProcessMessages(ctx context.Context) {
	for {
		msgBytes, err := eh.kafkaReader.ReadMessage(ctx)
		if err != nil {
			continue
		}

		var request domain.EventRegistrationRequest

		if err := json.Unmarshal(msgBytes, &request); err != nil {
			log.Printf("Error parsing message: %s\n", err)
			continue
		}

		response, err := eh.useCase.ProcessEventRegistration(ctx, request)

		if err != nil {
			log.Printf("Error processing event registration: %s\n", err)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshalling response: %s\n", err)
			continue
		}

		err = eh.kafkaWriter.WriteMessage(ctx, []byte(request.TransactionID), responseBytes)

		if err != nil {
			continue
		}
	}
}
