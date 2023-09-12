package adminstore

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type Store struct {
	AppRepo                  *AppRepo
	AppDemandProfileRepo     *AppDemandProfileRepo
	AuctionConfigurationRepo *AuctionConfigurationRepo
	CountryRepo              *CountryRepo
	DemandSourceRepo         *DemandSourceRepo
	DemandSourceAccountRepo  *DemandSourceAccountRepo
	LineItemRepo             *LineItemRepo
	SegmentRepo              *SegmentRepo
	UserRepo                 *UserRepo
}

func New(db *db.DB) *Store {
	return &Store{
		AppRepo:                  NewAppRepo(db),
		AppDemandProfileRepo:     NewAppDemandProfileRepo(db),
		AuctionConfigurationRepo: NewAuctionConfigurationRepo(db),
		CountryRepo:              NewCountryRepo(db),
		DemandSourceRepo:         NewDemandSourceRepo(db),
		DemandSourceAccountRepo:  NewDemandSourceAccountRepo(db),
		LineItemRepo:             NewLineItemRepo(db),
		SegmentRepo:              NewSegmentRepo(db),
		UserRepo:                 NewUserRepo(db),
	}
}

func (s *Store) Apps() admin.AppRepo {
	return s.AppRepo
}

func (s *Store) AppDemandProfiles() admin.AppDemandProfileRepo {
	return s.AppDemandProfileRepo
}

func (s *Store) AuctionConfigurations() admin.AuctionConfigurationRepo {
	return s.AuctionConfigurationRepo
}

func (s *Store) Countries() admin.CountryRepo {
	return s.CountryRepo
}

func (s *Store) DemandSources() admin.DemandSourceRepo {
	return s.DemandSourceRepo
}

func (s *Store) DemandSourceAccounts() admin.DemandSourceAccountRepo {
	return s.DemandSourceAccountRepo
}

func (s *Store) LineItems() admin.LineItemRepo {
	return s.LineItemRepo
}

func (s *Store) Segments() admin.SegmentRepo {
	return s.SegmentRepo
}

func (s *Store) Users() admin.UserRepo {
	return s.UserRepo
}

func platformID(platformID db.PlatformID) admin.PlatformID {
	switch platformID {
	case db.AndroidPlatformID:
		return admin.AndroidPlatformID
	case db.IOSPlatformID:
		return admin.IOSPlatformID
	default:
		return admin.UnknownPlatformID
	}
}

func dbPlatformID(platformID admin.PlatformID) db.PlatformID {
	switch platformID {
	case admin.AndroidPlatformID:
		return db.AndroidPlatformID
	case admin.IOSPlatformID:
		return db.IOSPlatformID
	default:
		return db.UnknownPlatformID
	}
}
