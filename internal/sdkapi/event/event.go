package event

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type TimingMap map[string][2]int64

type Event interface {
	Topic() config.Topic
	Children() []Event
	json.Marshaler
}

func NewConfig(request *schema.ConfigRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ConfigRequest]{
		timestamp: generateTimestamp(),
		topic:     config.ConfigTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewShow(request *schema.ShowRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ShowRequest]{
		timestamp: generateTimestamp(),
		topic:     config.ShowTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewClick(request *schema.ClickRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.ClickRequest]{
		timestamp: generateTimestamp(),
		topic:     config.ClickTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewReward(request *schema.RewardRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.RewardRequest]{
		timestamp: generateTimestamp(),
		topic:     config.RewardTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewRequest(request *schema.BaseRequest, adRequestParams AdRequestParams, geoData geocoder.GeoData) RequestEvent {
	requestEvent := newBaseRequest(request, geoData)

	requestEvent.EventType = adRequestParams.EventType
	requestEvent.Status = adRequestParams.Status
	requestEvent.AdType = adRequestParams.AdType
	requestEvent.AuctionID = adRequestParams.AuctionID
	requestEvent.AuctionConfigurationID = adRequestParams.AuctionConfigurationID
	requestEvent.AuctionConfigurationUID = adRequestParams.AuctionConfigurationUID
	requestEvent.RoundID = adRequestParams.RoundID
	requestEvent.RoundNumber = adRequestParams.RoundNumber
	requestEvent.ImpID = adRequestParams.ImpID
	requestEvent.DemandID = adRequestParams.DemandID
	requestEvent.Bidding = adRequestParams.Bidding
	requestEvent.AdUnitUID = adRequestParams.AdUnitUID
	requestEvent.AdUnitLabel = adRequestParams.AdUnitLabel
	requestEvent.Ecpm = adRequestParams.Ecpm
	requestEvent.PriceFloor = adRequestParams.PriceFloor
	requestEvent.RawRequest = adRequestParams.RawRequest
	requestEvent.RawResponse = adRequestParams.RawResponse
	requestEvent.Error = adRequestParams.Error
	requestEvent.TimingMap = adRequestParams.TimingMap

	return requestEvent
}

func NewLoss(request *schema.LossRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.LossRequest]{
		timestamp: generateTimestamp(),
		topic:     config.LossTopic,
		request:   request,
		geoData:   geoData,
	}
}

func NewWin(request *schema.WinRequest, geoData geocoder.GeoData) Event {
	return &simpleEvent[*schema.WinRequest]{
		timestamp: generateTimestamp(),
		topic:     config.WinTopic,
		request:   request,
		geoData:   geoData,
	}
}

func newBaseRequest(request *schema.BaseRequest, geoData geocoder.GeoData) RequestEvent {
	segmentUID, err := strconv.Atoi(request.Segment.UID)
	if err != nil {
		segmentUID = 0
	}

	return RequestEvent{
		Timestamp:                   generateTimestamp(),
		Manufacturer:                request.Device.Manufacturer,
		Model:                       request.Device.Model,
		Os:                          request.Device.OS,
		OsVersion:                   request.Device.OSVersion,
		ConnectionType:              request.Device.ConnectionType,
		DeviceType:                  string(request.Device.Type),
		SessionID:                   request.Session.ID,
		SessionUptime:               request.Session.Uptime(),
		Bundle:                      request.App.Bundle,
		Framework:                   request.App.Framework,
		FrameworkVersion:            request.App.FrameworkVersion,
		PluginVersion:               request.App.PluginVersion,
		PackageVersion:              request.App.Version,
		SdkVersion:                  request.App.SDKVersion,
		IDFA:                        request.User.IDFA,
		IDG:                         request.User.IDG,
		IDFV:                        request.User.IDFV,
		TrackingAuthorizationStatus: request.User.TrackingAuthorizationStatus,
		COPPA:                       request.GetRegulations().COPPA,
		GDPR:                        request.GetRegulations().GDPR,
		CountryCode:                 geoData.CountryCode,
		City:                        geoData.CityName,
		Ip:                          geoData.IPString,
		CountryID:                   geoData.CountryID,
		SegmentID:                   request.Segment.ID,
		SegmentUID:                  int64(segmentUID),
		Ext:                         request.Ext,
	}
}

type simpleEvent[T mapper] struct {
	timestamp float64
	topic     config.Topic
	request   T
	geoData   geocoder.GeoData
}

func (e *simpleEvent[T]) MarshalJSON() ([]byte, error) {
	payload, err := e.Payload()
	if err != nil {
		return nil, err
	}

	return json.Marshal(payload)
}

func (e *simpleEvent[T]) Topic() config.Topic {
	return e.topic
}

func (e *simpleEvent[T]) Payload() (map[string]any, error) {
	return prepareEventPayload(e.timestamp, e.request, e.geoData)
}

func (e *simpleEvent[T]) Children() []Event {
	return nil
}

type AdRequestParams struct {
	EventType               string
	AdType                  string
	AuctionID               string
	AuctionConfigurationID  int64
	AuctionConfigurationUID int64
	Status                  string
	RoundID                 string
	RoundNumber             int
	ImpID                   string
	DemandID                string
	Bidding                 bool
	AdUnitUID               int64
	AdUnitLabel             string
	Ecpm                    float64
	PriceFloor              float64
	RawRequest              string
	RawResponse             string
	Error                   string
	TimingMap               TimingMap
}

type RequestEvent struct {
	Timestamp                   float64   `json:"timestamp"`
	EventType                   string    `json:"event_type"`
	AdType                      string    `json:"ad_type"`
	AuctionID                   string    `json:"auction_id"`
	AuctionConfigurationID      int64     `json:"auction_configuration_id"`
	AuctionConfigurationUID     int64     `json:"auction_configuration_uid"`
	Status                      string    `json:"status"`
	RoundID                     string    `json:"round_id"`
	RoundNumber                 int       `json:"round_number"`
	ImpID                       string    `json:"impid"`
	DemandID                    string    `json:"demand_id"`
	Bidding                     bool      `json:"bidding"`
	AdUnitUID                   int64     `json:"ad_unit_uid"`
	AdUnitLabel                 string    `json:"ad_unit_label"`
	Ecpm                        float64   `json:"ecpm"`
	PriceFloor                  float64   `json:"price_floor"`
	RawRequest                  string    `json:"raw_request"`
	RawResponse                 string    `json:"raw_response"`
	Error                       string    `json:"error"`
	TimingMap                   TimingMap `json:"timing_map"`
	Manufacturer                string    `json:"manufacturer"`
	Model                       string    `json:"model"`
	Os                          string    `json:"os"`
	OsVersion                   string    `json:"os_version"`
	ConnectionType              string    `json:"connection_type"`
	DeviceType                  string    `json:"device_type"`
	SessionID                   string    `json:"session_id"`
	SessionUptime               int       `json:"session_uptime"`
	Bundle                      string    `json:"bundle"`
	Framework                   string    `json:"framework"`
	FrameworkVersion            string    `json:"framework_version"`
	PluginVersion               string    `json:"plugin_version"`
	PackageVersion              string    `json:"package_version"`
	SdkVersion                  string    `json:"sdk_version"`
	IDFA                        string    `json:"idfa"`
	IDG                         string    `json:"idg"`
	IDFV                        string    `json:"idfv"`
	TrackingAuthorizationStatus string    `json:"tracking_authorization_status"`
	COPPA                       bool      `json:"coppa"`
	GDPR                        bool      `json:"gdpr"`
	CountryCode                 string    `json:"country_code"`
	City                        string    `json:"city"`
	Ip                          string    `json:"ip"`
	CountryID                   int64     `json:"country_id"`
	SegmentID                   string    `json:"segment_id"`
	SegmentUID                  int64     `json:"segment_uid"`
	Ext                         string    `json:"ext"`
}

func (b RequestEvent) MarshalJSON() ([]byte, error) {
	type Alias RequestEvent
	return json.Marshal((Alias)(b))
}

func (b RequestEvent) Topic() config.Topic {
	return config.AdEventsTopic
}

func (b RequestEvent) Children() []Event {
	return nil
}

func generateTimestamp() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

type mapper interface {
	Map() map[string]any
}

func prepareEventPayload(timestamp float64, requestMapper mapper, geoData geocoder.GeoData) (map[string]any, error) {
	requestMap := requestMapper.Map()

	requestMap["timestamp"] = timestamp

	geo, _ := requestMap["geo"].(map[string]any)
	requestMap["geo"] = enhanceEventGeo(geo, geoData)

	ext, _ := requestMap["ext"].(string)
	eventExt, err := unmarshalEventExt(ext)
	requestMap["ext"] = eventExt

	if _, showPresent := requestMap["show"]; !showPresent {
		if bid, bidPresent := requestMap["bid"]; bidPresent {
			requestMap["show"] = bid
		}
	}

	return smashMap(requestMap, nil), err
}

func enhanceEventGeo(geo map[string]any, geoData geocoder.GeoData) map[string]any {
	if geo == nil {
		geo = make(map[string]any)
	}

	if geoData != (geocoder.GeoData{}) {
		geo["ip"] = geoData.IPString
		geo["country"] = geoData.CountryCode
		geo["country_id"] = geoData.CountryID
	}

	return geo
}

func unmarshalEventExt(ext string) (map[string]any, error) {
	result := make(map[string]any)

	if ext == "" {
		return result, nil
	}

	err := json.Unmarshal([]byte(ext), &result)
	if err != nil {
		return result, fmt.Errorf("unmarshal ext: %v", err)
	}

	return result, nil
}

func smashMap(src, dst map[string]any, nesting ...string) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}
	prefix := strings.Join(nesting, "__")

	for key, value := range src {
		switch mapValue := value.(type) {
		case map[string]any:
			n := slices.Clone(nesting)
			n = append(n, key)
			smashMap(mapValue, dst, n...)
		case []map[string]any:
			for i, v := range mapValue {
				n := slices.Clone(nesting)
				n = append(n, fmt.Sprintf("%s__%d", key, i))
				smashMap(v, dst, n...)
			}
		default:
			if prefix != "" {
				dst[fmt.Sprintf("%s__%s", prefix, key)] = value
			} else {
				dst[key] = value
			}
		}
	}

	return dst
}
