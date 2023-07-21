package schema

type User struct {
	IDFA                        string         `json:"idfa" validate:"uuid"`
	TrackingAuthorizationStatus string         `json:"tracking_authorization_status" validate:"required"`
	IDFV                        string         `json:"idfv" validate:"omitempty,uuid"`
	IDG                         string         `json:"idg" validate:"uuid"`
	Consent                     map[string]any `json:"consent"`
	COPPA                       *bool          `json:"coppa"`
}

func (u User) Map() map[string]any {
	m := map[string]any{
		"idfa":                          u.IDFA,
		"tracking_authorization_status": u.TrackingAuthorizationStatus,
		"idfv":                          u.IDFV,
		"idg":                           u.IDG,
	}

	if u.Consent != nil {
		m["consent"] = u.Consent
	}

	if u.COPPA != nil {
		m["coppa"] = u.COPPA
	}

	return m
}
