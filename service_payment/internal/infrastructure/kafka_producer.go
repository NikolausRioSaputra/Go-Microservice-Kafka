package infrastructure

import "github.com/segmentio/kafka-go"

func NewKafkaProducer(brokers []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{}, // -> Balancer ini akan mengarahkan pesan ke partisi dengan ukuran byte terkecil.
	})
}
