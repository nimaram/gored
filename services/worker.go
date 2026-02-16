package services

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func Worker() {
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
			time.Sleep(750 * time.Millisecond) // sth or a kind of heavy job
			log.Println("Job done: ", string(m.Value))
		}(msg)
	}
}
