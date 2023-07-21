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

func (e *Kafka) Produce(message event.LogMessage, handleErr func(error)) {
	topic := message.Topic
	topicStr := e.Topics[topic]
	if topicStr == "" {
		handleErr(fmt.Errorf("unknown topic: %v", topic))
	}

	record := &kgo.Record{
		Topic: topicStr,
		Value: message.Value,
	}
	e.Client.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			handleErr(fmt.Errorf("kafka produce record: %v", err))
		}
	})
}
