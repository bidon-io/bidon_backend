package admin

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/adapter"
)

const AuctionConfigurationV2ResourceKey = "auction_configuration_v2"

type AuctionConfigurationV2Resource struct {
	*AuctionConfigurationV2
	Permissions ResourceInstancePermissions `json:"_permissions"`
}

type AuctionConfigurationV2 struct {
	ID         int64  `json:"id"`
	PublicUID  string `json:"public_uid"`
	AuctionKey string `json:"auction_key"`
	AuctionConfigurationV2Attrs
	App     App      `json:"app"`
	Segment *Segment `json:"segment"`
}

// AuctionConfigurationV2Attrs is attributes of Configuration. Used to create and update configurations
type AuctionConfigurationV2Attrs struct {
	Name                     string         `json:"name"`
	AppID                    int64          `json:"app_id"`
	AdType                   ad.Type        `json:"ad_type"`
	Pricefloor               float64        `json:"pricefloor"`
	SegmentID                *int64         `json:"segment_id"`
	ExternalWinNotifications *bool          `json:"external_win_notifications"`
	Demands                  []adapter.Key  `json:"demands"`
	Bidding                  []adapter.Key  `json:"bidding"`
	AdUnitIDs                []int64        `json:"ad_unit_ids"`
	Timeout                  int32          `json:"timeout"`
	Settings                 map[string]any `json:"settings"`
}

type AuctionConfigurationV2Service struct {
	*ResourceService[AuctionConfigurationV2Resource, AuctionConfigurationV2, AuctionConfigurationV2Attrs]
}

func NewAuctionConfigurationV2Service(store Store) *AuctionConfigurationV2Service {
	s := &AuctionConfigurationV2Service{
		ResourceService: &ResourceService[AuctionConfigurationV2Resource, AuctionConfigurationV2, AuctionConfigurationV2Attrs]{},
	}

	s.resourceKey = AuctionConfigurationV2ResourceKey

	s.repo = store.AuctionConfigurationsV2()
	s.policy = newAuctionConfigurationV2Policy(store)

	s.prepareResource = func(authCtx AuthContext, config *AuctionConfigurationV2) AuctionConfigurationV2Resource {
		return AuctionConfigurationV2Resource{
			AuctionConfigurationV2: config,
			Permissions:            s.policy.instancePermissions(authCtx, config),
		}
	}

	return s
}

type AuctionConfigurationV2Repo interface {
	AllResourceQuerier[AuctionConfigurationV2]
	OwnedResourceQuerier[AuctionConfigurationV2]
	ResourceManipulator[AuctionConfigurationV2, AuctionConfigurationV2Attrs]
}

type auctionConfigurationV2Policy struct {
	repo AuctionConfigurationV2Repo

	appPolicy     *appPolicy
	segmentPolicy *segmentPolicy
}

func newAuctionConfigurationV2Policy(store Store) *auctionConfigurationV2Policy {
	return &auctionConfigurationV2Policy{
		repo: store.AuctionConfigurationsV2(),

		appPolicy:     newAppPolicy(store),
		segmentPolicy: newSegmentPolicy(store),
	}
}

func (p *auctionConfigurationV2Policy) getReadScope(authCtx AuthContext) resourceScope[AuctionConfigurationV2] {
	return &ownedResourceScope[AuctionConfigurationV2]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *auctionConfigurationV2Policy) getManageScope(authCtx AuthContext) resourceScope[AuctionConfigurationV2] {
	return &ownedResourceScope[AuctionConfigurationV2]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *auctionConfigurationV2Policy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *AuctionConfigurationV2Attrs) error {
	// Check if user can manage the app.
	_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
	if err != nil {
		return err
	}

	if attrs.SegmentID != nil {
		// Check if user can read the segment.
		_, err = p.segmentPolicy.getReadScope(authCtx).find(ctx, *attrs.SegmentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *auctionConfigurationV2Policy) authorizeUpdate(ctx context.Context, authCtx AuthContext, config *AuctionConfigurationV2, attrs *AuctionConfigurationV2Attrs) error {
	// If user tries to change the app and app is not the same as before, check if user can manage the new app.
	if attrs.AppID != 0 && attrs.AppID != config.AppID {
		_, err := p.appPolicy.getManageScope(authCtx).find(ctx, attrs.AppID)
		if err != nil {
			return err
		}
	}

	// If user tries to change the segment and segment is not the same as before, check if user can read the new segment.
	if attrs.SegmentID != nil && *attrs.SegmentID != *config.SegmentID {
		_, err := p.segmentPolicy.getReadScope(authCtx).find(ctx, *attrs.SegmentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *auctionConfigurationV2Policy) authorizeDelete(_ context.Context, _ AuthContext, _ *AuctionConfigurationV2) error {
	return nil
}

func (p *auctionConfigurationV2Policy) permissions(_ AuthContext) ResourcePermissions {
	return ResourcePermissions{
		Read:   true,
		Create: true,
	}
}

func (p *auctionConfigurationV2Policy) instancePermissions(_ AuthContext, _ *AuctionConfigurationV2) ResourceInstancePermissions {
	return ResourceInstancePermissions{
		Update: true,
		Delete: true,
	}
}
