mod echo;
mod proxy;

pub use echo::EchoBiddingService;
pub use proxy::ProxyBiddingService;
use std::{error, fmt};

use galaxy_bidon::com::iabtechlab::openrtb::v3::Openrtb;

#[async_trait::async_trait]
pub trait Api {
    /// Bidding service
    async fn bid(&mut self, bidding_request: Openrtb) -> Result<Openrtb, BiddingError>;
}

// todo check errors in openrtb
#[derive(Clone, Debug)]
pub struct BiddingError(pub String);

impl BiddingError {
    pub fn new(msg: String) -> Self {
        BiddingError(msg)
    }
}

impl fmt::Display for BiddingError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let debug: &dyn fmt::Debug = self;
        debug.fmt(f)
    }
}

impl error::Error for BiddingError {
    fn description(&self) -> &str {
        "Failed to produce a valid response."
    }
}
