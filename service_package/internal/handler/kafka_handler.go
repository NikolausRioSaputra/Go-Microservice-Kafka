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
		//mh.kafkaReader.ReadMessage(ctx) memanggil Kafka Reader untuk mengambil pesan. Jika terjadi error (misalnya karena Kafka tidak tersedia), loop akan melanjutkan iterasi berikutnya, melewatkan pesan yang tidak bisa dibaca.
		if err != nil {
			continue
		}

		var msg domain.Message
		//Pesan yang diterima dalam bentuk byte array kemudian di-deserialize menjadi struct domain.Message. Jika gagal, log error dan lanjut ke iterasi berikutnya.
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Error parsing message: %s\n", err)
			continue
		}

		response, err := mh.useCase.ActivatePackage(ctx, msg)
		// msg kemudian diproses menggunakan useCase untuk melakukan aktivasi paket. Jika terjadi kesalahan dalam logika bisnis, error akan dilog dan iterasi akan dilanjutkan.
		if err != nil {
			log.Printf("Error activating package: %s\n", err)
			continue
		}

		responseBytes, err := json.Marshal(response)
		// Respons yang dihasilkan dari pemrosesan pesan kemudian di-serialize menjadi JSON byte array. Jika gagal, error akan dilog dan iterasi akan dilanjutkan.
		if err != nil {
			log.Printf("Error marshalling response: %s\n", err)
			continue
		}

		err = mh.kafkaWriter.WriteMessage(ctx, []byte(msg.TransactionId), responseBytes)
		// Pesan respons kemudian dikirim kembali ke Kafka menggunakan kafkaWriter. Jika ada error saat menulis pesan, iterasi akan dilanjutkan tanpa menghentikan aplikasi.
		if err != nil {
			continue
		}
	}
}
