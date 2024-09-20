package schema

type LossRequest struct {
	ShowRequest
	ExternalWinner ExternalWinner `json:"external_winner"`
}

type ExternalWinner struct {
	DemandID string   `json:"demand_id"`
	ECPM     *float64 `json:"ecpm"` // Deprecated: ECPM is deprecated since 0.7, use Price instead
	Price    *float64 `json:"price"`
}

func (e *ExternalWinner) GetPrice() float64 {
	if e.Price != nil {
		return *e.Price
	}

	if e.ECPM != nil {
		return *e.ECPM
	}

	return 0
}
