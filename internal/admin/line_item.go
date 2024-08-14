package admin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
	v8n "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jszwec/csvutil"
	"github.com/shopspring/decimal"
)

const LineItemResourceKey = "line_item"

type LineItemResource struct {
	*LineItem
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

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
	IsBidding   *bool            `json:"is_bidding"`
	Extra       map[string]any   `json:"extra"`
}

type LineItemService struct {
	*ResourceService[LineItemResource, LineItem, LineItemAttrs]
	store Store
}

func NewLineItemService(store Store) *LineItemService {
	s := &LineItemService{
		ResourceService: &ResourceService[LineItemResource, LineItem, LineItemAttrs]{},
	}

	s.store = store

	s.resourceKey = LineItemResourceKey

	s.repo = store.LineItems()
	s.policy = newLineItemPolicy(store)

	s.prepareResource = func(authCtx AuthContext, lineItem *LineItem) LineItemResource {
		return LineItemResource{
			LineItem:    lineItem,
			Permissions: s.policy.instancePermissions(authCtx, lineItem),
		}
	}

	s.getValidator = func(attrs *LineItemAttrs) v8n.ValidatableWithContext {
		return &lineItemAttrsValidator{
			attrs:                   attrs,
			demandSourceAccountRepo: store.DemandSourceAccounts(),
		}
	}

	return s
}

type LineItemImportCSVAttrs struct {
	AppID     int64 `form:"app_id"`
	AccountID int64 `form:"account_id"`
	IsBidding bool  `form:"is_bidding"`
}

type admobLineItemCSV struct {
	AdFormat string          `csv:"ad_format"`
	BidFloor decimal.Decimal `csv:"bid_floor"`
	AdUnitID string          `csv:"ad_unit_id"`
}

func (csv admobLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"ad_unit_id": csv.AdUnitID,
		},
	}

	return lineItemAttrs, nil
}

type applovinLineItemCSV struct {
	AdFormat string          `csv:"ad_format"`
	BidFloor decimal.Decimal `csv:"bid_floor"`
	ZoneID   string          `csv:"zone_id"`
}

func (csv applovinLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"zone_id": csv.ZoneID,
		},
	}

	return lineItemAttrs, nil

}

type dtExchangeLineItemCSV struct {
	AdFormat string          `csv:"ad_format"`
	BidFloor decimal.Decimal `csv:"bid_floor"`
	SpotID   string          `csv:"spot_id"`
}

func (csv dtExchangeLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"spot_id": csv.SpotID,
		},
	}

	return lineItemAttrs, nil
}

type gamLineItemCSV struct {
	AdFormat string          `csv:"ad_format"`
	BidFloor decimal.Decimal `csv:"bid_floor"`
	AdUnitID string          `csv:"ad_unit_id"`
}

func (csv gamLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"ad_unit_id": csv.AdUnitID,
		},
	}

	return lineItemAttrs, nil
}

type inmobiLineItemCSV struct {
	AdFormat    string          `csv:"ad_format"`
	BidFloor    decimal.Decimal `csv:"bid_floor"`
	PlacementID string          `csv:"placement_id"`
}

func (csv inmobiLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"placement_id": csv.PlacementID,
		},
	}

	return lineItemAttrs, nil
}

type unityAdsLineItemCSV struct {
	AdFormat    string          `csv:"ad_format"`
	BidFloor    decimal.Decimal `csv:"bid_floor"`
	PlacementID string          `csv:"placement_id"`
}

func (csv unityAdsLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"placement_id": csv.PlacementID,
		},
	}

	return lineItemAttrs, nil
}

type yandexLineItemCSV struct {
	AdFormat string          `csv:"ad_format"`
	BidFloor decimal.Decimal `csv:"bid_floor"`
	AdUnitID string          `csv:"ad_unit_id"`
}

func (csv yandexLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"ad_unit_id": csv.AdUnitID,
		},
	}

	return lineItemAttrs, nil
}

