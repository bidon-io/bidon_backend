package event

import (
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
)

func TestNewAdEvent(t *testing.T) {
	request := &schema.BaseRequest{}
	adRequestParams := AdRequestParams{}
	geoData := geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1}
	event := NewAdEvent(request, adRequestParams, geoData)

	testImplementEvent(t, event)

	if event.Topic() != config.AdEventsTopic {
		t.Errorf("NewAdEvent: expected topic %v, got %v", config.AdEventsTopic, event.Topic())
	}
}

func TestNewNotificationEvent(t *testing.T) {
	params := NotificationParams{
		EventType:   "AdRequest",
		ImpID:       "imp-1",
		DemandID:    "test",
		LossReason:  0,
		URL:         "test",
		TemplateURL: "test",
		Error:       nil,
	}

	event := NewNotificationEvent(params)

	testImplementEvent(t, event)

	if event.Topic() != config.NotificationEventsTopic {
		t.Errorf("NewAdEvent: expected topic %v, got %v", config.NotificationEventsTopic, event.Topic())
	}
}

func testImplementEvent(t *testing.T, e any) {
	_, ok := e.(Event)
	if !ok {
		t.Errorf("NewAdEvent: expected event to implement Event interface")
	}
}

func TestNewAdEventModelField(t *testing.T) {
	testCases := []struct {
		name     string
		deviceOS string
		model    string
		hardware string
		expected string
	}{
		{
			name:     "Android",
			deviceOS: "android",
			model:    "Pixel 4",
			hardware: "qcom",
			expected: "Pixel 4",
		},
		{
			name:     "iPhone",
			deviceOS: "iOS",
			model:    "iPhone",
			hardware: "iPhone12,1",
			expected: "iPhone12,1",
		},
		{
			name:     "iPad",
			deviceOS: "iPadOS",
			model:    "iPad",
			hardware: "iPad7,11",
			expected: "iPad7,11",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &schema.BaseRequest{
				Device: schema.Device{
					OS:              tc.deviceOS,
					Model:           tc.model,
					HardwareVersion: tc.hardware,
				},
			}
			adRequestParams := AdRequestParams{}
			geoData := geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1}
			event := NewAdEvent(request, adRequestParams, geoData)

			if event.Model != tc.expected {
				t.Errorf("NewAdEvent: expected model %v, got %v", tc.expected, event.Model)
			}
		})
	}
}
