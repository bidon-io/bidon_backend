package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"golang.org/x/exp/maps"
)

type Event struct {
	Topic   Topic
	Payload map[string]any
}

type Topic string

const (
	UnknownTopic Topic = ""
	ConfigTopic  Topic = "config"
)

type RequestMapper interface {
	Map() map[string]any
}

// Prepare creates new Event for Logger to log.
// It is safe to use created Event when error is returned.
func Prepare(topic Topic, mapper RequestMapper, geoData geocoder.GeoData) (Event, error) {
	payload := mapper.Map()

	payload["timestamp"] = float64(time.Now().UnixNano()) / 1e9
	payload["geo"] = eventGeo(payload["geo"], geoData)

	ext, err := eventExt(payload["ext"])

	payload["ext"] = ext

	return Event{
		Topic:   topic,
		Payload: payload,
	}, err
}

func eventGeo(geo any, geoData geocoder.GeoData) map[string]any {
	var eventGeo map[string]any
	if geoData != (geocoder.GeoData{}) {
		eventGeo = map[string]any{
			"ip":         geoData.IPString,
			"country":    geoData.CountryCode,
			"country_id": geoData.CountryID,
		}
	} else {
		eventGeo = make(map[string]any)
	}

	payloadGeo, _ := geo.(map[string]any)
	if payloadGeo == nil {
		return eventGeo
	}

	maps.Copy(payloadGeo, eventGeo)

	return payloadGeo
}

func eventExt(ext any) (map[string]any, error) {
	eventExt := make(map[string]any)

	payloadExt, ok := ext.(string)
	if !ok || payloadExt == "" {
		return eventExt, nil
	}

	err := json.Unmarshal([]byte(payloadExt), &eventExt)
	if err != nil {
		return eventExt, fmt.Errorf("unmarshal ext: %v", err)
	}

	return eventExt, nil
}
