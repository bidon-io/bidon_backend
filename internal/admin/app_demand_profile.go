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

	s.policy = &appDemandProfilePolicy{
		repo: store.AppDemandProfiles(),
	}

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
}

func (p *appDemandProfilePolicy) scope(authCtx AuthContext) resourceScope[AppDemandProfile] {
	return &ownedResourceScope[AppDemandProfile]{
		repo:    p.repo,
		authCtx: authCtx,
	}
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
