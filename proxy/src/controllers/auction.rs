use crate::bidding::Api as BiddingAPI;
use axum::extract::State;
use axum::{
    extract::Path,
    http::StatusCode,
    response::IntoResponse,
};
use galaxy_bidon::extractor::BidonOpenRTBExtractor;
use galaxy_bidon::models::GetAuctionAdTypeParameter;
use prost::bytes::BytesMut;
use prost::Message;

// #[axum::debug_handler]
pub async fn get_auction_handler<A>(
    Path(ad_type): Path<String>,
    State(mut auction): State<Box<A>>,
    BidonOpenRTBExtractor(openrtb_request): BidonOpenRTBExtractor,
) -> impl IntoResponse
where
    A: BiddingAPI + Send + Sync,
{
    // TODO use ad_type to determine the bidding type.
    let _ad_type = match ad_type.parse::<GetAuctionAdTypeParameter>() {
        Ok(ad_type) => ad_type,
        Err(_) => return (StatusCode::BAD_REQUEST, "Invalid ad_type").into_response(),
    };

    // TODO use multiple bidding to avoid lock contention.
    match auction.bid(openrtb_request).await {
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
