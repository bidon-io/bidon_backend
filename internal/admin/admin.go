// Package admin implements services for managing resources.
package admin

// Service is a top-level service for managing resources.
type Service struct {
	AppService                  *AppService
	AppDemandProfileService     *AppDemandProfileService
	AuctionConfigurationService *AuctionConfigurationService
	CountryService              *CountryService
	DemandSourceService         *DemandSourceService
	DemandSourceAccountService  *DemandSourceAccountService
	LineItemService             *LineItemService
	SegmentService              *SegmentService
	UserService                 *UserService
}

// NewService creates a new Service.
func NewService(store Store) *Service {
	return &Service{
		AppService:                  NewAppService(store),
		AppDemandProfileService:     NewAppDemandProfileService(store),
		AuctionConfigurationService: NewAuctionConfigurationService(store),
		CountryService:              NewCountryService(store),
		DemandSourceService:         NewDemandSourceService(store),
		DemandSourceAccountService:  NewDemandSourceAccountService(store),
		LineItemService:             NewLineItemService(store),
		SegmentService:              NewSegmentService(store),
		UserService:                 NewUserService(store),
	}
}

// Store is an interface for accessing resources from storage.
type Store interface {
	Apps() AppRepo
	AppDemandProfiles() AppDemandProfileRepo
	AuctionConfigurations() AuctionConfigurationRepo
	Countries() CountryRepo
	DemandSources() DemandSourceRepo
	DemandSourceAccounts() DemandSourceAccountRepo
	LineItems() LineItemRepo
	Segments() SegmentRepo
	Users() UserRepo
}
