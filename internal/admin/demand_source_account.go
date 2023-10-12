package admin

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const DemandSourceAccountResourceKey = "demand_source_account"

type DemandSourceAccountResource struct {
	*DemandSourceAccount
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

type DemandSourceAccount struct {
	ID        int64  `json:"id"`
	PublicUID string `json:"public_uid"`
	DemandSourceAccountAttrs
	User         User         `json:"user"`
	DemandSource DemandSource `json:"demand_source"`
}

type DemandSourceAccountAttrs struct {
	UserID         int64          `json:"user_id"`
	Label          string         `json:"label"`
	Type           string         `json:"type"`
	DemandSourceID int64          `json:"demand_source_id"`
	IsBidding      *bool          `json:"is_bidding"`
	Extra          map[string]any `json:"extra"`
}

type DemandSourceAccountService struct {
	*ResourceService[DemandSourceAccountResource, DemandSourceAccount, DemandSourceAccountAttrs]
}

func NewDemandSourceAccountService(store Store) *DemandSourceAccountService {
	s := &DemandSourceAccountService{
		ResourceService: &ResourceService[DemandSourceAccountResource, DemandSourceAccount, DemandSourceAccountAttrs]{},
	}

	s.resourceKey = DemandSourceAccountResourceKey

	s.repo = store.DemandSourceAccounts()
	s.policy = newDemandSourceAccountPolicy(store)

	s.prepareResource = func(authCtx AuthContext, account *DemandSourceAccount) DemandSourceAccountResource {
		return DemandSourceAccountResource{
			DemandSourceAccount: account,
			Permissions:         s.policy.instancePermissions(authCtx, account),
		}
	}

	s.prepareCreateAttrs = func(authCtx AuthContext, attrs *DemandSourceAccountAttrs) {
		if attrs.UserID == 0 && !authCtx.IsAdmin() {
			attrs.UserID = authCtx.UserID()
		}
	}

	s.getValidator = func(attrs *DemandSourceAccountAttrs) v8n.ValidatableWithContext {
		return &demandSourceAccountValidator{
			attrs:            attrs,
			demandSourceRepo: store.DemandSources(),
		}
	}

	return s
}

//go:generate go run -mod=mod github.com/matryer/moq@latest -out demand_source_account_mocks_test.go . DemandSourceAccountRepo
type DemandSourceAccountRepo interface {
	AllResourceQuerier[DemandSourceAccount]
	OwnedResourceQuerier[DemandSourceAccount]
	OwnedOrSharedResourceQuerier[DemandSourceAccount]
	ResourceManipulator[DemandSourceAccount, DemandSourceAccountAttrs]
}

type demandSourceAccountPolicy struct {
	repo DemandSourceAccountRepo

	userPolicy         *userPolicy
	demandSourcePolicy *demandSourcePolicy
}

func newDemandSourceAccountPolicy(store Store) *demandSourceAccountPolicy {
	return &demandSourceAccountPolicy{
		repo: store.DemandSourceAccounts(),

		userPolicy:         newUserPolicy(store),
		demandSourcePolicy: newDemandSourcePolicy(store),
	}
}

func (p *demandSourceAccountPolicy) getReadScope(authCtx AuthContext) resourceScope[DemandSourceAccount] {
	return &ownedOrSharedResourceScope[DemandSourceAccount]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *demandSourceAccountPolicy) getManageScope(authCtx AuthContext) resourceScope[DemandSourceAccount] {
	return &ownedResourceScope[DemandSourceAccount]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *demandSourceAccountPolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *DemandSourceAccountAttrs) error {
	// If user is not the owner, check if user can manage the owner.
	if attrs.UserID != authCtx.UserID() {
		_, err := p.userPolicy.getManageScope(authCtx).find(ctx, attrs.UserID)
		return err
	}

	// Check if user can read the account.
	_, err := p.demandSourcePolicy.getReadScope(authCtx).find(ctx, attrs.DemandSourceID)
	return err
}

func (p *demandSourceAccountPolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, account *DemandSourceAccount, attrs *DemandSourceAccountAttrs) error {
	// If user tries to change the owner and owner is not the same as before, check if user can manage the new owner.
	if attrs.UserID != 0 && attrs.UserID != account.UserID {
		_, err := p.userPolicy.getManageScope(authCtx).find(ctx, attrs.UserID)
		return err
	}

	// If user tries to change the demand source and demand source is not the same as before, check if user can read the new demand source.
	if attrs.DemandSourceID != 0 && attrs.DemandSourceID != account.DemandSourceID {
		_, err := p.demandSourcePolicy.getReadScope(authCtx).find(ctx, attrs.DemandSourceID)
		return err
	}

	return nil
}

func (p *demandSourceAccountPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *DemandSourceAccount) error {
	return nil
}

func (p *demandSourceAccountPolicy) permissions(_ AuthContext) ResourcePermissions {
	return ResourcePermissions{
		Read:   true,
		Create: true,
	}
}

func (p *demandSourceAccountPolicy) instancePermissions(authCtx AuthContext, account *DemandSourceAccount) ResourceInstancePermissions {
	return ResourceInstancePermissions{
		Update: authCtx.IsAdmin() || authCtx.UserID() == account.UserID,
		Delete: authCtx.IsAdmin() || authCtx.UserID() == account.UserID,
	}
}

type demandSourceAccountValidator struct {
	attrs *DemandSourceAccountAttrs

	demandSourceRepo DemandSourceRepo
}

func (v *demandSourceAccountValidator) ValidateWithContext(ctx context.Context) error {
	demandSource, err := v.demandSourceRepo.Find(ctx, v.attrs.DemandSourceID)
	if err != nil {
		return v8n.NewInternalError(err)
	}

	return v8n.ValidateStruct(v.attrs,
		v8n.Field(&v.attrs.Extra, v.extraRule(demandSource)),
	)
}

func (v *demandSourceAccountValidator) extraRule(demandSource *DemandSource) v8n.Rule {
	var rule v8n.MapRule

	switch adapter.Key(demandSource.ApiKey) {
	case adapter.AmazonKey:
		rule = v8n.Map(v8n.Key("price_points_map", isMap, v8n.Each(v8n.Map(
			v8n.Key("name", v8n.Required, isString),
			v8n.Key("price_point", v8n.Required, isString),
			v8n.Key("price", v8n.Required, isFloat),
		))))
	case adapter.ApplovinKey:
		rule = v8n.Map(
			v8n.Key("sdk_key", v8n.Required, isString),
		)
	case adapter.BidmachineKey:
		rule = v8n.Map(
			v8n.Key("seller_id", v8n.Required, isString),
			v8n.Key("endpoint", v8n.Required, Any(is.URL, is.Host)),
			v8n.Key("mediation_config", v8n.NotNil, v8n.Each(isString)),
		)
	case adapter.BigoAdsKey:
		rule = v8n.Map(
			v8n.Key("publisher_id", v8n.Required, isString),
			v8n.Key("endpoint", v8n.Required, is.URL),
		)
	case adapter.InmobiKey:
		rule = v8n.Map(
			v8n.Key("account_id", v8n.Required, isString),
		)
	case adapter.MintegralKey:
		rule = v8n.Map(
			v8n.Key("app_key", v8n.Required, isString),
			v8n.Key("publisher_id", v8n.Required, isString),
		)
	case adapter.MobileFuseKey:
		rule = v8n.Map(
			v8n.Key("publisher_id", v8n.Required, isString),
		)
	case adapter.VungleKey:
		rule = v8n.Map(
			v8n.Key("account_id", v8n.Required, isString),
		)
	}

	return rule.AllowExtraKeys()
}
