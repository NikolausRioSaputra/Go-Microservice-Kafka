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

		var topic string

		switch incoming.OrderType {
		case "Buy Package":
			switch incoming.OrderService {
			case "":
				incoming.OrderService = "validateUser"
				topic = "topic_validateUser"

			case "validateUser":
				incoming.OrderService = "validatePackage"
				topic = "topic_validatePackage"

			case "validatePackage":
				incoming.OrderService = "processPayment"
				topic = "topic_processPayment"

			case "processPayment":
				log.Printf("Transaction ID %s for order type '%s' is COMPLETED\n", incoming.TransactionId, incoming.OrderType)
				continue
			}

			responseBytes, _ := json.Marshal(incoming)
			err = uc.kafkaWriter.WriteMessage(ctx, topic, kafka.Message{
				Key:   []byte(incoming.TransactionId),
				Value: responseBytes,
			})
			if err != nil {
				log.Printf("Error writing message to %s: %v\n", topic, err)
				continue
			}
			log.Printf("Message sent to %s: %s\n", topic, string(responseBytes))

		default:
			log.Printf("Received unsupported message format: %v\n", incoming)
		}
	}
}
