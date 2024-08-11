package repository

import (
	"context"

	"github.com/segmentio/kafka-go"
)
type KafkaReaderRepository interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type kafkaReader struct {
	reader *kafka.Reader
}

// fungsi nya membuat instansi baru dari kafkareader
func NewKafkaReader(brokers []string, topic, groupID string) KafkaReaderRepository {
	return &kafkaReader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers, // -> Broker adalah node dalam cluster Kafka yang menyimpan dan mengelola data Kafka.
			Topic:    topic, // -> Nama topic Kafka tempat pesan akan dibaca.
			GroupID:  groupID, // group untuk melacak offset dari pesan yang di konsomsi group
		}),
	}
}

func (kr *kafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return kr.reader.ReadMessage(ctx) // -> metode yang digunakan untuk membaca pesan dari kafka
}

func (kr *kafkaReader) Close() error {
	return kr.reader.Close()
}
