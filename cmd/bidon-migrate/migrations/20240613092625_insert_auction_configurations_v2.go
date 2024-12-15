package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upInsertAuctionConfigurationsV2, downInsertAuctionConfigurationsV2)
}

func upInsertAuctionConfigurationsV2(ctx context.Context, tx *sql.Tx) error {
	type roundConfig struct {
		ID      string   `json:"id"`
		Demands []string `json:"demands"`
		Bidding []string `json:"bidding"`
		Timeout int      `json:"timeout"`
	}

	type AuctionConfiguration struct {
		ID                       int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
		Name                     sql.NullString `gorm:"column:name;type:character varying" json:"name"`
		AppID                    int64          `gorm:"column:app_id;type:bigint;not null;index:index_auction_configurations_on_app_id,priority:1" json:"app_id"`
		AdType                   int32          `gorm:"column:ad_type;type:integer;not null" json:"ad_type"`
		Rounds                   []roundConfig  `gorm:"column:rounds;type:jsonb;default:'[]';serializer:json" json:"rounds"`
		Status                   sql.NullInt32  `gorm:"column:status;type:integer" json:"status"`
		Settings                 map[string]any `gorm:"column:settings;type:jsonb;default:'{}';serializer:json" json:"settings"`
		Pricefloor               float64        `gorm:"column:pricefloor;type:double precision;not null" json:"pricefloor"`
		CreatedAt                time.Time      `gorm:"column:created_at;type:timestamp(6) without time zone;not null" json:"created_at"`
		UpdatedAt                time.Time      `gorm:"column:updated_at;type:timestamp(6) without time zone;not null" json:"updated_at"`
		SegmentID                *sql.NullInt64 `gorm:"column:segment_id;type:bigint;index:index_auction_configurations_on_segment_id,priority:1" json:"segment_id"`
		ExternalWinNotifications *bool          `gorm:"column:external_win_notifications;type:boolean;not null;default:false" json:"external_win_notifications"`
		PublicUID                sql.NullInt64  `gorm:"column:public_uid;type:bigint;uniqueIndex:index_auction_configurations_on_public_uid,priority:1" json:"public_uid"`
		Timeout                  int32          `gorm:"column:timeout;type:integer;not null" json:"timeout"`
		Demands                  pq.StringArray `gorm:"column:demands;type:character varying[];default:ARRAY[]" json:"demands"`
		Bidding                  pq.StringArray `gorm:"column:bidding;type:character varying[];default:ARRAY[]" json:"bidding"`
		AdUnitIds                pq.Int64Array  `gorm:"column:ad_unit_ids;type:bigint[];default:ARRAY[]" json:"ad_unit_ids"`
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: tx}))
	if err != nil {
		return err
	}

	var auctionConfigurations []AuctionConfiguration
	if err := gormDB.Find(&auctionConfigurations).Error; err != nil {
		return err
	}

	for _, aucConf := range auctionConfigurations {
		log.Printf("Processing auction configuration: %v", aucConf.ID)
		if len(aucConf.Rounds) == 0 {
			log.Printf("No rounds found for auction configuration: %v", aucConf.ID)
			continue
		}
		demands := aucConf.Rounds[0].Demands
		bidding := aucConf.Rounds[0].Bidding
		timeout := aucConf.Rounds[0].Timeout
		settings := aucConf.Settings
		settings["v2"] = true
		settings["reference_id"] = aucConf.ID
		snowflakeID, err := generateSnowflakeID()
		time.Sleep(100 * time.Millisecond) // Sleep for 100ms to avoid duplicate snowflake IDs
		if err != nil {
			return fmt.Errorf("generate snowflake id: %v", err)
		}
		adapters := append(demands, bidding...)

		adUnitIDs, err := fetchAdUnitIDs(gormDB, &adUnitParams{
			AppID:      aucConf.AppID,
			AdType:     aucConf.AdType,
			Adapters:   adapters,
			PriceFloor: &aucConf.Pricefloor,
		})
		if err != nil {
			log.Printf("Error fetching ad unit IDs: %v", err)
		}

		aucConfV2 := AuctionConfiguration{
			Name:       aucConf.Name,
			AppID:      aucConf.AppID,
			AdType:     aucConf.AdType,
			Status:     aucConf.Status,
			Pricefloor: aucConf.Pricefloor,
			Settings:   settings,
			Demands:    demands,
			Bidding:    bidding,
			AdUnitIds:  adUnitIDs,
			Timeout:    int32(timeout),
			UpdatedAt:  time.Now(),
			CreatedAt:  time.Now(),
			PublicUID: sql.NullInt64{
				Int64: snowflakeID,
				Valid: true,
			},
			ExternalWinNotifications: aucConf.ExternalWinNotifications,
		}

		if err = gormDB.Create(&aucConfV2).Error; err != nil {
			log.Printf("Error creating auction configuration: %v", err)
		}
	}

	return nil
}

func downInsertAuctionConfigurationsV2(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

func generateSnowflakeID() (int64, error) {
	snowflakeNodeID, _ := strconv.ParseInt("1", 10, 64)
	node, err := snowflake.NewNode(snowflakeNodeID)
	if err != nil {
		return 0, fmt.Errorf("snowflake.NewNode(%v): %v", snowflakeNodeID, err)
	}

	return node.Generate().Int64(), nil
}

type LineItem struct {
	ID int64 `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true" json:"id"`
}

type adUnitParams struct {
	AppID      int64
	AdType     int32
	AdFormat   string
	Adapters   []string
	PriceFloor *float64
}

func fetchAdUnitIDs(gormDB *gorm.DB, params *adUnitParams) ([]int64, error) {
	var pricefloor float32
	if params.PriceFloor != nil {
		pricefloor = float32(*params.PriceFloor)
	}
	sql := `
SELECT line_items.id FROM line_items
INNER JOIN demand_source_accounts AS account ON line_items.account_id = account.id
INNER JOIN demand_sources AS demand_source ON account.demand_source_id = demand_source.id
WHERE app_id = ? AND ad_type = ? AND demand_source.api_key IN (?) AND (bid_floor >= ? OR line_items.bidding)
`
	var adUnits []LineItem
	if err := gormDB.Raw(sql, params.AppID, params.AdType, params.Adapters, pricefloor).Scan(&adUnits).Error; err != nil {
		return nil, err
	}

	adUnitIDs := make([]int64, 0, len(adUnits))
	for _, adUnit := range adUnits {
		adUnitIDs = append(adUnitIDs, adUnit.ID)
	}

	return adUnitIDs, nil
}
