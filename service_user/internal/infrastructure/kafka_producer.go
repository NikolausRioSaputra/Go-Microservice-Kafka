package infrastructure

import "github.com/segmentio/kafka-go"

func NewKafkaProducer(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "topic_0",
		Balancer: &kafka.LeastBytes{},
	})
}
