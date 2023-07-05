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
