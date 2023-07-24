package event

import (
	"errors"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/google/go-cmp/cmp"
)

type testMapper map[string]any

func (m testMapper) Map() map[string]any {
	return m
}

func TestPrepareEventPayload(t *testing.T) {
	type in struct {
		timestamp     float64
		requestMapper mapper
		geoData       geocoder.GeoData
	}
	type want struct {
		payload map[string]any
		err     error
	}

	timestamp := generateTimestamp()

	tests := []struct {
		name string
		in   in
		want want
	}{
		{
			"empty request and no geo data",
			in{
				timestamp,
				testMapper{},
				geocoder.GeoData{},
			},
			want{
				map[string]any{
					"timestamp": timestamp,
				},
				nil,
			},
		},
		{
			"filled request and no geo data",
			in{
				timestamp,
				testMapper{
					"foo": map[string]any{
						"foo": map[string]any{
							"foo": "foo",
						},
						"bar": "bar",
						"baz": []map[string]any{
							{"foo": "foo"},
							{"bar": map[string]any{"bar": "bar"}},
						},
					},
					"bar": map[string]any{
						"bar": "bar",
					},
					"baz": "baz",
				},
				geocoder.GeoData{},
			},
			want{
				map[string]any{
					"timestamp":             timestamp,
					"foo__foo__foo":         "foo",
					"foo__bar":              "bar",
					"foo__baz__0__foo":      "foo",
					"foo__baz__1__bar__bar": "bar",
					"bar__bar":              "bar",
					"baz":                   "baz",
				},
				nil,
			},
		},
		{
			"empty request and present geo data",
			in{
				timestamp,
				testMapper{},
				geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1},
			},
			want{
				map[string]any{
					"timestamp":       timestamp,
					"geo__ip":         "8.8.8.8",
					"geo__country":    "US",
					"geo__country_id": int64(1),
				},
				nil,
			},
		},
		{
			"request with geo and present geo data",
			in{
				timestamp,
				testMapper{
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
				map[string]any{
					"timestamp":       timestamp,
					"geo__ip":         "8.8.8.8",
					"geo__country":    "US",
					"geo__country_id": int64(1),
					"geo__ext":        "something",
				},
				nil,
			},
		},
		{
			"request with empty ext and no geo data",
			in{
				timestamp,
				testMapper{
					"ext": "",
				},
				geocoder.GeoData{},
			},
			want{
				map[string]any{
					"timestamp": timestamp,
				},
				nil,
			},
		},
		{
			"request with filled ext and no geo data",
			in{
				timestamp,
				testMapper{
					"ext": `{"foo": {"foo": "foo"}, "bar": "bar"}`,
				},
				geocoder.GeoData{},
			},
			want{
				map[string]any{
					"timestamp":     timestamp,
					"ext__foo__foo": "foo",
					"ext__bar":      "bar",
				},
				nil,
			},
		},
		{
			"request with invalid ext and no geo data",
			in{
				timestamp,
				testMapper{
					"ext": `"foo": "foo"`,
				},
				geocoder.GeoData{},
			},
			want{
				map[string]any{
					"timestamp": timestamp,
				},
				errors.New("message not important"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := prepareEventPayload(test.in.timestamp, test.in.requestMapper, test.in.geoData)

			if diff := cmp.Diff(test.want.payload, got); diff != "" {
				t.Errorf("%v: prepareEventPayload() mismatch (-want +got):\n%s", test.name, diff)
			}

			if err != test.want.err && (err == nil || test.want.err == nil) {
				t.Errorf("%v: prepareEventPayload() got error %q, want error %q", test.name, err, test.want.err)
			}
		})
	}
}
