package admin

type AuctionConfiguration struct {
	ID int64 `json:"id"`
	AuctionConfigurationAttrs
}

// AuctionConfigurationAttrs is attributes of Configuration. Used to create and update configurations
type AuctionConfigurationAttrs struct {
	Name       string                      `json:"name"`
	AppID      int64                       `json:"app_id"`
	AdType     AdType                      `json:"ad_type"`
	Rounds     []AuctionRoundConfiguration `json:"rounds"`
	Pricefloor float64                     `json:"pricefloor"`
}

type AdType string

const (
	UnknownAdType      AdType = ""
	BannerAdType       AdType = "banner"
	InterstitialAdType AdType = "interstitial"
	RewardedAdType     AdType = "rewarded"
)

type AuctionRoundConfiguration struct {
	ID      string   `json:"id"`
	Demands []string `json:"demands"`
	Timeout int      `json:"timeout"`
}

type AuctionConfigurationService = resourceService[AuctionConfiguration, AuctionConfigurationAttrs]
