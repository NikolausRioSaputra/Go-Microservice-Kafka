package repository

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaReaderRepository interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type kafkaReader struct {
	reader *kafka.Reader
}

func NewKafkaReader(brokers []string, topic, groupID string) KafkaReaderRepository {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &kafkaReader{
		reader: r,
	}
}

func (kr *kafkaReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	message, err := kr.reader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message: %v", err)
		return kafka.Message{}, err
	}
	return message, nil
}

func (kr *kafkaReader) Close() error {
	return kr.reader.Close()
}
