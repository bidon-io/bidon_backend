package event

import (
	"encoding/json"
	"fmt"
)

type Logger struct {
	Engine LoggerEngine
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . LoggerEngine
type LoggerEngine interface {
	Produce(message LogMessage, handleErr func(error))
}

type LogMessage struct {
	Topic Topic
	Value []byte
}

func (l *Logger) Log(event Event, handleErr func(error)) {
	events := append(event.Children(), event)
	messages := make([]LogMessage, len(events))

	goodToProduce := true
	for i, event := range events {
		topic := event.Topic()

		payload, err := event.Payload()
		if err != nil {
			handleErr(fmt.Errorf("prepare %q event payload: %v", topic, err))
		}

		message, err := json.Marshal(payload)
		if err != nil {
			handleErr(fmt.Errorf("marshal %q event payload: %v", topic, err))
			goodToProduce = false
		}

		messages[i] = LogMessage{
			Topic: topic,
			Value: message,
		}
	}

	// Event must be batched together with children. Either all or none.
	if !goodToProduce {
		return
	}

	for _, message := range messages {
		message := message
		l.Engine.Produce(message, func(err error) {
			handleErr(fmt.Errorf("produce %q message: %v", message.Topic, err))
		})
	}
}
