package event_test

import (
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/event"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type requestMapper map[string]any

func (m requestMapper) Map() map[string]any {
	return m
}

func TestPrepare(t *testing.T) {
	type in struct {
		topic   event.Topic
		mapper  event.RequestMapper
		geoData geocoder.GeoData
	}
	type want struct {
		event event.Event
		err   error
	}

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			"empty request and no geo data",
			in{
				event.ConfigTopic,
				requestMapper{},
				geocoder.GeoData{},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{},
						"ext": map[string]any{},
					},
				},
				nil,
			},
		},
		{
			"filled request and no geo data",
			in{
				event.ConfigTopic,
				requestMapper{
					"foo": "foo",
					"bar": "bar",
				},
				geocoder.GeoData{},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"foo": "foo",
						"bar": "bar",
						"geo": map[string]any{},
						"ext": map[string]any{},
					},
				},
				nil,
			},
		},
		{
			"empty request and present geo data",
			in{
				event.ConfigTopic,
				requestMapper{},
				geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{
							"ip":         "8.8.8.8",
							"country":    "US",
							"country_id": int64(1),
						},
						"ext": map[string]any{},
					},
				},
				nil,
			},
		},
		{
			"request with geo and present geo data",
			in{
				event.ConfigTopic,
				requestMapper{
					"geo": map[string]any{
						"ip":         "1.1.1.1",
						"country":    "GB",
						"country_id": 2,
						"ext":        "something",
					},
				},
				geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{
							"ip":         "8.8.8.8",
							"country":    "US",
							"country_id": int64(1),
							"ext":        "something",
						},
						"ext": map[string]any{},
					},
				},
				nil,
			},
		},
		{
			"request with empty ext and no geo data",
			in{
				event.ConfigTopic,
				requestMapper{
					"ext": "",
				},
				geocoder.GeoData{},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{},
						"ext": map[string]any{},
					},
				},
				nil,
			},
		},
		{
			"request with filled ext and no geo data",
			in{
				event.ConfigTopic,
				requestMapper{
					"ext": `{"foo": "foo"}`,
				},
				geocoder.GeoData{},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{},
						"ext": map[string]any{
							"foo": "foo",
						},
					},
				},
				nil,
			},
		},
		{
			"request with invalid ext and no geo data",
			in{
				event.ConfigTopic,
				requestMapper{
					"ext": `"foo": "foo"`,
				},
				geocoder.GeoData{},
			},
			want{
				event.Event{
					Topic: event.ConfigTopic,
					Payload: map[string]any{
						"geo": map[string]any{},
						"ext": map[string]any{},
					},
				},
				errors.New("message not important"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := event.Prepare(test.in.topic, test.in.mapper, test.in.geoData)

			ignoreTSEntry := cmpopts.IgnoreMapEntries(func(key string, _ any) bool {
				return key == "timestamp"
			})
			if diff := cmp.Diff(test.want.event, got, ignoreTSEntry); diff != "" {
				t.Errorf("%v: event.Prepare() mismatch (-out +got):\n%s", test.name, diff)
			}

			tsVal := got.Payload["timestamp"]
			if ts, _ := tsVal.(float64); ts == 0 {
				t.Errorf("%v: event.Prepare() got timestamp %T(%v), out non-zero float64", test.name, tsVal, tsVal)
			}

			if err != test.want.err && (err == nil || test.want.err == nil) {
				t.Errorf("%v: event.Prepare() got error %q, out error %q", test.name, err, test.want.err)
			}
		})
	}
}
