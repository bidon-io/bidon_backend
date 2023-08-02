package admin

type AppDemandProfile struct {
	ID int64 `json:"id"`
	AppDemandProfileAttrs
	App          App                 `json:"app"`
	Account      DemandSourceAccount `json:"account"`
	DemandSource DemandSource        `json:"demand_source"`
}

type AppDemandProfileAttrs struct {
	AppID          int64          `json:"app_id"`
	DemandSourceID int64          `json:"demand_source_id"`
	AccountID      int64          `json:"account_id"`
	Data           map[string]any `json:"data"`
	AccountType    string         `json:"account_type"`
}

type AppDemandProfileService = resourceService[AppDemandProfile, AppDemandProfileAttrs]
