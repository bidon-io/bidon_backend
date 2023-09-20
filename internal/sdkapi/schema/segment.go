package schema

import "strconv"

type Segment struct {
	ID  string `json:"id"`
	UID string `json:"uid"`
	Ext string `json:"ext"`
}

func (s Segment) Map() map[string]any {
	uid, err := strconv.Atoi(s.UID)
	if err != nil {
		uid = 0
	}

	m := map[string]any{
		"id":  s.ID,
		"uid": uid,
	}

	if s.Ext != "" {
		m["ext"] = s.Ext
	}

	return m
}
