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
