package event

import (
	"strconv"
	"time"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/geocoder"
	"github.com/bidon-io/bidon-backend/internal/sdkapi/schema"
)

type TimingMap map[string][2]int64

type Event interface {
	Topic() config.Topic
}

func NewAdEvent(request *schema.BaseRequest, adRequestParams AdRequestParams, geoData geocoder.GeoData) *AdEvent {
	requestEvent := newBaseRequest(request, geoData)

	requestEvent.EventType = adRequestParams.EventType
	requestEvent.Status = adRequestParams.Status
	requestEvent.AdType = adRequestParams.AdType
	requestEvent.AdFormat = adRequestParams.AdFormat
	requestEvent.AuctionID = adRequestParams.AuctionID
	requestEvent.AuctionConfigurationID = adRequestParams.AuctionConfigurationID
	requestEvent.AuctionConfigurationUID = adRequestParams.AuctionConfigurationUID
	requestEvent.RoundID = adRequestParams.RoundID
	requestEvent.RoundNumber = adRequestParams.RoundNumber
	requestEvent.ImpID = adRequestParams.ImpID
	requestEvent.DemandID = adRequestParams.DemandID
	requestEvent.Bidding = adRequestParams.Bidding
	requestEvent.AdUnitUID = adRequestParams.AdUnitUID
	requestEvent.AdUnitInternalID = adRequestParams.AdUnitInternalID
	requestEvent.AdUnitLabel = adRequestParams.AdUnitLabel
	requestEvent.AdUnitCredentials = adRequestParams.AdUnitCredentials
	requestEvent.ECPM = adRequestParams.ECPM
	requestEvent.PriceFloor = adRequestParams.PriceFloor
	requestEvent.RawRequest = adRequestParams.RawRequest
	requestEvent.RawResponse = adRequestParams.RawResponse
	requestEvent.Error = adRequestParams.Error
	if adRequestParams.TimingMap == nil {
		requestEvent.TimingMap = make(TimingMap)
	} else {
		requestEvent.TimingMap = adRequestParams.TimingMap
	}
	requestEvent.ExternalWinnerDemandID = adRequestParams.ExternalWinnerDemandID
	requestEvent.ExternalWinnerEcpm = adRequestParams.ExternalWinnerEcpm
	requestEvent.Badv = adRequestParams.Badv
	requestEvent.Bcat = adRequestParams.Bcat
	requestEvent.Bapp = adRequestParams.Bapp

	return requestEvent
}

func newBaseRequest(request *schema.BaseRequest, geoData geocoder.GeoData) *AdEvent {
	segmentUID, err := strconv.Atoi(request.Segment.UID)
	if err != nil {
		segmentUID = 0
	}

	var model string
	if request.Device.OS == "android" {
		model = request.Device.Model
	} else {
		// for iOS and iPadOS, use hardware version like "iPad7,11" instead of model that is just "iPad"
		model = request.Device.HardwareVersion
	}

	return &AdEvent{
		Timestamp:                   generateTimestamp(),
		Manufacturer:                request.Device.Manufacturer,
		Model:                       model,
		Os:                          request.Device.OS,
		OsVersion:                   request.Device.OSVersion,
		ConnectionType:              request.Device.ConnectionType,
		DeviceType:                  string(request.Device.Type),
		UserAgent:                   request.Device.UserAgent,
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
		AppSetID:                    request.User.AppSetID,
		AppSetIDScope:               request.User.AppSetIDScope,
		COPPA:                       request.GetRegulations().COPPA,
		GDPR:                        request.GetRegulations().GDPR,
		CountryCode:                 geoData.CountryCode,
		City:                        geoData.CityName,
		Ip:                          geoData.IPString,
		CountryID:                   geoData.CountryID,
		SegmentID:                   request.Segment.ID,
		SegmentUID:                  int64(segmentUID),
		Ext:                         request.Ext,
		MediationMode:               request.GetMediationMode(),
		Mediator:                    request.GetMediator(),
		Session: Session{
			ID:                        request.Session.ID,
			LaunchTS:                  request.Session.LaunchTS,
			LaunchMonotonicTS:         request.Session.LaunchMonotonicTS,
			StartTS:                   request.Session.StartTS,
			StartMonotonicTS:          request.Session.StartMonotonicTS,
			TS:                        request.Session.TS,
			MonotonicTS:               request.Session.MonotonicTS,
			MemoryWarningsTS:          request.Session.MemoryWarningsTS,
			MemoryWarningsMonotonicTS: request.Session.MemoryWarningsMonotonicTS,
			RAMUsed:                   request.Session.RAMUsed,
			RAMSize:                   request.Session.RAMSize,
			StorageFree:               request.Session.StorageFree,
			StorageUsed:               request.Session.StorageUsed,
			Battery:                   request.Session.Battery,
			CPUUsage:                  request.Session.CPUUsage,
		},
	}
}

