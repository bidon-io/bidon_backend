// Package store implements a database store for entities.
package store

import (
	"time"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"gorm.io/gorm"
)

// AutoMigrate migrates all of the models in store
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&auctionConfiguration{},
		&segment{},
	)
}

// Model different thant default gorm.Model, because we already have schema from Rails
type Model struct {
	ID        int64     `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null"`
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
