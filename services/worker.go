package services

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)


func main(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{RedpandaUrl},
		Topic:   "tasks",
		GroupID: "workers",
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		go func(m kafka.Message) {
			log.Println("processing:", string(m.Value))
			time.Sleep(2 * time.Second) // sth or a kind of heavy job
		}(msg)
	}
}
