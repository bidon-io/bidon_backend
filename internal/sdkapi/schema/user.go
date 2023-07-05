package schema

type User struct {
	IDFA                        string         `json:"idfa" validate:"uuid_rfc4122"`
	TrackingAuthorizationStatus string         `json:"tracking_authorization_status" validate:"required"`
	IDFV                        string         `json:"idfv" validate:"omitempty,uuid_rfc4122"`
	IDG                         string         `json:"idg" validate:"uuid_rfc4122"`
	Consent                     map[string]any `json:"consent"`
	COPPA                       *bool          `json:"coppa"`
}
