package usecase

// use case yang bertanggung jawab untuk mengkonsumsi pesan dari Kafka, memprosesnya, dan mengirimkannya kembali ke Kafka pada topik yang berbeda berdasarkan tipe pesan.

import (
	"context"
	"encoding/json"
	"log"
	"service-orchestration/m/internal/domain"
	"service-orchestration/m/internal/repository"

	"github.com/segmentio/kafka-go"
)

type KafkaUseCase interface {
	ConsumeMessages(ctx context.Context)
}
type kafkaUseCase struct {
	kafkaReader        repository.KafkaReaderRepository // -> digunakan membaca pesan dikafka
	kafkaWriter        repository.KafkaWriterRepository // -> menuliskan pesan di kafka untuk topic valdasi
	kafkaActivate      repository.KafkaWriterRepository // -> digunakan untuk menulisakan pesan untuk aktivasi package
	kafkaPaymentWriter repository.KafkaWriterRepository // -> digunakan untuk menulisakan pesan untuk aktivasi package
}

// fungsi yang di gunakan untuk membuat fungsi baru
func NewKafkaUseCase(kr repository.KafkaReaderRepository, kw repository.KafkaWriterRepository, ka repository.KafkaWriterRepository, kp repository.KafkaWriterRepository) KafkaUseCase {
	return &kafkaUseCase{
		kafkaReader:        kr,
		kafkaWriter:        kw,
		kafkaActivate:      ka,
		kafkaPaymentWriter: kp,
	}
}

func (uc *kafkaUseCase) ConsumeMessages(ctx context.Context) {
	// fungsi loop tak henti digunakan untuk secara terus menerus mendeangar pesan dari kafka
	for {
		message, err := uc.kafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error while reading message: %v\n", err)
		}

		log.Printf("Received message: %s\n", string(message.Value))

		var incoming domain.IncomingMessage
		if err := json.Unmarshal(message.Value, &incoming); err != nil {
			log.Printf("Error parsing message: %v\n", err)
			continue
		}

		switch incoming.OrderType {
		case "Buy Package":
			switch incoming.OrderService {
			case "":
				// Step 1: Validate User
				incoming.OrderService = "validateUser"
				responseBytes, _ := json.Marshal(incoming)
				err = uc.kafkaWriter.WriteMessage(ctx, kafka.Message{
					Key:   []byte(incoming.TransactionId),
					Value: responseBytes,
				})
				if err != nil {
					log.Printf("Error writing message to topic_validateUser: %v\n", err)
					continue
				}
				log.Printf("Message sent to topic_validateUser: %s\n", string(responseBytes))

			case "validateUser":
				// Step 2: Validate Package
				incoming.OrderService = "validatePackage"
				responseBytes, _ := json.Marshal(incoming)
				err = uc.kafkaActivate.WriteMessage(ctx, kafka.Message{
					Key:   []byte(incoming.TransactionId),
					Value: responseBytes,
				})
				if err != nil {
					log.Printf("Error writing message to topic_validatePackage: %v\n", err)
					continue
				}
				log.Printf("Message sent to topic_validatePackage: %s\n", string(responseBytes))

			case "validatePackage":
				// Step 3: Process Payment
				incoming.OrderService = "processPayment"
				responseBytes, _ := json.Marshal(incoming)
				err = uc.kafkaPaymentWriter.WriteMessage(ctx, kafka.Message{
					Key:   []byte(incoming.TransactionId),
					Value: responseBytes,
				})
				if err != nil {
					log.Printf("Error writing message to topic_processPayment: %v\n", err)
					continue
				}
				log.Printf("Message sent to topic_processPayment: %s\n", string(responseBytes))

			case "processPayment":
				// Step 4: Complete Transaction
				log.Printf("===============================================================================================")
				log.Printf("Transaction ID %s for order type '%s' is COMPLETED\n", incoming.TransactionId, incoming.OrderType)
				log.Printf("===============================================================================================")
			}

		default:
			log.Printf("Received unsupported message format: %v\n", incoming)
		}
	}
}
