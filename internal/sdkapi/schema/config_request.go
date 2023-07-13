package schema

type ConfigRequest struct {
	BaseRequest
	Adapters Adapters `json:"adapters"`
}

func (r *ConfigRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	if r.Adapters != nil {
		m["adapters"] = r.Adapters.Map()
	}

	return m
}
