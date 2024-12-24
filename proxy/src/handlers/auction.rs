use crate::adapter;
use crate::bidding::BiddingService;
use crate::extract::AuctionRequestPayload;
use axum::extract::State;
use axum::{
    http::StatusCode,
    response::{IntoResponse, Json},
};

pub async fn get_auction_handler<S>(
    State(bidding_service): State<Box<S>>,
    AuctionRequestPayload(request): AuctionRequestPayload,
) -> impl IntoResponse
where
    S: BiddingService + Send + Sync,
{
    match bidding_service.bid(request).await {
        Ok(response) => {
            tracing::trace!("Bidding response: {:?}", response);
            match adapter::try_into(response) {
                Ok(auction_response) => Json::from(auction_response).into_response(),
                Err(err) => (StatusCode::INTERNAL_SERVER_ERROR, err.to_string()).into_response(),
            }
        }
        Err(err) => {
            tracing::error!("Bidding error: {:?}", err);
            (StatusCode::INTERNAL_SERVER_ERROR, "Internal Server Error").into_response()
        }
    }
}
