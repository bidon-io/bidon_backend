use crate::bidding::Api;
use crate::bidding::BiddingError;
use galaxy_bidon::com::iabtechlab::openrtb::v3::Openrtb;

// This is a simple echo bidding that returns the same request as the response.
// It's useful for testing the API.
#[derive(Debug, Clone)]
pub struct EchoBiddingService;

#[allow(dead_code)]
impl EchoBiddingService {
    pub fn new() -> Self {
        EchoBiddingService
    }
}

#[async_trait::async_trait]
impl Api for EchoBiddingService {
    async fn bid(&mut self, bidding_request: Openrtb) -> Result<Openrtb, BiddingError> {
        Ok(bidding_request)
    }
}
