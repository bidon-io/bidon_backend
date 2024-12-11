pub(crate) use crate::bidding::BiddingService;
pub(crate) use crate::bidding::BiddingError;
use crate::com::iabtechlab::openrtb::v3::Openrtb;

// This is a simple echo bidding that returns the same request as the response.
// It's useful for testing the API.
#[derive(Debug, Clone, Default)]
pub struct EchoBiddingService;

#[allow(dead_code)]
impl EchoBiddingService {
    pub fn new() -> Self {
        EchoBiddingService
    }
}

#[async_trait::async_trait]
impl BiddingService for EchoBiddingService {
    async fn bid(&self, request: Openrtb) -> Result<Openrtb, BiddingError> {
        Ok(request)
    }
}
