package schema

type Regulations struct {
	COPPA bool `json:"coppa"`
	GDPR  bool `json:"gdpr"`
}

func (r Regulations) Map() map[string]any {
	m := map[string]any{
		"coppa": r.COPPA,
		"gdpr":  r.GDPR,
	}

	return m
}
