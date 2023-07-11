package schema

type ConfigRequest struct {
	BaseRequest
	Adapters Adapters `json:"adapters"`
}

func (r *ConfigRequest) Map() map[string]any {
	m := r.BaseRequest.Map()

	for key, adapter := range r.Adapters {
		m[string(key)] = adapter.Map()
	}

	return m
}
