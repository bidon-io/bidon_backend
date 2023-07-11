package event

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Logger struct {
	Engine LoggerEngine
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out mocks_test.go . LoggerEngine
type LoggerEngine interface {
	Produce(ctx context.Context, topic Topic, message []byte, handleErr func(error))
}

func (l *Logger) Log(ctx context.Context, event Event, handleErr func(error)) {
	payload := make(map[string]any)
	smashMap(payload, event.Payload)

	message, err := json.Marshal(payload)
	if err != nil {
		handleErr(fmt.Errorf("marshal event payload: %v", err))
	}

	l.Engine.Produce(ctx, event.Topic, message, handleErr)
}

func smashMap(dst, src map[string]any, nesting ...string) {
	prefix := strings.Join(nesting, "__")

	for key, value := range src {
		mapValue, ok := value.(map[string]any)
		if ok {
			n := slices.Clone(nesting)
			n = append(n, key)
			smashMap(dst, mapValue, n...)
		} else if prefix != "" {
			dst[fmt.Sprintf("%s__%s", prefix, key)] = value
		} else {
			dst[key] = value
		}
	}
}
