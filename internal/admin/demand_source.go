package admin

//go:generate go run -mod=mod github.com/matryer/moq@latest -out demand_source_mocks_test.go . DemandSourceRepo

type DemandSource struct {
	ID int64 `json:"id"`
	DemandSourceAttrs
}

type DemandSourceAttrs struct {
	HumanName string `json:"human_name"`
	ApiKey    string `json:"api_key"`
}

type DemandSourceService = ResourceService[DemandSource, DemandSourceAttrs]

func NewDemandSourceService(store Store) *DemandSourceService {
	return &DemandSourceService{
		repo: store.DemandSources(),
		policy: &demandSourcePolicy{
			repo: store.DemandSources(),
		},
	}
}

type DemandSourceRepo interface {
	AllResourceQuerier[DemandSource]
	ResourceManipulator[DemandSource, DemandSourceAttrs]
}

type demandSourcePolicy struct {
	repo DemandSourceRepo
}

func (p *demandSourcePolicy) scope(_ AuthContext) resourceScope[DemandSource] {
	return &publicResourceScope[DemandSource]{
		repo: p.repo,
	}
}
