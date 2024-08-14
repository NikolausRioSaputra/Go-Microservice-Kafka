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

func (kw *kafkaWriter) Close() error {
	return nil
}
