package store

import "github.com/bidon-io/bidon-backend/internal/admin"

type DemandSourceRepo = resourceRepo[admin.DemandSource, admin.DemandSourceAttrs, demandSource, *demandSource]

type demandSource struct {
	Model
	ApiKey    string `gorm:"column:api_key;type:varchar;not null"`
	HumanName string `gorm:"column:human_name;type:varchar;not null"`
}

//lint:ignore U1000 this method is used by generic struct
func (s *demandSource) initFromResourceAttrs(attrs *admin.DemandSourceAttrs) {
	s.ApiKey = attrs.ApiKey
	s.HumanName = attrs.HumanName
}

//lint:ignore U1000 this method is used by generic struct
func (s *demandSource) toResource() admin.DemandSource {
	return admin.DemandSource{
		ID: s.ID,
		DemandSourceAttrs: admin.DemandSourceAttrs{
			ApiKey:    s.ApiKey,
			HumanName: s.HumanName,
		},
	}
}
