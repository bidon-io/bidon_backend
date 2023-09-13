package admin

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

type AppDemandProfile struct {
	ID int64 `json:"id"`
	AppDemandProfileAttrs
	App          App                 `json:"app"`
	Account      DemandSourceAccount `json:"account"`
	DemandSource DemandSource        `json:"demand_source"`
}

type AppDemandProfileAttrs struct {
	AppID          int64          `json:"app_id"`
	DemandSourceID int64          `json:"demand_source_id"`
	AccountID      int64          `json:"account_id"`
	Data           map[string]any `json:"data"`
	AccountType    string         `json:"account_type"`
}

type AppDemandProfileService = ResourceService[AppDemandProfile, AppDemandProfileAttrs]

func NewAppDemandProfileService(store Store) *AppDemandProfileService {
	s := &AppDemandProfileService{
		repo: store.AppDemandProfiles(),
	}

	s.policy = newAppDemandProfilePolicy(store)

	s.getValidator = func(attrs *AppDemandProfileAttrs) v8n.ValidatableWithContext {
		return &appDemandProfileAttrsValidator{
			attrs:            attrs,
			demandSourceRepo: store.DemandSources(),
		}
	}

	return s
}

type AppDemandProfileRepo interface {
	AllResourceQuerier[AppDemandProfile]
	OwnedResourceQuerier[AppDemandProfile]
	ResourceManipulator[AppDemandProfile, AppDemandProfileAttrs]
}

type appDemandProfilePolicy struct {
	repo AppDemandProfileRepo

	appPolicy                 *appPolicy
	demandSourceAccountPolicy *demandSourceAccountPolicy
	demandSourcePolicy        *demandSourcePolicy
}

func newAppDemandProfilePolicy(store Store) *appDemandProfilePolicy {
	return &appDemandProfilePolicy{
		repo: store.AppDemandProfiles(),

		appPolicy:                 newAppPolicy(store),
		demandSourceAccountPolicy: newDemandSourceAccountPolicy(store),
		demandSourcePolicy:        newDemandSourcePolicy(store),
	}
}

func (p *appDemandProfilePolicy) getReadScope(authCtx AuthContext) resourceScope[AppDemandProfile] {
	return &ownedResourceScope[AppDemandProfile]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *appDemandProfilePolicy) getManageScope(authCtx AuthContext) resourceScope[AppDemandProfile] {
	return &ownedResourceScope[AppDemandProfile]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *appDemandProfilePolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *AppDemandProfileAttrs) error {
	// Check if user can manage the app.
	_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
	if err != nil {
		return err
	}

	// Check if user can read the account.
	_, err = p.demandSourceAccountPolicy.getReadScope(authCtx).find(ctx, attrs.AccountID)
	if err != nil {
		return err
	}

	// Check if user can read the demand source.
	_, err = p.demandSourcePolicy.getReadScope(authCtx).find(ctx, attrs.DemandSourceID)
	if err != nil {
		return err
	}

	return nil
}

func (p *appDemandProfilePolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, profile *AppDemandProfile, attrs *AppDemandProfileAttrs) error {
	// If user tries to change the app and app is not the same as before, check if user can manage the new app.
	if attrs.AppID != 0 && attrs.AppID != profile.AppID {
		_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
		if err != nil {
			return err
		}
	}

	// If user tries to change the account and account is not the same as before, check if user can read the new account.
	if attrs.AccountID != 0 && attrs.AccountID != profile.AccountID {
		_, err := p.demandSourceAccountPolicy.getReadScope(authCtx).find(ctx, attrs.AccountID)
		if err != nil {
			return err
		}
	}

	// If user tries to change the demand source and demand source is not the same as before, check if user can read the new demand source.
	if attrs.DemandSourceID != 0 && attrs.DemandSourceID != profile.DemandSourceID {
		_, err := p.demandSourcePolicy.getReadScope(authCtx).find(ctx, attrs.DemandSourceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *appDemandProfilePolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *AppDemandProfile) error {
	return nil
}

type appDemandProfileAttrsValidator struct {
	attrs *AppDemandProfileAttrs

	demandSourceRepo DemandSourceRepo
}

func (v *appDemandProfileAttrsValidator) ValidateWithContext(ctx context.Context) error {
	demandSource, err := v.demandSourceRepo.Find(ctx, v.attrs.DemandSourceID)
	if err != nil {
		return v8n.NewInternalError(err)
	}

	return v8n.ValidateStruct(v.attrs,
		v8n.Field(&v.attrs.Data, v.dataRule(demandSource)),
	)
}

func (v *appDemandProfileAttrsValidator) dataRule(demandSource *DemandSource) v8n.Rule {
	var rule v8n.MapRule

	switch adapter.Key(demandSource.ApiKey) {
	case adapter.AdmobKey, adapter.BigoAdsKey, adapter.DTExchangeKey, adapter.MintegralKey, adapter.VungleKey:
		rule = v8n.Map(
			v8n.Key("app_id", v8n.Required, isString),
		)
	case adapter.MetaKey:
		rule = v8n.Map(
			v8n.Key("app_id", v8n.Required, isString),
			v8n.Key("app_secret", v8n.Required, isString),
		)
	case adapter.UnityAdsKey:
		rule = v8n.Map(
			v8n.Key("game_id", v8n.Required, isString),
		)
	}

	return rule.AllowExtraKeys()
}
