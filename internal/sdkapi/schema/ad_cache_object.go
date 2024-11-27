package schema

type AdCacheObject struct {
	DemandID  string  `json:"demand_id" validate:"required"`
	Timestamp int64   `json:"timestamp" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
	AdUnitUID *string `json:"ad_unit_uid,omitempty"`
}
