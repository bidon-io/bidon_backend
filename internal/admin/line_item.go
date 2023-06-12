package admin

import (
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/shopspring/decimal"
)

type LineItem struct {
	ID int64 `json:"id"`
	LineItemAttrs
}

type LineItemAttrs struct {
	HumanName   string           `json:"human_name"`
	AppID       int64            `json:"app_id"`
	BidFloor    *decimal.Decimal `json:"bid_floor"`
	AdType      ad.Type          `json:"ad_type"`
	Format      *ad.Format       `json:"format"`
	AccountID   int64            `json:"account_id"`
	AccountType string           `json:"account_type"`
	Code        *string          `json:"code"`
	Extra       map[string]any   `json:"extra"`
}

type LineItemService = resourceService[LineItem, LineItemAttrs]
