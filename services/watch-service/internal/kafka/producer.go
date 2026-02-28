package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

var (
	writer *kafka.Writer
)

func InitKafkaProducer(brokerHost string) {
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerHost},
		Topic:   "anime-watch"})
}

func WriteMessage(ctx context.Context, msg kafka.Message) error {
	return writer.WriteMessages(ctx, msg)
}

func Close() {
	writer.Close()
}
