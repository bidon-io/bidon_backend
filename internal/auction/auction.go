// Package auction is the core package that should contain all the domain logic related to auctioning.
// Currently, the system is really simple and is just a basic CRUD. So this package just defines domain structs and repository interfaces for actors to use directly.
package auction

import "context"

type Configuration struct {
	ID         uint                 `json:"id"`
	Name       string               `json:"name"`
	AppID      uint                 `json:"app_id"`
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
	Find(ctx context.Context, id uint) (*Configuration, error)
	Create(ctx context.Context, configuration *Configuration) error
	Update(ctx context.Context, configuration *Configuration) error
	Delete(ctx context.Context, id uint) error
}
