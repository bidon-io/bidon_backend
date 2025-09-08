package admin

import "context"

//go:generate go run -mod=mod github.com/matryer/moq@v0.5.3 -out demand_source_mocks_test.go . DemandSourceRepo

const DemandSourceResourceKey = "demand_source"

type DemandSourceResource struct {
	*DemandSource
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

type DemandSource struct {
	ID int64 `json:"id"`
	DemandSourceAttrs
}

type DemandSourceAttrs struct {
	HumanName string `json:"human_name"`
	ApiKey    string `json:"api_key"`
}

type DemandSourceService struct {
	*ResourceService[DemandSourceResource, DemandSource, DemandSourceAttrs]
}

func NewDemandSourceService(store Store) *DemandSourceService {
	s := &DemandSourceService{
		ResourceService: &ResourceService[DemandSourceResource, DemandSource, DemandSourceAttrs]{},
	}

	s.resourceKey = DemandSourceResourceKey

	s.repo = store.DemandSources()
	s.policy = newDemandSourcePolicy(store)

	s.prepareResource = func(authCtx AuthContext, source *DemandSource) DemandSourceResource {
		return DemandSourceResource{
			DemandSource: source,
			Permissions:  s.policy.instancePermissions(authCtx, source),
		}
	}

	return s
}

type DemandSourceRepo interface {
	AllResourceQuerier[DemandSource]
	ResourceManipulator[DemandSource, DemandSourceAttrs]
}

type demandSourcePolicy struct {
	repo DemandSourceRepo
}

func newDemandSourcePolicy(store Store) *demandSourcePolicy {
	return &demandSourcePolicy{
		repo: store.DemandSources(),
	}
}

func (p *demandSourcePolicy) getReadScope(_ AuthContext) resourceScope[DemandSource] {
	return &publicResourceScope[DemandSource]{
		repo: p.repo,
	}
}

func (p *demandSourcePolicy) getManageScope(authCtx AuthContext) resourceScope[DemandSource] {
	return &privateResourceScope[DemandSource]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *demandSourcePolicy) authorizeCreate(_ context.Context, authCtx AuthContext, _ *DemandSourceAttrs) error {
	if !authCtx.IsAdmin() {
		return ErrActionForbidden
	}

	return nil
}

func (p *demandSourcePolicy) authorizeUpdate(_ context.Context, _ AuthContext, _ *DemandSource, _ *DemandSourceAttrs) error {
	return nil
}

func (p *demandSourcePolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *DemandSource) error {
	return nil
}

func (p *demandSourcePolicy) permissions(authCtx AuthContext) ResourcePermissions {
	return ResourcePermissions{
		Read:   true,
		Create: authCtx.IsAdmin(),
	}
}

func (p *demandSourcePolicy) instancePermissions(authCtx AuthContext, _ *DemandSource) ResourceInstancePermissions {
	return ResourceInstancePermissions{
		Update: authCtx.IsAdmin(),
		Delete: authCtx.IsAdmin(),
	}
}
