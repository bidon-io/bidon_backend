package engine

import (
	"log/slog"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
)

type Log struct{}

func (e *Log) Produce(message event.LogMessage, _ func(error)) {
	topic := message.Topic
	value := message.Value
	slog.Info("PRODUCE EVENT", "topic", topic, "value", value)
}
