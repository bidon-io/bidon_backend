use crate::{models, ServiceError};
use swagger::{ApiError, ContextWrapper};
use async_trait::async_trait;
use std::error::Error;
use std::task::{Context, Poll};
use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Serialize, Deserialize)]
pub enum GetAuctionResponse {
    /// Auction response
    AuctionResponse(models::AuctionResponse)
}

/// API
#[async_trait]
#[allow(clippy::too_many_arguments, clippy::ptr_arg)]
pub trait Api<C: Send + Sync> {
    fn poll_ready(&self, _cx: &mut Context) -> Poll<Result<(), Box<dyn Error + Send + Sync + 'static>>> {
        Poll::Ready(Ok(()))
    }

    /// Auction
    async fn get_auction(&self, ad_type: models::GetAuctionAdTypeParameter,
                         auction_request: models::AuctionRequest,
                         context: &C) -> Result<GetAuctionResponse, ApiError>;

}
