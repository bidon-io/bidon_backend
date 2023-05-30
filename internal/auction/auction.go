// Package auction is the core package that should contain all the domain logic related to auctioning.
// Currently, the system is really simple and is just a basic CRUD. So this package just defines domain structs and repository interfaces for actors to use directly.
package auction

import "context"

type Configuration struct {
	ID int64 `json:"id"`
	ConfigurationAttrs
}

// ConfigurationAttrs is attributes of Configuration. Used to create and update configurations
type ConfigurationAttrs struct {
	Name       string               `json:"name"`
	AppID      int64                `json:"app_id"`
	AdType     AdType               `json:"ad_type"`
	Rounds     []RoundConfiguration `json:"rounds"`
	Pricefloor float64              `json:"pricefloor"`
}

type AdType string

const (
	UnknownAdType      AdType = ""
	BannerAdType       AdType = "banner"
	InterstitialAdType AdType = "interstitial"
	RewardedAdType     AdType = "rewarded"
)

type RoundConfiguration struct {
	ID      string   `json:"id"`
	Demands []string `json:"demands"`
	Timeout int      `json:"timeout"`
}

type ConfigurationRepo interface {
	List(ctx context.Context) ([]Configuration, error)
	Find(ctx context.Context, id int64) (*Configuration, error)
	Create(ctx context.Context, attrs *ConfigurationAttrs) (*Configuration, error)
	Update(ctx context.Context, id int64, attrs *ConfigurationAttrs) (*Configuration, error)
	Delete(ctx context.Context, id int64) error
}