type IronSourceLineItemCSV struct {
	AdFormat   string          `csv:"ad_format"`
	BidFloor   decimal.Decimal `csv:"bid_floor"`
	InstanceID string          `csv:"instance_id"`
}

func (csv IronSourceLineItemCSV) buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error) {
	adType, format := parseCSVAdFormat(csv.AdFormat)
	if adType == ad.UnknownType {
		return LineItemAttrs{}, fmt.Errorf("unknown ad format %q", csv.AdFormat)
	}

	lineItemAttrs := LineItemAttrs{
		HumanName:   strings.ToLower(fmt.Sprintf("%v_%v_%v", account.DemandSource.ApiKey, csv.AdFormat, csv.BidFloor)),
		AppID:       attrs.AppID,
		BidFloor:    &csv.BidFloor,
		AdType:      adType,
		Format:      format,
		AccountID:   account.ID,
		AccountType: account.Type,
		IsBidding:   &attrs.IsBidding,
		Extra: map[string]any{
			"instance_id": csv.InstanceID,
		},
	}

	return lineItemAttrs, nil
}

func parseCSVAdFormat(adFormat string) (ad.Type, *ad.Format) {
	switch strings.ToLower(adFormat) {
	case "banner":
		return ad.BannerType, ptr(ad.BannerFormat)
	case "interstitial":
		return ad.InterstitialType, nil
	case "rewarded":
		return ad.RewardedType, nil
	default:
		return ad.UnknownType, nil
	}
}

func ptr[T any](v T) *T {
	return &v
}

type LineItemCSV interface {
	buildLineItemAttrs(account *DemandSourceAccount, attrs LineItemImportCSVAttrs) (LineItemAttrs, error)
}

func (s *LineItemService) ImportCSV(ctx context.Context, _ AuthContext, reader io.Reader, attrs LineItemImportCSVAttrs) error {
	account, err := s.store.DemandSourceAccounts().Find(ctx, attrs.AccountID)
	if err != nil {
		return fmt.Errorf("find account: %v", err)
	}

	csvInput, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read csv: %v", err)
	}

	var csvLineItems []LineItemCSV
	switch adapter.Key(account.DemandSource.ApiKey) {
	case adapter.AdmobKey:
		var admobLineItems []admobLineItemCSV
		err = csvutil.Unmarshal(csvInput, &admobLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(admobLineItems))
		for i, admobLineItem := range admobLineItems {
			csvLineItems[i] = admobLineItem
		}
	case adapter.ApplovinKey:
		var applovinLineItems []applovinLineItemCSV
		err = csvutil.Unmarshal(csvInput, &applovinLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(applovinLineItems))
		for i, applovinLineItem := range applovinLineItems {
			csvLineItems[i] = applovinLineItem
		}
	case adapter.DTExchangeKey:
		var dtExchangeLineItems []dtExchangeLineItemCSV
		err = csvutil.Unmarshal(csvInput, &dtExchangeLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(dtExchangeLineItems))
		for i, dtExchangeLineItem := range dtExchangeLineItems {
			csvLineItems[i] = dtExchangeLineItem
		}
	case adapter.GAMKey:
		var gamLineItems []gamLineItemCSV
		err = csvutil.Unmarshal(csvInput, &gamLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(gamLineItems))
		for i, gamLineItem := range gamLineItems {
			csvLineItems[i] = gamLineItem
		}
	case adapter.InmobiKey:
		var inmobiLineItems []inmobiLineItemCSV
		err = csvutil.Unmarshal(csvInput, &inmobiLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(inmobiLineItems))
		for i, inmobiLineItem := range inmobiLineItems {
			csvLineItems[i] = inmobiLineItem
		}
	case adapter.UnityAdsKey:
		var unityLineItems []unityAdsLineItemCSV
		err = csvutil.Unmarshal(csvInput, &unityLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(unityLineItems))
		for i, unityAdsLineItem := range unityLineItems {
			csvLineItems[i] = unityAdsLineItem
		}
	case adapter.YandexKey:
		var yandexLineItems []yandexLineItemCSV
		err = csvutil.Unmarshal(csvInput, &yandexLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(yandexLineItems))
		for i, yandexLineItem := range yandexLineItems {
			csvLineItems[i] = yandexLineItem
		}
	case adapter.IronSourceKey:
		var ironSourceLineItems []IronSourceLineItemCSV
		err = csvutil.Unmarshal(csvInput, &ironSourceLineItems)
		if err != nil {
			return fmt.Errorf("unmarshal csv: %v", err)
		}

		csvLineItems = make([]LineItemCSV, len(ironSourceLineItems))
		for i, ironSourceLineItem := range ironSourceLineItems {
			csvLineItems[i] = ironSourceLineItem
		}
	default:
		return fmt.Errorf("unsupported demand source: %s", account.DemandSource.ApiKey)
	}
	if len(csvLineItems) == 0 {
		return errors.New("csv empty")
	}

	lineItemsAttrs := make([]LineItemAttrs, len(csvLineItems))
	for i, csvLineItem := range csvLineItems {
		lineItemsAttrs[i], err = csvLineItem.buildLineItemAttrs(account, attrs)
		if err != nil {
			return fmt.Errorf("build line item attrs: %v", err)
		}
	}

	err = s.store.LineItems().CreateMany(ctx, lineItemsAttrs)
	if err != nil {
		return fmt.Errorf("create line items: %v", err)
	}

	return nil
}

