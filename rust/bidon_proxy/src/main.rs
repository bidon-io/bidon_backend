use axum::{
    routing::post,
    Router,
};
use std::sync::Arc;
use tokio::sync::Mutex;
use crate::bidon_version::XBidonVersionString;
use swagger::{ContextBuilder, XSpanIdString, AuthData};
use crate::context::BidonContext;
use crate::auction;
use crate::controllers;

// mod models;
// mod auction;
// mod controllers;

#[tokio::main]
async fn main() {

    // Create a ProxyServer instance
    let auction = Arc::new(Mutex::new(auction::SimpleAuction::new()));

    // Define the context
    let context = BidonContext::new()
        .push(XSpanIdString::default())
        .push(XBidonVersionString::default())
        .push(None::<AuthData>)
        .build();

    // Define the routes
    let app = Router::new()
        .route("/v2/auction/:ad_type", post(controllers::auction::get_auction_handler))
        .layer(axum::extract::Extension(auction))
        .layer(axum::extract::Extension(context));

    // Start the server
    axum::Server::bind(&"127.0.0.1:3030".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