type AdRequestParams struct {
	EventType               string
	AdType                  string
	AdFormat                string
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
	AdUnitInternalID        int64
	AdUnitLabel             string
	AdUnitCredentials       map[string]string
	ECPM                    float64
	PriceFloor              float64
	RawRequest              string
	RawResponse             string
	Error                   string
	TimingMap               TimingMap
	ExternalWinnerDemandID  string
	ExternalWinnerEcpm      float64
	Badv                    string
	Bcat                    string
	Bapp                    string
}

const (
	SuccessAdRequestStatus = "SUCCESS"
	ErrorAdRequestStatus   = "ERROR"
)

type AdEvent struct {
	Timestamp                   float64           `json:"timestamp"`
	EventType                   string            `json:"event_type"`
	AdType                      string            `json:"ad_type"`
	AdFormat                    string            `json:"ad_format"`
	AuctionID                   string            `json:"auction_id"`
	AuctionConfigurationID      int64             `json:"auction_configuration_id"`
	AuctionConfigurationUID     int64             `json:"auction_configuration_uid"`
	Status                      string            `json:"status"`
	RoundID                     string            `json:"round_id"`
	RoundNumber                 int               `json:"round_number"`
	ImpID                       string            `json:"impid"`
	DemandID                    string            `json:"demand_id"`
	Bidding                     bool              `json:"bidding"`
	AdUnitUID                   int64             `json:"ad_unit_uid"`
	AdUnitInternalID            int64             `json:"ad_unit_internal_id"`
	AdUnitLabel                 string            `json:"ad_unit_label"`
	AdUnitCredentials           map[string]string `json:"ad_unit_credentials"`
	ECPM                        float64           `json:"ecpm"`
	PriceFloor                  float64           `json:"price_floor"`
	RawRequest                  string            `json:"raw_request"`
	RawResponse                 string            `json:"raw_response"`
	Error                       string            `json:"error"`
	TimingMap                   TimingMap         `json:"timing_map"`
	ExternalWinnerDemandID      string            `json:"external_winner_demand_id"`
	ExternalWinnerEcpm          float64           `json:"external_winner_ecpm"`
	Manufacturer                string            `json:"manufacturer"`
	Model                       string            `json:"model"`
	Os                          string            `json:"os"`
	OsVersion                   string            `json:"os_version"`
	ConnectionType              string            `json:"connection_type"`
	DeviceType                  string            `json:"device_type"`
	UserAgent                   string            `json:"user_agent"`
	SessionID                   string            `json:"session_id"`
	SessionUptime               int               `json:"session_uptime"`
	Bundle                      string            `json:"bundle"`
	Framework                   string            `json:"framework"`
	FrameworkVersion            string            `json:"framework_version"`
	PluginVersion               string            `json:"plugin_version"`
	PackageVersion              string            `json:"package_version"`
	SdkVersion                  string            `json:"sdk_version"`
	IDFA                        string            `json:"idfa"`
	IDG                         string            `json:"idg"`
	IDFV                        string            `json:"idfv"`
	TrackingAuthorizationStatus string            `json:"tracking_authorization_status"`
	AppSetID                    string            `json:"app_set_id"`
	AppSetIDScope               string            `json:"app_set_id_scope"`
	COPPA                       bool              `json:"coppa"`
	GDPR                        bool              `json:"gdpr"`
	CountryCode                 string            `json:"country_code"`
	City                        string            `json:"city"`
	Ip                          string            `json:"ip"`
	CountryID                   int64             `json:"country_id"`
	SegmentID                   string            `json:"segment_id"`
	SegmentUID                  int64             `json:"segment_uid"`
	Ext                         string            `json:"ext"`
	Session                     Session           `json:"session"`
	MediationMode               string            `json:"mediation_mode"`
	Mediator                    string            `json:"mediator"`
	Badv                        string            `json:"badv,omitempty"`
	Bcat                        string            `json:"bcat,omitempty"`
	Bapp                        string            `json:"bapp,omitempty"`
}

