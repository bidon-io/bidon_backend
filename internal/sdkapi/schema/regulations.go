package schema

type Regulations struct {
	COPPA     bool           `json:"coppa"`
	GDPR      bool           `json:"gdpr"`
	USPrivacy string         `json:"us_privacy"`
	EUPrivacy string         `json:"eu_privacy"`
	IAB       map[string]any `json:"iab"`
}
