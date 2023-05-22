// Package store implements a database store for entities.
package store

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"gorm.io/gorm"
	"time"
)

type AuctionConfigurationRepo struct {
	DB *gorm.DB
}

func (r *AuctionConfigurationRepo) List(ctx context.Context) ([]auction.Configuration, error) {
	var dbConfigurations []auctionConfiguration
	if err := r.DB.WithContext(ctx).Find(&dbConfigurations).Error; err != nil {
		return nil, err
	}

	configurations := make([]auction.Configuration, len(dbConfigurations))
	for i, configuration := range dbConfigurations {
		configurations[i] = configuration.auctionConfiguration()
	}

	return configurations, nil
}

func (r *AuctionConfigurationRepo) Find(ctx context.Context, id uint) (*auction.Configuration, error) {
	var dbConfiguration auctionConfiguration
	if err := r.DB.WithContext(ctx).First(&dbConfiguration, id).Error; err != nil {
		return nil, err
	}

	configuration := dbConfiguration.auctionConfiguration()
	return &configuration, nil
}

func (r *AuctionConfigurationRepo) Create(ctx context.Context, configuration *auction.Configuration) error {
	dbConfiguration := fromAuctionConfiguration(configuration)
	if err := r.DB.WithContext(ctx).Create(&dbConfiguration).Error; err != nil {
		return err
	}

	*configuration = dbConfiguration.auctionConfiguration()

	return nil
}

func (r *AuctionConfigurationRepo) Update(ctx context.Context, configuration *auction.Configuration) error {
	dbConfiguration := fromAuctionConfiguration(configuration)
	if err := r.DB.WithContext(ctx).Save(&dbConfiguration).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuctionConfigurationRepo) Delete(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&auctionConfiguration{}, id).Error
}

// Model different thant default gorm.Model, because we already have schema from Rails
type Model struct {
	ID        int64     `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null"`
}

type auctionConfiguration struct {
	Model
	Name       *string                      `gorm:"column:name;type:varchar"`
	AppID      int64                        `gorm:"column:app_id;type:bigint;not null"`
	AdType     int32                        `gorm:"column:ad_type;type:integer;not null"`
	Rounds     []auction.RoundConfiguration `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json"`
	Pricefloor float64                      `gorm:"column:pricefloor;type:double precision;not null"`
}

func fromAuctionConfiguration(configuration *auction.Configuration) auctionConfiguration {
	return auctionConfiguration{
		Model:      Model{ID: int64(configuration.ID)},
		Name:       &configuration.Name,
		AppID:      int64(configuration.AppID),
		AdType:     dbAdType(configuration.AdType),
		Rounds:     configuration.Rounds,
		Pricefloor: configuration.Pricefloor,
	}
}

func (a *auctionConfiguration) auctionConfiguration() auction.Configuration {
	return auction.Configuration{
		ID:         uint(a.ID),
		Name:       *a.Name,
		AppID:      uint(a.AppID),
		AdType:     adType(a.AdType),
		Rounds:     a.Rounds,
		Pricefloor: a.Pricefloor,
	}
}

func adType(adType int32) auction.AdType {
	switch adType {
	case 1:
		return auction.InterstitialAdType
	case 3:
		return auction.BannerAdType
	case 6:
		return auction.RewardedAdType
	default:
		return auction.UnknownAdType
	}
}

func dbAdType(adType auction.AdType) int32 {
	switch adType {
	case auction.InterstitialAdType:
		return 1
	case auction.BannerAdType:
		return 3
	case auction.RewardedAdType:
		return 6
	default:
		return 0
	}
}
