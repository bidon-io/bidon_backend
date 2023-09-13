package admin

import "context"

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
		repo:   store.Countries(),
		policy: newCountryPolicy(store),
	}
}

type CountryRepo interface {
	AllResourceQuerier[Country]
	ResourceManipulator[Country, CountryAttrs]
}

type countryPolicy struct {
	repo CountryRepo
}

func newCountryPolicy(store Store) *countryPolicy {
	return &countryPolicy{
		repo: store.Countries(),
	}
}

func (p *countryPolicy) getReadScope(_ AuthContext) resourceScope[Country] {
	return &publicResourceScope[Country]{
		repo: p.repo,
	}
}

func (p *countryPolicy) getManageScope(authCtx AuthContext) resourceScope[Country] {
	return &privateResourceScope[Country]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *countryPolicy) authorizeCreate(_ context.Context, authCtx AuthContext, _ *CountryAttrs) error {
	if !authCtx.IsAdmin() {
		return ErrActionForbidden
	}

	return nil
}

func (p *countryPolicy) authorizeUpdate(_ context.Context, _ AuthContext, _ *Country, _ *CountryAttrs) error {
	return nil
}

func (p *countryPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *Country) error {
	return nil
}
