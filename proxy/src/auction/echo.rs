use crate::auction::Api;
use crate::auction::AuctionError;
use crate::com::iabtechlab::openrtb::v3::Openrtb;

// This is a simple echo auction that returns the same request as the response.
// It's useful for testing the API.
#[derive(Debug, Clone)]
pub struct EchoAuction;

#[allow(dead_code)]
impl EchoAuction {
    pub fn new() -> Self {
        EchoAuction
    }
}

#[async_trait::async_trait]
impl Api for EchoAuction {
    async fn bid(&mut self, auction_request: Openrtb) -> Result<Openrtb, AuctionError> {
        Ok(auction_request)
    }
}
