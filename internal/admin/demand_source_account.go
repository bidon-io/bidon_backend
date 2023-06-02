package admin

type DemandSourceAccount struct {
	ID int64 `json:"id"`
	DemandSourceAccountAttrs
}

type DemandSourceAccountAttrs struct {
	UserID         int64          `json:"user_id"`
	Type           string         `json:"type"`
	DemandSourceID int64          `json:"demand_source_id"`
	IsBidding      *bool          `json:"is_bidding"`
	Extra          map[string]any `json:"extra"`
}

type DemandSourceAccountService = resourceService[DemandSourceAccount, DemandSourceAccountAttrs]
