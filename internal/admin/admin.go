// Package admin implements services for managing resources.
package admin

//go:generate go run -mod=mod github.com/matryer/moq@latest -out admin_mocks_test.go . Store

// Service is a top-level service for managing resources.
type Service struct {
	AppService                    *AppService
	AppDemandProfileService       *AppDemandProfileService
	AuctionConfigurationService   *AuctionConfigurationService
	AuctionConfigurationV2Service *AuctionConfigurationV2Service
	CountryService                *CountryService
	DemandSourceService           *DemandSourceService
	DemandSourceAccountService    *DemandSourceAccountService
	LineItemService               *LineItemService
	SegmentService                *SegmentService
	UserService                   *UserService
}

// NewService creates a new Service.
func NewService(store Store) *Service {
	return &Service{
		AppService:                    NewAppService(store),
		AppDemandProfileService:       NewAppDemandProfileService(store),
		AuctionConfigurationService:   NewAuctionConfigurationService(store),
		AuctionConfigurationV2Service: NewAuctionConfigurationV2Service(store),
		CountryService:                NewCountryService(store),
		DemandSourceService:           NewDemandSourceService(store),
		DemandSourceAccountService:    NewDemandSourceAccountService(store),
		LineItemService:               NewLineItemService(store),
		SegmentService:                NewSegmentService(store),
		UserService:                   NewUserService(store),
	}
}

// Store is an interface for accessing resources from storage.
type Store interface {
	Apps() AppRepo
	AppDemandProfiles() AppDemandProfileRepo
	AuctionConfigurations() AuctionConfigurationRepo
	AuctionConfigurationsV2() AuctionConfigurationV2Repo
	Countries() CountryRepo
	DemandSources() DemandSourceRepo
	DemandSourceAccounts() DemandSourceAccountRepo
	LineItems() LineItemRepo
	Segments() SegmentRepo
	Users() UserRepo
}
