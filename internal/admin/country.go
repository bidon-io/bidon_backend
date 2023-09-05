package admin

type Country struct {
	ID int64 `json:"id"`
	CountryAttrs
}

type CountryAttrs struct {
	HumanName  string `json:"human_name"`
	Alpha2Code string `json:"alpha2_code"`
	Alpha3Code string `json:"alpha3_code"`
}

type CountryService = ResourceService[Country, CountryAttrs]

func NewCountryService(store Store) *CountryService {
	return &CountryService{
		repo: store.Countries(),
		policy: &countryPolicy{
			repo: store.Countries(),
		},
	}
}

type CountryRepo interface {
	AllResourceQuerier[Country]
	ResourceManipulator[Country, CountryAttrs]
}

type countryPolicy struct {
	repo CountryRepo
}

func (p *countryPolicy) scope(_ AuthContext) resourceScope[Country] {
	return &publicResourceScope[Country]{
		repo: p.repo,
	}
}
