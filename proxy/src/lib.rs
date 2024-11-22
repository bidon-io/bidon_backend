use crate::bidding::Api as BiddingApi;
use axum::routing::post;
use axum::Router;

// Define the routes
pub fn create_app<A>(auction: Box<A>) -> Router
where
    A: Clone + 'static,
    A: BiddingApi + Send + Sync,
{
    Router::new()
        .route(
            "/v2/auction/:ad_type",
            post(controllers::auction::get_auction_handler),
        )
        .with_state(auction)
}

pub mod bidding;

pub mod controllers {
    pub mod auction;
}