type LineItemRepo interface {
	AllResourceQuerier[LineItem]
	OwnedResourceQuerier[LineItem]
	ResourceManipulator[LineItem, LineItemAttrs]

	CreateMany(ctx context.Context, items []LineItemAttrs) error
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

func (p *lineItemPolicy) permissions(_ AuthContext) ResourcePermissions {
	return ResourcePermissions{
		Read:   true,
		Create: true,
	}
}

func (p *lineItemPolicy) instancePermissions(_ AuthContext, _ *LineItem) ResourceInstancePermissions {
	return ResourceInstancePermissions{
		Update: true,
		Delete: true,
	}
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
	case adapter.AmazonKey:
		rule = v8n.Map(
			v8n.Key("slot_uuid", v8n.Required, isString),
			v8n.Key("format", v8n.Required, isString, v8n.In("BANNER", "MREC", "INTERSTITIAL", "VIDEO", "REWARDED")),
		)
	case adapter.ApplovinKey:
		rule = v8n.Map(
			v8n.Key("zone_id", v8n.Required, isString),
		)
	case adapter.BigoAdsKey:
		rule = v8n.Map(
			v8n.Key("slot_id", v8n.Required, isString),
		)
	case adapter.MetaKey, adapter.UnityAdsKey, adapter.VungleKey, adapter.MobileFuseKey:
		rule = v8n.Map(
			v8n.Key("placement_id", v8n.Required, isString),
		)
	case adapter.DTExchangeKey:
		rule = v8n.Map(
			v8n.Key("spot_id", v8n.Required, isString),
		)
	case adapter.IronSourceKey:
		rule = v8n.Map(
			v8n.Key("instance_id", v8n.Required, isString),
		)
	case adapter.MintegralKey:
		rule = v8n.Map(
			v8n.Key("placement_id", v8n.Required, isString),
			v8n.Key("unit_id", v8n.Required, isString),
		)
	case adapter.VKAdsKey:
		rule = v8n.Map(
			v8n.Key("slot_id", v8n.Required, isString),
			v8n.Key("mediation", isString),
		)
	case adapter.YandexKey:
		rule = v8n.Map(
			v8n.Key("ad_unit_id", v8n.Required, isString),
		)
	}

	return rule.AllowExtraKeys()
}
