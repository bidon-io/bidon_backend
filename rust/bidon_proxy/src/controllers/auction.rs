use crate::auction::Api as AuctionApi;
use crate::controllers::adapter;
use crate::models::{AuctionRequest, GetAuctionAdTypeParameter};
use axum::{
    extract::{Extension, Json, Path},
    http::StatusCode,
    response::IntoResponse
    ,
};
use prost::bytes::BytesMut;
use prost::Message;
use std::convert::TryFrom;
use std::sync::Arc;
use tokio::sync::Mutex;
use crate::bidon_version::XBidonVersionString;

pub async fn get_auction_handler<A>(
    Path(ad_type): Path<String>,
    Json(auction_request): Json<AuctionRequest>,
    Extension(auction): Extension<Arc<Mutex<A>>>,
    Extension(xbidon_version_string): Extension<XBidonVersionString>,
) -> impl IntoResponse
where
    A: AuctionApi,
{
    let ad_type = match ad_type.parse::<GetAuctionAdTypeParameter>() {
        Ok(ad_type) => ad_type,
        Err(_) => return (StatusCode::BAD_REQUEST, "Invalid ad_type").into_response(),
    };

    let openrtb_request = match adapter::try_from(auction_request) {
        Ok(req) => req,
        Err(_) => return (StatusCode::BAD_REQUEST, "Invalid auction request").into_response(),
    };
    // TODO use xbidon_version_string and ad_type to determine the auction type.
   
    // TODO use multiple auction to avoid lock contention.
    let mut auction = auction.lock().await;
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
