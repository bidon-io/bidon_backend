package schema

type LossRequest struct {
	ShowRequest
	ExternalWinner ExternalWinner `json:"external_winner"`
}

type ExternalWinner struct {
	DemandID string  `json:"demand_id"`
	ECPM     float64 `json:"ecpm" validate:"required"`
}
