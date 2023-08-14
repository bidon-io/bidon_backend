package admin

type DemandSource struct {
	ID int64 `json:"id"`
	DemandSourceAttrs
}

type DemandSourceAttrs struct {
	HumanName string `json:"human_name"`
	ApiKey    string `json:"api_key"`
}

type DemandSourceRepo = ResourceRepo[DemandSource, DemandSourceAttrs]

type DemandSourceService = ResourceService[DemandSource, DemandSourceAttrs]

func NewDemandSourceService(store Store) *DemandSourceService {
	return &DemandSourceService{
		ResourceRepo: store.DemandSources(),
	}
}
