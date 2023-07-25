package schema

import "github.com/bidon-io/bidon-backend/internal/ad"

type RewardRequest struct {
	ShowRequest
	AdType ad.Type `param:"ad_type" validate:"eq=rewarded"`
}
