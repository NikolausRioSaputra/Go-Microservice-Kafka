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
	viewTopic          repository.OcresRepository
	saveTransaction    repository.OcresRepository
}

// fungsi yang di gunakan untuk membuat fungsi baru
func NewKafkaUseCase(
	kr repository.KafkaReaderRepository,
	kw repository.KafkaWriterRepository,
	ka repository.KafkaWriterRepository,
	kp repository.KafkaWriterRepository,
	vt repository.OcresRepository,
	st repository.OcresRepository,
) KafkaUseCase {
	return &kafkaUseCase{
		kafkaReader:        kr,
		kafkaWriter:        kw,
		kafkaActivate:      ka,
		kafkaPaymentWriter: kp,
		viewTopic:          vt,
		saveTransaction:    st,
	}
}

func (uc *kafkaUseCase) ConsumeMessages(ctx context.Context) {
	for {
		message, err := uc.kafkaReader.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error while reading message: %v\n", err)
		}

		log.Printf("Received message: %s\n", string(message.Value))

		var incoming domain.Message
		if err := json.Unmarshal(message.Value, &incoming); err != nil {
			log.Printf("Error parsing message: %v\n", err)
			continue
		}

		// Ambil topik dan order service berikutnya dari database berdasarkan OrderType dan OrderService saat ini
		nextTopic, err := uc.viewTopic.ViewTopic(incoming.OrderType, incoming.OrderService)
		if err != nil {
			log.Printf("Error while retrieving topic: %v\n", err)
			continue
		}

		// Simpan transaksi ke dalam database dengan status "PROCESSED"
		transactionID, err := uc.saveTransaction.SaveTransaction(incoming, nextTopic, "SUCCESS")
		if err != nil {
			log.Printf("Error saving transaction: %v\n", err)
			continue
		}
		log.Printf("Transaction saved with ID: %d\n", transactionID)

		// Periksa apakah langkah berikutnya adalah "finish"
		if nextTopic == "finish" {
			_, err := uc.saveTransaction.SaveTransaction(incoming, nextTopic, "COMPLETED")
			if err != nil {
				log.Printf("Error saving final transaction: %v\n", err)
			}
			log.Printf("Transaction ID %s for order type '%s' is COMPLETED\n", incoming.TransactionId, incoming.OrderType)
			continue
		}

		// Update OrderService untuk langkah berikutnya
		incoming.OrderService = nextTopic

		// Serialisasi pesan kembali menjadi JSON
		responseBytes, _ := json.Marshal(incoming)

		// Menulis pesan ke Kafka pada topik yang ditentukan
		err = uc.kafkaWriter.WriteMessage(ctx, nextTopic, kafka.Message{
			Key:   []byte(incoming.TransactionId),
			Value: responseBytes,
		})
		if err != nil {
			log.Printf("Error writing message to %s: %v\n", nextTopic, err)
			continue
		}
		log.Printf("Message sent to %s: %s\n", nextTopic, string(responseBytes))
	}
}
