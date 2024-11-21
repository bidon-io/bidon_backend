use galaxy::bidding::EchoBiddingService;
// use std::sync::Arc;
// use tokio::sync::Mutex;

#[tokio::main]
async fn main() {
    // Create a ProxyServer instance
    // let bidding = Arc::new(Mutex::new(
    //     galaxy::bidding::SimpleAuction::new("http://localhost:50051".to_string())
    //         .await
    //         .unwrap(),
    // ));
    let auction = Box::new(EchoBiddingService::new());

    let app = galaxy::create_app(auction);
    // Start the server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
