package event

import (
	"encoding/json"
	"fmt"

	"github.com/bidon-io/bidon-backend/config"
)

type Logger struct {
	Engine LoggerEngine
}

type LoggerEngine interface {
	Produce(message LogMessage, handleErr func(error))
}

type LogMessage struct {
	Topic config.Topic
	Value []byte
}

func (l *Logger) Log(event Event, handleErr func(error)) {
	topic := event.Topic()

	message, err := json.Marshal(event)
	if err != nil {
		handleErr(fmt.Errorf("marshal %q event payload: %v", topic, err))
	}

	logMessage := LogMessage{
		Topic: topic,
		Value: message,
	}

	l.Engine.Produce(logMessage, func(err error) {
		handleErr(fmt.Errorf("produce %q message: %v", logMessage.Topic, err))
	})
}
