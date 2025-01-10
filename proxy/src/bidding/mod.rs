mod echo;
pub mod error;
mod proxy;

pub use echo::EchoBiddingService;
pub use proxy::ProxyBiddingService;

use crate::bidding::error::BiddingError;
use crate::com::iabtechlab::openrtb::v3::Openrtb;

#[async_trait::async_trait]
pub trait BiddingService {
    /// Bidding service
    async fn bid(&self, request: Openrtb) -> Result<Openrtb, BiddingError>;
}
