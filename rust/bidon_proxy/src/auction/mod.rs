mod simple;

use std::{error, fmt};
pub use simple::SimpleAuction;

use crate::com::iabtechlab::openrtb::v3::Openrtb;

#[async_trait::async_trait]
pub trait Api {
    /// Auction
    async fn bid(&mut self, auction_request: Openrtb ) -> Result<Openrtb, AuctionError>;
}

#[derive(Clone, Debug)]
pub struct AuctionError(pub String);

impl AuctionError {
    pub fn new(msg: String) -> Self {
        AuctionError(msg)
    }
}

impl fmt::Display for AuctionError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let debug: &dyn fmt::Debug = self;
        debug.fmt(f)
    }
}

impl error::Error for AuctionError {
    fn description(&self) -> &str {
        "Failed to produce a valid response."
    }
}
