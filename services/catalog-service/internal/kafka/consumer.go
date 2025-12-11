package kafka

import "github.com/segmentio/kafka-go"

func InitReader(kafkaUrl, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  []string{kafkaUrl},
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		},
	)
}
