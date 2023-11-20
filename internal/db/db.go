package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bwmarrin/snowflake"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	*gorm.DB
}

type Option func(*DB)

func WithSnowflakeNode(node *snowflake.Node) Option {
	return func(db *DB) {
		db.DB = db.Set(snowflakeNodeKey, node)
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

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(db)
		}
		// Mark as safe to share after applying options
		// Calling Session() is important
		// For more information refer to https://gorm.io/docs/method_chaining.html#New-Session-Method
		db.DB = db.Session(&gorm.Session{})
	}

	return db, nil
}

const snowflakeNodeKey = "snowflake:node"

func generateSnowflakeID(db *gorm.DB) (int64, error) {
	node, ok := db.Get(snowflakeNodeKey)
	if !ok {
		return 0, fmt.Errorf("snowflake node not set")
	}
	if node, ok := node.(*snowflake.Node); ok {
		return node.Generate().Int64(), nil
	}
	return 0, fmt.Errorf("snowflake node has wrong type")
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
