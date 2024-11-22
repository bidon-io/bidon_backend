#![allow(
    missing_docs,
    trivial_casts,
    unused_variables,
    unused_mut,
    unused_imports,
    unused_extern_crates,
    unused_attributes,
    non_camel_case_types
)]
#![allow(clippy::derive_partial_eq_without_eq, clippy::disallowed_names)]

use crate::bidding::Api as BiddingApi;
use axum::routing::post;
use axum::{middleware, Router, async_trait};
use axum::extract::{Json, Request, FromRequest, Extension, State};
use axum::http::StatusCode;
use std::error::Error;
use axum::extract::rejection::JsonRejection;

pub const BASE_PATH: &str = "";
pub const API_VERSION: &str = "1.0.0";

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

