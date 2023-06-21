package schema

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/device"
)

type Request struct {
	AdType      ad.Type                 `param:"ad_type"`
	AdObject    AdObject                `json:"ad_object"`
	Device      Device                  `json:"device"`
	Session     Session                 `json:"session"`
	App         App                     `json:"app"`
	User        User                    `json:"user"`
	Geo         *Geo                    `json:"geo"`
	Regulations *Regulations            `json:"regs"`
	Adapters    map[adapter.Key]Adapter `json:"adapters"`
	Segment     Segment                 `json:"segment"`
	Token       string                  `json:"token"`
	Ext         string                  `json:"ext"`
}

type Geo struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Accuracy  float64 `json:"accuracy"`
	LastFix   int     `json:"lastfix"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	ZIP       string  `json:"zip"`
	UTCOffset int     `json:"utcoffset"`
}

type Device struct {
	Geo             *Geo        `json:"geo"`
	UserAgent       string      `json:"ua"`
	Manufacturer    string      `json:"make"`
	Model           string      `json:"model"`
	OS              string      `json:"os"`
	OSVersion       string      `json:"osv"`
	HardwareVersion string      `json:"hwv"`
	Height          int         `json:"h"`
	Width           int         `json:"w"`
	PPI             int         `json:"ppi"`
	PXRatio         float64     `json:"pxratio"`
	JS              int         `json:"js"`
	Language        string      `json:"language"`
	Carrier         string      `json:"carrier"`
	MCCMNC          string      `json:"mccmnc"`
	ConnectionType  string      `json:"connection_type"`
	Type            device.Type `json:"type"`
}

type Session struct {
	ID                        string  `json:"id"`
	LaunchTS                  int     `json:"launch_ts"`
	LaunchMonotonicTS         int     `json:"launch_monotonic_ts"`
	StartTS                   int     `json:"start_ts"`
	StartMonotonicTS          int     `json:"start_monotonic_ts"`
	TS                        int     `json:"ts"`
	MonotonicTS               int     `json:"monotonic_ts"`
	MemoryWarningsTS          []int   `json:"memory_warnings_ts"`
	MemoryWarningsMonotonicTS []int   `json:"memory_warnings_monotonic_ts"`
	RAMUsed                   int     `json:"ram_used"`
	RAMSize                   int     `json:"ram_size"`
	StorageFree               int     `json:"storage_free"`
	StorageUsed               int     `json:"storage_used"`
	Battery                   float64 `json:"battery"`
	CPUUsage                  float64 `json:"cpu_usage"`
}

type App struct {
	Bundle           string   `json:"bundle"`
	Key              string   `json:"key"`
	Framework        string   `json:"framework"`
	Version          string   `json:"version"`
	FrameworkVersion string   `json:"framework_version"`
	PluginVersion    string   `json:"plugin_version"`
	Skadn            []string `json:"skadn"`
}

type User struct {
	IDFA                        string         `json:"idfa"`
	TrackingAuthorizationStatus string         `json:"tracking_authorization_status"`
	IDFV                        string         `json:"idfv"`
	IDG                         string         `json:"idg"`
	Consent                     map[string]any `json:"consent"`
	COPPA                       *bool          `json:"coppa"`
}

type Regulations struct {
	COPPA bool `json:"coppa"`
	GDPR  bool `json:"gdpr"`
}

type Adapter struct {
	Version    string `json:"version"`
	SDKVersion string `json:"sdk_version"`
}

type AdObject struct {
	PlacementID  string                `json:"placement_id"`
	AuctionID    string                `json:"auction_id"`
	Orientation  string                `json:"orientation"`
	PriceFloor   float64               `json:"pricefloor"`
	Banner       *BannerAdObject       `json:"banner"`
	Interstitial *InterstitialAdObject `json:"interstitial"`
	Rewarded     *RewardedAdObject     `json:"rewarded"`
}

type Segment struct {
	ID  string `json:"id"`
	Ext string `json:"ext"`
}

func (o *AdObject) AdFormat() ad.Format {
	if o.Banner != nil {
		return o.Banner.Format
	}

	return ad.EmptyFormat
}

type BannerAdObject struct {
	Format ad.Format `json:"format"`
}

type InterstitialAdObject struct {
}

type RewardedAdObject struct {
}
