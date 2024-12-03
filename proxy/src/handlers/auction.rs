use crate::bidding::BiddingService;
use crate::extract::AuctionRequestPayload;
use crate::sdk::GetAuctionAdTypeParameter;
use axum::extract::State;
use axum::{extract::Path, http::StatusCode, response::IntoResponse};
use prost::bytes::BytesMut;
use prost::Message;

pub async fn get_auction_handler<S>(
    Path(ad_type): Path<String>,
    State(bidding_service): State<Box<S>>,
    AuctionRequestPayload(request): AuctionRequestPayload,
) -> impl IntoResponse
where
    S: BiddingService + Send + Sync,
{
    // TODO use ad_type to determine the bidding type.
    let _ad_type = match ad_type.parse::<GetAuctionAdTypeParameter>() {
        Ok(ad_type) => ad_type,
        Err(_) => return (StatusCode::BAD_REQUEST, "Invalid ad_type").into_response(),
    };

    match bidding_service.bid(request).await {
        Ok(response) => {
            let mut buf = BytesMut::with_capacity(128);
            match response.encode(&mut buf) {
                Ok(()) => buf.into_response(),
                Err(err) => (StatusCode::INTERNAL_SERVER_ERROR, err.to_string()).into_response(),
            }
        }
        Err(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Internal Server Error").into_response(),
    }
}
