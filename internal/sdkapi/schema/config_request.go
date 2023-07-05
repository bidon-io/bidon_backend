package schema

type ConfigRequest struct {
	BaseRequest
	Adapters Adapters `json:"adapters"`
}
