package event_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/google/go-cmp/cmp"
)

func TestLogger_Log_ProcessesPayload(t *testing.T) {
	tests := []struct {
		name string
		in   map[string]any
		want map[string]any
	}{
		{
			name: "SimpleNestedMaps",
			in: map[string]any{
				"a": map[string]any{"c": "c"},
				"b": map[string]any{"d": "d"},
			},
			want: map[string]any{"a__c": "c", "b__d": "d"},
		},
		{
			name: "NestedMapsWithDifferentDepths",
			in: map[string]any{
				"a": map[string]any{"c": map[string]any{"e": "e"}},
				"b": map[string]any{"d": "d"},
			},
			want: map[string]any{"a__c__e": "e", "b__d": "d"},
		},
		{
			name: "KeyBHasAdditionalDepth",
			in: map[string]any{
				"a": map[string]any{"c": "c"},
				"b": map[string]any{"d": map[string]any{"f": "f"}},
			},
			want: map[string]any{"a__c": "c", "b__d__f": "f"},
		},
		{
			name: "KeyAHasDeepNesting",
			in: map[string]any{
				"a": map[string]any{"c": map[string]any{"e": map[string]any{"g": "g"}}},
				"b": map[string]any{"d": map[string]any{"f": "f"}},
			},
			want: map[string]any{"a__c__e__g": "g", "b__d__f": "f"},
		},
		{
			name: "MixedNestingAndPlainKey",
			in: map[string]any{
				"z": "z",
				"a": map[string]any{"c": map[string]any{"e": map[string]any{"g": "g", "h": "h"}}},
				"b": map[string]any{"d": map[string]any{"f": "f", "i": "i"}},
			},
			want: map[string]any{"z": "z", "a__c__e__g": "g", "a__c__e__h": "h", "b__d__f": "f", "b__d__i": "i"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ev := event.Event{Payload: test.in}
			mock := &event.LoggerEngineMock{
				ProduceFunc: func(_ event.Topic, message []byte, _ func(error)) {
					var unmarshalledMessage map[string]any
					err := json.Unmarshal(message, &unmarshalledMessage)
					if err != nil {
						t.Fatalf("%v: json.Unmarshal() failed: %v", test.name, err)
					}

					if diff := cmp.Diff(test.want, unmarshalledMessage); diff != "" {
						t.Errorf("%v: logger.Engine.Produce() got message diff (-want +got):\n%s", test.name, diff)
					}
				},
			}

			logger := &event.Logger{Engine: mock}
			logger.Log(ev, nil)
		})
	}
}

func TestLogger_Log_PassesTopic(t *testing.T) {
	tests := []struct {
		name string
		in   event.Topic
		want event.Topic
	}{
		{
			name: "KnownTopic",
			in:   event.ConfigTopic,
			want: event.ConfigTopic,
		},
		{
			name: "EmptyTopic",
			in:   event.Topic(""),
			want: event.Topic(""),
		},
		{
			name: "RandomTopic",
			in:   event.Topic("random"),
			want: event.Topic("random"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ev := event.Event{Topic: test.in}
			mock := &event.LoggerEngineMock{
				ProduceFunc: func(topic event.Topic, _ []byte, _ func(error)) {
					if topic != test.want {
						t.Errorf("%v: logger.Engine.Produce() got topic %v, want %v", test.name, topic, test.want)
					}
				},
			}
			logger := &event.Logger{Engine: mock}

			logger.Log(ev, nil)
		})
	}
}

func TestLogger_Log_CallsErrHandlerWithErrorFromEngine(t *testing.T) {
	err := errors.New("engine error")
	mock := &event.LoggerEngineMock{
		ProduceFunc: func(_ event.Topic, _ []byte, handleErr func(error)) {
			handleErr(err)
		},
	}
	logger := &event.Logger{Engine: mock}

	var handlerCalled bool
	var gotErr error
	logger.Log(event.Event{}, func(err error) {
		handlerCalled = true
		gotErr = err
	})

	if !handlerCalled {
		t.Fatalf("logger.Log() did not call error handler")
	}

	if gotErr == nil {
		t.Fatalf("logger.Log() called error handler with nil error")
	}
}

func TestLogger_Log_CallsErrHandlerWithErrorUnmarshallingInvalidJSON(t *testing.T) {
	mock := &event.LoggerEngineMock{ProduceFunc: func(_ event.Topic, _ []byte, _ func(error)) {}}
	logger := &event.Logger{Engine: mock}

	var handlerCalled bool
	var gotErr error
	logger.Log(event.Event{Payload: map[string]any{"foo": func() {}}}, func(err error) {
		handlerCalled = true
		gotErr = err
	})

	if !handlerCalled {
		t.Fatalf("logger.Log() did not call error handler")
	}

	if gotErr == nil {
		t.Fatalf("logger.Log() called error handler with nil error")
	}
}
