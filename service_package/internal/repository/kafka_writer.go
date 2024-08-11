package repository

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaWriterRepository interface {
	WriteMessage(ctx context.Context, key, value []byte) error
	/*
		ctx: context.Context yang digunakan untuk mengatur pembatalan atau timeout.
		key: Kunci pesan, biasanya digunakan untuk partisi pesan di Kafka.
		value: Nilai pesan yang akan dikirim.
	*/
	Close() error
}

type kafkaWriter struct {
	writer *kafka.Writer
}

func NewKafkaWriterRepository(brokers []string, topic string) KafkaWriterRepository {
	return &kafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (kw *kafkaWriter) WriteMessage(ctx context.Context, key, value []byte) error {
	//WriteMessages: Fungsi ini digunakan untuk mengirim pesan ke Kafka. Pesan ini dibungkus dalam kafka.Message, yang memiliki Key dan Value
	err := kw.writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Printf("Error writing message: %s\n", err)
		return err
	}
	log.Printf("RESPONSE SENT: %s\n", string(value))
	return nil
}

func (kw *kafkaWriter) Close() error {
	return kw.writer.Close()
}
