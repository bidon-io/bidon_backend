package admin

import (
	"context"

	"github.com/bidon-io/bidon-backend/internal/ad"
	"github.com/bidon-io/bidon-backend/internal/auction"
)

type AuctionConfiguration struct {
	ID int64 `json:"id"`
	AuctionConfigurationAttrs
	App     App      `json:"app"`
	Segment *Segment `json:"segment"`
}

// AuctionConfigurationAttrs is attributes of Configuration. Used to create and update configurations
type AuctionConfigurationAttrs struct {
	Name                     string                `json:"name"`
	AppID                    int64                 `json:"app_id"`
	AdType                   ad.Type               `json:"ad_type"`
	Rounds                   []auction.RoundConfig `json:"rounds"`
	Pricefloor               float64               `json:"pricefloor"`
	SegmentID                *int64                `json:"segment_id"`
	ExternalWinNotifications *bool                 `json:"external_win_notifications"`
}

type AuctionConfigurationService = ResourceService[AuctionConfiguration, AuctionConfigurationAttrs]

func NewAuctionConfigurationService(store Store) *AuctionConfigurationService {
	return &AuctionConfigurationService{
		repo:   store.AuctionConfigurations(),
		policy: newAuctionConfigurationPolicy(store),
	}
}

type AuctionConfigurationRepo interface {
	AllResourceQuerier[AuctionConfiguration]
	OwnedResourceQuerier[AuctionConfiguration]
	ResourceManipulator[AuctionConfiguration, AuctionConfigurationAttrs]
}

type auctionConfigurationPolicy struct {
	repo AuctionConfigurationRepo

	appPolicy     *appPolicy
	segmentPolicy *segmentPolicy
}

func newAuctionConfigurationPolicy(store Store) *auctionConfigurationPolicy {
	return &auctionConfigurationPolicy{
		repo: store.AuctionConfigurations(),

		appPolicy:     newAppPolicy(store),
		segmentPolicy: newSegmentPolicy(store),
	}
}

func (p *auctionConfigurationPolicy) getReadScope(authCtx AuthContext) resourceScope[AuctionConfiguration] {
	return &ownedResourceScope[AuctionConfiguration]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *auctionConfigurationPolicy) getManageScope(authCtx AuthContext) resourceScope[AuctionConfiguration] {
	return &ownedResourceScope[AuctionConfiguration]{
		repo:    p.repo,
		authCtx: authCtx,
	}
}

func (p *auctionConfigurationPolicy) authorizeCreate(ctx context.Context, authCtx AuthContext, attrs *AuctionConfigurationAttrs) error {
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

func (p *auctionConfigurationPolicy) authorizeUpdate(ctx context.Context, authCtx AuthContext, config *AuctionConfiguration, attrs *AuctionConfigurationAttrs) error {
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

func (p *auctionConfigurationPolicy) authorizeDelete(_ context.Context, _ AuthContext, _ *AuctionConfiguration) error {
	return nil
}
