package repository

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaReaderRepository interface {
	ReadMessage(ctx context.Context) ([]byte, error) // ->  Metode ini membaca pesan dari Kafka dan mengembalikan isi pesan dalam bentuk byte array ([]byte).
	Close() error // -> Metode ini digunakan untuk menutup koneksi reader Kafka dengan aman
}

type kafkaReader struct {
	reader *kafka.Reader // -> Field ini digunakan untuk mengakses Kafka.
}

func NewKafkaReaderRepository(brokers []string, topic, groupID string) KafkaReaderRepository {
	return &kafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
		}),
	}
}

func (kr *kafkaReader) ReadMessage(ctx context.Context) ([]byte, error) {
	m, err := kr.reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message: %s\n", err)
		return nil, err
	}
	log.Printf("RECEIVED MESSAGE : %s\n", string(m.Value))
	return m.Value, nil
}

func (kr *kafkaReader) Close() error {
	return kr.reader.Close()
}
