package usecase

import (
	"context"
	"encoding/json"
	"fmt"


	"service-orchestration/m/internal/repository"

	"github.com/segmentio/kafka-go"
)

type TransactionUseCase struct {
	repo        repository.OcresRepository
	kafkaWriter repository.KafkaWriterRepository
}

func NewTransactionUseCase(repo repository.OcresRepository, kafkaWriter repository.KafkaWriterRepository) *TransactionUseCase {
	return &TransactionUseCase{
		repo:        repo,
		kafkaWriter: kafkaWriter,
	}
}
func (uc *TransactionUseCase) UpdateItemIdAndSendKafkaMessage(transactionId, newItemId string) error {
    // Ambil transaksi yang ada
    transaction, err := uc.repo.GetTransactionByID(transactionId)
    if err != nil {
        return err
    }

    // Unmarshal payload JSON
    var payload map[string]interface{}
    if err := json.Unmarshal([]byte(transaction.Payload), &payload); err != nil {
        return err
    }

    // Update itemId
    payload["itemId"] = newItemId

      // Marshal kembali payload
	  updatedPayload, err := json.Marshal(payload)
	  if err != nil {
		  return fmt.Errorf("error marshaling updated payload: %v", err)
	  }

     // Update payload di database
	 if err := uc.repo.UpdateTransactionPayload(transactionId, string(updatedPayload), "validateItem"); err != nil {
        return fmt.Errorf("error updating transaction payload: %v", err)
    }

    // Kirim pesan ke Kafka
    message := kafka.Message{
        Key:   []byte(transactionId),
        Value: updatedPayload,
    }

    if err := uc.kafkaWriter.WriteMessage(context.Background(), "topic_validateItem", message); err != nil {
        return err
    }

    return nil
}