type Session struct {
	ID                        string   `json:"id"`
	LaunchTS                  int      `json:"launch_ts"`
	LaunchMonotonicTS         int      `json:"launch_monotonic_ts"`
	StartTS                   int      `json:"start_ts"`
	StartMonotonicTS          int      `json:"start_monotonic_ts"`
	TS                        int      `json:"ts"`
	MonotonicTS               int      `json:"monotonic_ts"`
	MemoryWarningsTS          []int    `json:"memory_warnings_ts"`
	MemoryWarningsMonotonicTS []int    `json:"memory_warnings_monotonic_ts"`
	RAMUsed                   int      `json:"ram_used"`
	RAMSize                   int      `json:"ram_size"`
	StorageFree               int      `json:"storage_free"`
	StorageUsed               int      `json:"storage_used"`
	Battery                   float64  `json:"battery"`
	CPUUsage                  *float64 `json:"cpu_usage"`
}

func (e *AdEvent) Topic() config.Topic {
	return config.AdEventsTopic
}

func generateTimestamp() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

func NewNotificationEvent(params NotificationParams) *NotificationEvent {
	errorString := ""
	if params.Error != nil {
		errorString = params.Error.Error()
	}

	return &NotificationEvent{
		Timestamp:   generateTimestamp(),
		EventType:   params.EventType,
		ImpID:       params.ImpID,
		Bundle:      params.Bundle,
		AdType:      params.AdType,
		AuctionID:   params.AuctionID,
		DemandID:    params.DemandID,
		LossReason:  params.LossReason,
		Price:       params.Price,
		FirstPrice:  params.FirstPrice,
		SecondPrice: params.SecondPrice,
		URL:         params.URL,
		TemplateURL: params.TemplateURL,
		Error:       errorString,
	}
}

type NotificationParams struct {
	EventType   string
	ImpID       string
	Bundle      string
	AdType      string
	AuctionID   string
	DemandID    string
	LossReason  int64
	Price       float64
	FirstPrice  float64
	SecondPrice float64
	URL         string
	TemplateURL string
	Error       error
}

type NotificationEvent struct {
	Timestamp   float64 `json:"timestamp"`
	EventType   string  `json:"event_type"`
	Bundle      string  `json:"bundle"`
	AdType      string  `json:"ad_type"`
	DemandID    string  `json:"demand_id"`
	AuctionID   string  `json:"auction_id"`
	ImpID       string  `json:"imp_id"`
	LossReason  int64   `json:"loss_reason"`
	Price       float64 `json:"ecpm"`
	FirstPrice  float64 `json:"first_price"`
	SecondPrice float64 `json:"second_price"`
	URL         string  `json:"url"`
	TemplateURL string  `json:"template_url"`
	Error       string  `json:"error"`
}

func (e *NotificationEvent) Topic() config.Topic {
	return config.NotificationEventsTopic
}
