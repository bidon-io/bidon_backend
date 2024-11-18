use axum::routing::post;
use axum::{middleware, Router};
use galaxy::auction::SimpleAuction;
use galaxy::bidon_version::XBidonVersionString;
use galaxy::controllers;
use std::sync::Arc;
use tokio::sync::Mutex;

#[tokio::main]
async fn main() {
    // Create a ProxyServer instance
    let auction = Arc::new(Mutex::new(
        galaxy::auction::SimpleAuction::new("http://localhost:50051".to_string())
            .await
            .unwrap(),
    ));

    // Define the routes
    let app = Router::new()
        .route(
            "/v2/auction/:ad_type",
            post(controllers::auction::get_auction_handler::<SimpleAuction>),
        )
        .route_layer(middleware::from_fn(
            XBidonVersionString::extract_header_middleware,
        ))
        .layer(axum::extract::Extension(auction));

    // Start the server
    axum::Server::bind(&"127.0.0.1:3030".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
