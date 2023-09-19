package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"github.com/bwmarrin/snowflake"
	"github.com/shopspring/decimal"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
	snowflakeNode *snowflake.Node
}

type Option func(*DB)

func WithSnowflakeNode(node *snowflake.Node) Option {
	return func(db *DB) {
		db.snowflakeNode = node
	}
}

func Open(databaseURL string, opts ...Option) (*DB, error) {
	gormDB, err := gorm.Open(postgres.Open(databaseURL))
	if err != nil {
		return nil, err
	}

	err = gormDB.Use(otelgorm.NewPlugin())
	if err != nil {
		return nil, err
	}

	db := &DB{DB: gormDB}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}

func (db *DB) Begin(opts ...*sql.TxOptions) *DB {
	return &DB{DB: db.DB.Begin(opts...)}
}

func (db *DB) SetDebug() {
	db.Logger = db.Logger.LogMode(logger.Info)
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

func (db *DB) GenerateSnowflakeID() int64 {
	if db.snowflakeNode == nil {
		return 0
	}

	return db.snowflakeNode.Generate().Int64()
}

// Model is different from default gorm.Model, because we already have schema from Rails
type Model struct {
	ID        int64     `gorm:"primaryKey;column:id;type:bigint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp(6);not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp(6);not null"`
}

type AppDemandProfile struct {
	Model
	AppID          int64               `gorm:"column:app_id;type:bigint;not null"`
	App            App                 `gorm:"foreignKey:AppID"`
	AccountType    string              `gorm:"column:account_type;type:varchar;not null"`
	AccountID      int64               `gorm:"column:account_id;type:bigint;not null"`
	Account        DemandSourceAccount `gorm:"foreignKey:AccountID"`
	DemandSourceID int64               `gorm:"column:demand_source_id;type:bigint;not null"`
	DemandSource   DemandSource        `gorm:"foreignKey:DemandSourceID"`
	Data           datatypes.JSON      `gorm:"column:data;type:jsonb;default:'{}'"`
	PublicUID      sql.NullInt64       `gorm:"column:public_uid;type:bigint"`
}

type App struct {
	Model
	UserID      int64          `gorm:"column:user_id;type:bigint;not null"`
	User        User           `gorm:"foreignKey:UserID"`
	PlatformID  PlatformID     `gorm:"column:platform_id;type:integer;not null"`
	HumanName   string         `gorm:"column:human_name;type:varchar;not null"`
	PackageName sql.NullString `gorm:"column:package_name;type:varchar"`
	AppKey      sql.NullString `gorm:"column:app_key;type:varchar"`
	Settings    map[string]any `gorm:"column:settings;type:jsonb;default:'{}';serializer:json"`
	PublicUID   sql.NullInt64  `gorm:"column:public_uid;type:bigint"`
}

type AuctionConfiguration struct {
	Model
	Name                     sql.NullString        `gorm:"column:name;type:varchar"`
	AppID                    int64                 `gorm:"column:app_id;type:bigint;not null"`
	App                      App                   `gorm:"foreignKey:AppID"`
	AdType                   AdType                `gorm:"column:ad_type;type:integer;not null"`
	Rounds                   []auction.RoundConfig `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json"`
	Pricefloor               float64               `gorm:"column:pricefloor;type:numeric;not null"`
	SegmentID                *sql.NullInt64        `gorm:"column:segment_id;type:bigint"`
	Segment                  *Segment              `gorm:"foreignKey:SegmentID"`
	ExternalWinNotifications *bool                 `gorm:"column:external_win_notifications;type:boolean;default:false;not null"`
	PublicUID                sql.NullInt64         `gorm:"column:public_uid;type:bigint"`
}

type Country struct {
	Model
	Alpha2Code string         `gorm:"column:alpha_2_code;type:varchar;not null"`
	Alpha3Code string         `gorm:"column:alpha_3_code;type:varchar;not null"`
	HumanName  sql.NullString `gorm:"column:human_name;type:varchar"`
}

type DemandSourceAccount struct {
	Model
	DemandSourceID int64          `gorm:"column:demand_source_id;type:bigint;not null"`
	DemandSource   DemandSource   `gorm:"foreignKey:DemandSourceID"`
	UserID         int64          `gorm:"column:user_id;type:bigint;not null"`
	User           User           `gorm:"foreignKey:UserID"`
	Label          string         `gorm:"column:label;type:varchar;not null"`
	Type           string         `gorm:"column:type;type:varchar;not null"`
	Extra          datatypes.JSON `gorm:"column:extra;type:jsonb;default:'{}'"`
	IsBidding      *bool          `gorm:"column:bidding;type:boolean;default:false"`
	IsDefault      sql.NullBool   `gorm:"column:is_default;type:boolean"`
	PublicUID      sql.NullInt64  `gorm:"column:public_uid;type:bigint"`
}

type DemandSource struct {
	Model
	APIKey    string `gorm:"column:api_key;type:varchar;not null"`
	HumanName string `gorm:"column:human_name;type:varchar;not null"`
}

type LineItem struct {
	Model
	AppID       int64               `gorm:"column:app_id;type:bigint;not null"`
	App         App                 `gorm:"foreignKey:AppID"`
	AccountType string              `gorm:"column:account_type;type:varchar;not null"`
	AccountID   int64               `gorm:"column:account_id;type:bigint;not null"`
	Account     DemandSourceAccount `gorm:"foreignKey:AccountID"`
	HumanName   string              `gorm:"column:human_name;type:varchar;not null"`
	Code        *string             `gorm:"column:code;type:varchar;not null"`
	BidFloor    decimal.NullDecimal `gorm:"column:bid_floor;type:numeric"`
	AdType      AdType              `gorm:"column:ad_type;type:integer;not null"`
	Extra       map[string]any      `gorm:"column:extra;type:jsonb;default:'{}';serializer:json"`
	Width       int32               `gorm:"column:width;type:integer;default:0;not null"`
	Height      int32               `gorm:"column:height;type:integer;default:0;not null"`
	Format      sql.NullString      `gorm:"column:format;type:varchar"`
	PublicUID   sql.NullInt64       `gorm:"column:public_uid;type:bigint"`
}

type Segment struct {
	Model
	Name        string           `gorm:"column:name;type:varchar;not null"`
	Description string           `gorm:"column:description;type:text;not null"`
	Filters     []segment.Filter `gorm:"column:filters;type:jsonb;not null;default:'[]';serializer:json"`
	Enabled     *bool            `gorm:"column:enabled;type:bool;not null;default:true"`
	AppID       int64            `gorm:"column:app_id;type:bigint;not null"`
	App         App              `gorm:"foreignKey:AppID"`
	Priority    int32            `gorm:"column:priority;type:integer;default:0;not null"`
	PublicUID   sql.NullInt64    `gorm:"column:public_uid;type:bigint"`
}

type User struct {
	Model
	Email        string        `gorm:"column:email;type:varchar;not null"`
	IsAdmin      *bool         `gorm:"column:is_admin;type:boolean;default:false;not null"`
	PasswordHash string        `gorm:"column:password_hash;type:varchar;not null"`
	PublicUID    sql.NullInt64 `gorm:"column:public_uid;type:bigint"`
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

func (t *AdType) Scan(v any) (err error) {
	if v, ok := v.(int64); ok {
		*t = AdType(v)
		return nil
	}

	return fmt.Errorf("db: unsupported value %v (type %T) converting to AdType", v, v)
}

func (t AdType) Value() (driver.Value, error) {
	return int64(t), nil
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

func (id *PlatformID) Scan(v any) (err error) {
	if v, ok := v.(int64); ok {
		*id = PlatformID(v)
		return nil
	}

	return fmt.Errorf("db: unsupported value %v (type %T) converting to PlatformID", v, v)
}

func (id PlatformID) Value() (driver.Value, error) {
	return int64(id), nil
}
