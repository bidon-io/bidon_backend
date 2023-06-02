// Package store implements a database store for entities.
package store

import (
	"time"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"gorm.io/gorm"
)

// AutoMigrate migrates all of the models in store
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&auctionConfiguration{},
		&segment{},
		&app{},
		&appDemandProfile{},
		&demandSource{},
		&demandSourceAccount{},
		&lineItem{},
		&country{},
		&user{},
	)
}

// Model different thant default gorm.Model, because we already have schema from Rails
type Model struct {
	ID        int64     `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null"`
}

func adType(adType int32) admin.AdType {
	switch adType {
	case 1:
		return admin.InterstitialAdType
	case 3:
		return admin.BannerAdType
	case 6:
		return admin.RewardedAdType
	default:
		return admin.UnknownAdType
	}
}

func dbAdType(adType admin.AdType) int32 {
	switch adType {
	case admin.InterstitialAdType:
		return 1
	case admin.BannerAdType:
		return 3
	case admin.RewardedAdType:
		return 6
	default:
		return 0
	}
}

func platformID(platformID int32) admin.PlatformID {
	switch platformID {
	case 1:
		return admin.AndroidPlatformID
	case 4:
		return admin.IOSPlatformID
	default:
		return admin.UnknownPlatformID
	}
}

func dbPlatformID(platformID admin.PlatformID) int32 {
	switch platformID {
	case admin.AndroidPlatformID:
		return 1
	case admin.IOSPlatformID:
		return 4
	default:
		return 0
	}
}
