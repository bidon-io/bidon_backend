package schema

type AuctionRequest struct {
	BaseRequest
	Adapters Adapters `json:"adapters" validate:"required"`
	AdObject AdObject `json:"ad_object" validate:"required"`
}
