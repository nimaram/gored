package services

import (
	"context"
	"os"

	"github.com/segmentio/kafka-go"
)

var RedpandaUrl string = os.Getenv("REDPANDA_STRING_URL")

var writer = kafka.NewWriter(kafka.WriterConfig{
	Brokers: []string{RedpandaUrl},
	Topic:   "tasks",
})

func Publish(ctx context.Context, msg []byte) error {
	return writer.WriteMessages(ctx, kafka.Message{
		Value: msg,
	})
}

// Close writer gracefully
func Close() error {
	if writer != nil {
		return writer.Close()
	}
	return nil
}
