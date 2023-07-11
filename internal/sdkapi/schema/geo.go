package schema

type Geo struct {
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Accuracy  float64 `json:"accuracy"`
	LastFix   int     `json:"lastfix"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	ZIP       string  `json:"zip"`
	UTCOffset int     `json:"utcoffset"`
}

func (g Geo) Map() map[string]any {
	m := map[string]any{
		"lat":       g.Lat,
		"lon":       g.Lon,
		"accuracy":  g.Accuracy,
		"lastfix":   g.LastFix,
		"country":   g.Country,
		"city":      g.City,
		"zip":       g.ZIP,
		"utcoffset": g.UTCOffset,
	}

	return m
}
