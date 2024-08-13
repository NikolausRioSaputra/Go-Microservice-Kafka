package infrastructure

import "github.com/segmentio/kafka-go"

func NewKafkaConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,          // Gunakan parameter brokers
		Topic:       topic,            // Gunakan parameter topic
		GroupID:     groupID,          // Gunakan parameter groupID
		StartOffset: kafka.LastOffset, // Mulai membaca dari pesan terbaru
	})
}
