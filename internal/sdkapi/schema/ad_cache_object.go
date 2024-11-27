package schema

type AdCacheObject struct {
	DemandID  string  `json:"demand_id" validate:"required"`
	Timestamp int64   `json:"timestamp" validate:"required"`
	Price     float64 `json:"price" validate:"required"`
	AdUnitID  *string `json:"ad_unit_id,omitempty"`
}
