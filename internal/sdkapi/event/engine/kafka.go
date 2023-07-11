package engine

import (
	"context"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Kafka struct {
	Topics map[event.Topic]string
	Client *kgo.Client
}

func (e *Kafka) Produce(ctx context.Context, topic event.Topic, message []byte, handleErr func(error)) {
	topicStr, ok := e.Topics[topic]
	if !ok {
		handleErr(fmt.Errorf("unknown topic: %v", topic))
	}

	record := &kgo.Record{
		Topic: topicStr,
		Value: message,
	}
	e.Client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			handleErr(fmt.Errorf("kafka produce record: %v", err))
		}
	})
}
