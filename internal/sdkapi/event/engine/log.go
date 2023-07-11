package engine

import (
	"context"
	"log"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
)

type Log struct {
	Topics map[event.Topic]string
}

func (e *Log) Produce(ctx context.Context, topic event.Topic, message []byte, handleErr func(error)) {
	log.Printf("PRODUCE EVENT %T(%v): %s", topic, topic, message)
}
