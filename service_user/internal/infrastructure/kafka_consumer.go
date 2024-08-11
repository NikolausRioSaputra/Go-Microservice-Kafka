package infrastructure

import "github.com/segmentio/kafka-go"

func NewKafkaConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       "topic_validateUser",
		GroupID:     "my-consumer-group",
		StartOffset: kafka.LastOffset,
	})
}
