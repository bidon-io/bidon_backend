package db

import (
	"database/sql"
	"time"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/shopspring/decimal"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Open(databaseURL string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL))
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Begin(opts ...*sql.TxOptions) *DB {
	return &DB{db.DB.Begin(opts...)}
}

func (db *DB) AutoMigrate() error {
	return db.DB.AutoMigrate(
		&App{},
		&AppDemandProfile{},
		&AuctionConfiguration{},
		&Country{},
		&DemandSourceAccount{},
		&DemandSource{},
		&LineItem{},
		&Segment{},
		&User{},
	)
}

// Model different than default gorm.Model, because we already have schema from Rails
type Model struct {
	ID        int64     `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null"`
}

type AppDemandProfile struct {
	Model
	AppID          int64          `gorm:"column:app_id;type:bigint;not null"`
	AccountType    string         `gorm:"column:account_type;type:varchar;not null"`
	AccountID      int64          `gorm:"column:account_id;type:bigint;not null"`
	DemandSourceID int64          `gorm:"column:demand_source_id;type:bigint;not null"`
	Data           map[string]any `gorm:"column:data;type:jsonb;default:'{}';serializer:json"`
}

type App struct {
	Model
	UserID      int64          `gorm:"column:user_id;type:bigint;not null"`
	PlatformID  PlatformID     `gorm:"column:platform_id;type:integer;not null"`
	HumanName   string         `gorm:"column:human_name;type:varchar;not null"`
	PackageName sql.NullString `gorm:"column:package_name;type:varchar"`
	AppKey      sql.NullString `gorm:"column:app_key;type:varchar"`
	Settings    map[string]any `gorm:"column:settings;type:jsonb;default:'{}';serializer:json"`
}

type AuctionConfiguration struct {
	Model
	Name       sql.NullString                    `gorm:"column:name;type:varchar"`
	AppID      int64                             `gorm:"column:app_id;type:bigint;not null"`
	AdType     AdType                            `gorm:"column:ad_type;type:integer;not null"`
	Rounds     []admin.AuctionRoundConfiguration `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json"`
	Pricefloor float64                           `gorm:"column:pricefloor;type:double precision;not null"`
}

type Country struct {
	Model
	Alpha2Code string         `gorm:"column:alpha_2_code;type:varchar;not null"`
	Alpha3Code string         `gorm:"column:alpha_3_code;type:varchar;not null"`
	HumanName  sql.NullString `gorm:"column:human_name;type:varchar"`
}

type DemandSourceAccount struct {
	Model
	DemandSourceID int64 `gorm:"column:demand_source_id;type:bigint;not null"`
	DemandSource   DemandSource
	UserID         int64          `gorm:"column:user_id;type:bigint;not null"`
	Type           string         `gorm:"column:type;type:varchar;not null"`
	Extra          map[string]any `gorm:"column:extra;type:jsonb;default:'{}';serializer:json"`
	IsBidding      *bool          `gorm:"column:bidding;type:boolean;default:false"`
	IsDefault      sql.NullBool   `gorm:"column:is_default;type:boolean"`
}

type DemandSource struct {
	Model
	APIKey    string `gorm:"column:api_key;type:varchar;not null"`
	HumanName string `gorm:"column:human_name;type:varchar;not null"`
}

type LineItem struct {
	Model
	AppID       int64  `gorm:"column:app_id;type:bigint;not null"`
	AccountType string `gorm:"column:account_type;type:varchar;not null"`
	AccountID   int64  `gorm:"column:account_id;type:bigint;not null"`
	Account     DemandSourceAccount
	HumanName   string              `gorm:"column:human_name;type:varchar;not null"`
	Code        *string             `gorm:"column:code;type:varchar;not null"`
	BidFloor    decimal.NullDecimal `gorm:"column:bid_floor;type:numeric"`
	AdType      AdType              `gorm:"column:ad_type;type:integer;not null"`
	Extra       map[string]any      `gorm:"column:extra;type:jsonb;default:'{}';serializer:json"`
	Width       int32               `gorm:"column:width;type:integer;default:0;not null"`
	Height      int32               `gorm:"column:height;type:integer;default:0;not null"`
	Format      sql.NullString      `gorm:"column:format;type:varchar"`
}

type Segment struct {
	Model
	Name        string                `gorm:"column:name;type:varchar;not null"`
	Description string                `gorm:"column:description;type:text;not null"`
	Filters     []admin.SegmentFilter `gorm:"column:filters;type:jsonb;not null;default:'[]';serializer:json"`
	Enabled     *bool                 `gorm:"column:enabled;type:bool;not null;default:true"`
	AppID       int64                 `gorm:"column:app_id;type:bigint;not null"`
}

type User struct {
	Model
	Email string `gorm:"column:email;type:varchar;not null"`
}

type AdType int32

const (
	UnknownAdType      AdType = 0
	InterstitialAdType AdType = 1
	BannerAdType       AdType = 3
	RewardedAdType     AdType = 6
)

func AdTypeFromDomain(t ad.Type) AdType {
	switch t {
	case ad.InterstitialType:
		return InterstitialAdType
	case ad.BannerType:
		return BannerAdType
	case ad.RewardedType:
		return RewardedAdType
	default:
		return UnknownAdType
	}
}

func (t AdType) Domain() ad.Type {
	switch t {
	case InterstitialAdType:
		return ad.InterstitialType
	case BannerAdType:
		return ad.BannerType
	case RewardedAdType:
		return ad.RewardedType
	default:
		return ad.UnknownType
	}
}

type PlatformID int32

const (
	UnknownPlatformID PlatformID = 0
	AndroidPlatformID PlatformID = 1
	IOSPlatformID     PlatformID = 4
)