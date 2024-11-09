use axum::{
    routing::post,
    Router,
};
use std::sync::Arc;
use tokio::sync::Mutex;
use galaxy::bidon_version::XBidonVersionString;
// use crate::auction;
use galaxy::controllers;
use swagger::{ContextBuilder, XSpanIdString, AuthData};
use galaxy::context::{MyContext, MyEmpContext};
use swagger::Push;
use galaxy::auction::SimpleAuction;

#[tokio::main]
async fn main() {

    // Create a ProxyServer instance
    let auction = Arc::new(Mutex::new(galaxy::auction::SimpleAuction::new("http://localhost:50051".to_string()).await.unwrap()));

    // Define the context
    // let context = swagger::make_context!(MyContext, MyEmpContext,
    //     XSpanIdString::default(),
    //     XBidonVersionString::default(),
    //     None::<AuthData>);

    // Define the routes
    let app = Router::new()
        .route("/v2/auction/:ad_type", post(controllers::auction::get_auction_handler::<SimpleAuction>))
        .layer(axum::extract::Extension(auction));
        // .layer(axum::extract::Extension(context) );

    // Start the server
    axum::Server::bind(&"127.0.0.1:3030".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
