package schema

type LossRequest struct {
	ShowRequest
	ExternalWinnder ExternalWinner `json:"external_winner"`
}

type ExternalWinner struct {
	DemandID string  `json:"demand_id" validate:"required"`
	ECPM     float64 `json:"ecpm" validate:"required"`
}

func (w *ExternalWinner) Map() map[string]any {
	m := map[string]any{
		"demand_id": w.DemandID,
		"ecpm":      w.ECPM,
	}

	return m
}

func (r *LossRequest) Map() map[string]any {
	m := r.ShowRequest.Map()

	m["external_winner"] = r.ExternalWinnder.Map()

	return m
}
