package schema

type Segment struct {
	ID  string `json:"id"`
	Ext string `json:"ext"`
}

func (s Segment) Map() map[string]any {
	m := map[string]any{
		"id": s.ID,
	}

	if s.Ext != "" {
		m["ext"] = s.Ext
	}

	return m
}