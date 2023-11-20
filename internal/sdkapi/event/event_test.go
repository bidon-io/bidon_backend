package event

import (
	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
)

func TestNewRequest(t *testing.T) {
	request := &schema.BaseRequest{}
	adRequestParams := AdRequestParams{}
	geoData := geocoder.GeoData{IPString: "8.8.8.8", CountryCode: "US", CountryID: 1}
	event := NewRequest(request, adRequestParams, geoData)

	testImplementEvent(t, event)

	if event.Topic() != config.AdEventsTopic {
		t.Errorf("NewRequest: expected topic %v, got %v", config.AdEventsTopic, event.Topic())
	}
}

func testImplementEvent(t *testing.T, e any) {
	_, ok := e.(Event)
	if !ok {
		t.Errorf("NewRequest: expected event to implement Event interface")
	}
}
