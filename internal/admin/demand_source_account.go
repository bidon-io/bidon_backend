package admin

type DemandSourceAccount struct {
	ID int64 `json:"id"`
	DemandSourceAccountAttrs
	User         User         `json:"user"`
	DemandSource DemandSource `json:"demand_source"`
}

type DemandSourceAccountAttrs struct {
	UserID         int64             `json:"user_id"`
	Type           string            `json:"type"`
	DemandSourceID int64             `json:"demand_source_id"`
	IsBidding      *bool             `json:"is_bidding"`
	Extra          map[string]string `json:"extra"`
}

type DemandSourceAccountService = resourceService[DemandSourceAccount, DemandSourceAccountAttrs]
