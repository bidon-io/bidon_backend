package store

import (
	"context"
	"database/sql"

	"github.com/bidon-io/bidon-backend/internal/auction"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type auctionConfiguration struct {
	Model
	Name       sql.NullString               `gorm:"column:name;type:varchar"`
	AppID      int64                        `gorm:"column:app_id;type:bigint;not null"`
	AdType     int32                        `gorm:"column:ad_type;type:integer;not null"`
	Rounds     []auction.RoundConfiguration `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json"`
	Pricefloor float64                      `gorm:"column:pricefloor;type:double precision;not null"`
}

func fromAuctionConfigurationAttrs(config *auction.ConfigurationAttrs) auctionConfiguration {
	name := sql.NullString{}
	if config.Name != "" {
		name.String = config.Name
		name.Valid = true
	}

	return auctionConfiguration{
		Name:       name,
		AppID:      config.AppID,
		AdType:     dbAdType(config.AdType),
		Rounds:     config.Rounds,
		Pricefloor: config.Pricefloor,
	}
}

func (a *auctionConfiguration) auctionConfiguration() auction.Configuration {
	return auction.Configuration{
		ID: a.ID,
		ConfigurationAttrs: auction.ConfigurationAttrs{
			Name:       a.Name.String,
			AppID:      a.AppID,
			AdType:     adType(a.AdType),
			Rounds:     a.Rounds,
			Pricefloor: a.Pricefloor,
		},
	}
}

type AuctionConfigurationRepo struct {
	DB *gorm.DB
}

func (r *AuctionConfigurationRepo) List(ctx context.Context) ([]auction.Configuration, error) {
	var dbConfigs []auctionConfiguration
	if err := r.DB.WithContext(ctx).Find(&dbConfigs).Error; err != nil {
		return nil, err
	}

	configs := make([]auction.Configuration, len(dbConfigs))
	for i, config := range dbConfigs {
		configs[i] = config.auctionConfiguration()
	}

	return configs, nil
}

func (r *AuctionConfigurationRepo) Find(ctx context.Context, id int64) (*auction.Configuration, error) {
	var dbConfig auctionConfiguration
	if err := r.DB.WithContext(ctx).First(&dbConfig, id).Error; err != nil {
		return nil, err
	}

	config := dbConfig.auctionConfiguration()
	return &config, nil
}

func (r *AuctionConfigurationRepo) Create(ctx context.Context, attrs *auction.ConfigurationAttrs) (*auction.Configuration, error) {
	dbConfig := fromAuctionConfigurationAttrs(attrs)
	if err := r.DB.WithContext(ctx).Create(&dbConfig).Error; err != nil {
		return nil, err
	}

	config := dbConfig.auctionConfiguration()
	return &config, nil
}

func (r *AuctionConfigurationRepo) Update(ctx context.Context, id int64, attrs *auction.ConfigurationAttrs) (*auction.Configuration, error) {
	dbConfig := fromAuctionConfigurationAttrs(attrs)
	dbConfig.ID = id

	if err := r.DB.WithContext(ctx).Model(&dbConfig).Clauses(clause.Returning{}).Updates(&dbConfig).Error; err != nil {
		return nil, err
	}

	config := dbConfig.auctionConfiguration()
	return &config, nil
}

func (r *AuctionConfigurationRepo) Delete(ctx context.Context, id int64) error {
	return r.DB.WithContext(ctx).Delete(&auctionConfiguration{}, id).Error
}
