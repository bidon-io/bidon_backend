package schema

type Regulations struct {
	COPPA     bool           `json:"coppa"`
	GDPR      bool           `json:"gdpr"`
	USPrivacy string         `json:"us_privacy"`
	EUPrivacy string         `json:"eu_privacy"`
	IAB       map[string]any `json:"iab"`
}

func (r Regulations) Map() map[string]any {
	m := map[string]any{
		"coppa":      r.COPPA,
		"gdpr":       r.GDPR,
		"us_privacy": r.USPrivacy,
		"eu_privacy": r.EUPrivacy,
		"iab":        r.IAB,
	}

	return m
}
