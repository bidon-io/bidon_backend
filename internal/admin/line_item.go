package admin

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/shopspring/decimal"
)

type LineItem struct {
	ID        int64  `json:"id"`
	PublicUID string `json:"public_uid"`
	LineItemAttrs
	App     App                 `json:"app"`
	Account DemandSourceAccount `json:"account"`
}

type LineItemAttrs struct {
	HumanName   string           `json:"human_name"`
	AppID       int64            `json:"app_id"`
	BidFloor    *decimal.Decimal `json:"bid_floor"`
	AdType      ad.Type          `json:"ad_type"`
	Format      *ad.Format       `json:"format"`
	AccountID   int64            `json:"account_id"`
	AccountType string           `json:"account_type"`
	Code        *string          `json:"code"`
	IsBidding   *bool            `json:"is_bidding"`
	Extra       map[string]any   `json:"extra"`
}

type LineItemService = ResourceService[LineItem, LineItemAttrs]

func NewLineItemService(store Store) *LineItemService {
	s := &LineItemService{
		repo: store.LineItems(),
	}

	s.policy = newLineItemPolicy(store)

	s.getValidator = func(attrs *LineItemAttrs) v8n.ValidatableWithContext {
		return &lineItemAttrsValidator{
			attrs:                   attrs,
			demandSourceAccountRepo: store.DemandSourceAccounts(),
		}
	}

	return s
}

type LineItemRepo interface {
	AllResourceQuerier[LineItem]
	OwnedResourceQuerier[LineItem]
	ResourceManipulator[LineItem, LineItemAttrs]
}

type lineItemPolicy struct {
	repo LineItemRepo

	appPolicy                 *appPolicy
	demandSourceAccountPolicy *demandSourceAccountPolicy
}

func newLineItemPolicy(store Store) *lineItemPolicy {
	return &lineItemPolicy{
		repo: store.LineItems(),

		appPolicy:                 newAppPolicy(store),
		demandSourceAccountPolicy: newDemandSourceAccountPolicy(store),
	}
}

func (p *lineItemPolicy) getReadScope(authCtx AuthContext) resourceScope[LineItem] {
	return &ownedResourceScope[LineItem]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *lineItemPolicy) getManageScope(authCtx AuthContext) resourceScope[LineItem] {
	return &ownedResourceScope[LineItem]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *lineItemPolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *LineItemAttrs) error {
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

	return nil
}

func (p *lineItemPolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, profile *LineItem, attrs *LineItemAttrs) error {
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

	return nil
}

func (p *lineItemPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *LineItem) error {
	return nil
}

type lineItemAttrsValidator struct {
	attrs *LineItemAttrs

	demandSourceAccountRepo DemandSourceAccountRepo
}

func (v *lineItemAttrsValidator) ValidateWithContext(ctx context.Context) error {
	account, err := v.demandSourceAccountRepo.Find(ctx, v.attrs.AccountID)
	if err != nil {
		return v8n.NewInternalError(err)
	}

	return v8n.ValidateStruct(v.attrs,
		v8n.Field(&v.attrs.Extra, v.extraRule(account)),
	)
}

func (v *lineItemAttrsValidator) extraRule(account *DemandSourceAccount) v8n.Rule {
	var rule v8n.MapRule

	switch adapter.Key(account.DemandSource.ApiKey) {
	case adapter.AdmobKey:
		rule = v8n.Map(
			v8n.Key("ad_unit_id", v8n.Required, isString),
		)
	case adapter.ApplovinKey:
		rule = v8n.Map(
			v8n.Key("zone_id", v8n.Required, isString),
		)
	case adapter.BigoAdsKey:
		rule = v8n.Map(
			v8n.Key("slot_id", v8n.Required, isString),
		)
	case adapter.DTExchangeKey, adapter.MetaKey, adapter.UnityAdsKey, adapter.VungleKey, adapter.MobileFuseKey:
		rule = v8n.Map(
			v8n.Key("placement_id", v8n.Required, isString),
		)
	case adapter.MintegralKey:
		rule = v8n.Map(
			v8n.Key("placement_id", v8n.Required, isString),
			v8n.Key("ad_unit_id", v8n.Required, isString),
		)
	}

	return rule.AllowExtraKeys()
}
