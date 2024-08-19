package repository

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaWriterRepository interface {
	WriteMessage(ctx context.Context, topic string, message kafka.Message) error
	Close() error
}

type kafkaWriter struct {
	brokers []string
}

func NewKafkaWriter(brokers []string) KafkaWriterRepository {
	return &kafkaWriter{
		brokers: brokers,
	}
}

func (kw *kafkaWriter) WriteMessage(ctx context.Context, topic string, message kafka.Message) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: kw.brokers,
		Topic:   topic,
	})
	defer writer.Close()

	return writer.WriteMessages(ctx, message)
}
// ->Metode ini bertanggung jawab untuk menulis atau mengirim satu pesan ke Kafka
func (kw *kafkaWriter) Close() error {
	// Kafka writer tidak memiliki resource yang tetap, jadi tidak ada yang perlu di-close di sini
	return nil
}
