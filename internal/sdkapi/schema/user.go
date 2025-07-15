package schema

type User struct {
	IDFA                        string         `json:"idfa" validate:"uuid"`
	TrackingAuthorizationStatus string         `json:"tracking_authorization_status" validate:"required"`
	IDFV                        string         `json:"idfv" validate:"omitempty,uuid"`
	IDG                         string         `json:"idg" validate:"uuid"`
	Consent                     map[string]any `json:"consent"`
	COPPA                       *bool          `json:"coppa"`
	AppSetID                    string         `json:"app_set_id"`
	AppSetIDScope               string         `json:"app_set_id_scope" validate:"omitempty,oneof=app developer"`
}
