pub(crate) use bidon::bidding::EchoBiddingService;

#[tokio::main]
async fn main() {
    let echo_bidding_service = Box::new(EchoBiddingService::new());

    let app = bidon::create_app(echo_bidding_service);
    // Start the server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
