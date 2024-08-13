package repository

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaWriterRepository interface {
	WriteMessage(ctx context.Context, message kafka.Message) error
	Close() error
}

type kafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriter(brokers []string, topic string) KafkaWriterRepository {
	return &kafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}
// ->Metode ini bertanggung jawab untuk menulis atau mengirim satu pesan ke Kafka
func (kw *kafkaWriter) WriteMessage(ctx context.Context, message kafka.Message) error {
	return kw.writer.WriteMessages(ctx, message)
}

func (kw *kafkaWriter) Close() error {
	return kw.writer.Close()
}
