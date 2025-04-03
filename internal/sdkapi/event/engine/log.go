package engine

import (
	"context"
	"log"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
)

type Log struct{}

func (e *Log) Produce(message event.LogMessage, _ func(error)) {
	topic := message.Topic
	value := message.Value
	log.Printf("PRODUCE EVENT %T(%v): %s", topic, topic, value)
}

func (e *Log) Ping(_ context.Context) error {
	return nil
}
