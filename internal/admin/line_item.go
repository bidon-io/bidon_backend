package admin

import "github.com/shopspring/decimal"

type LineItem struct {
	ID int64 `json:"id"`
	LineItemAttrs
}

type LineItemAttrs struct {
	HumanName   string           `json:"human_name"`
	AppID       int64            `json:"app_id"`
	BidFloor    *decimal.Decimal `json:"bid_floor"`
	AdType      AdType           `json:"ad_type"`
	Format      *LineItemFormat  `json:"format"`
	AccountID   int64            `json:"account_id"`
	AccountType string           `json:"account_type"`
	Code        *string          `json:"code"`
	Extra       map[string]any   `json:"extra"`
}

type LineItemFormat string

const (
	EmptyLineItemFormat       LineItemFormat = ""
	BannerLineItemFormat      LineItemFormat = "BANNNER"
	LeaderboardLineItemFormat LineItemFormat = "LEADERBOARD"
	MRECLineItemFormat        LineItemFormat = "MREC"
	AdaptiveLineItemFormat    LineItemFormat = "ADAPTIVE"
)

type LineItemService = resourceService[LineItem, LineItemAttrs]
