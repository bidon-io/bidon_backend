use galaxy::auction::EchoAuction;
// use std::sync::Arc;
// use tokio::sync::Mutex;

#[tokio::main]
async fn main() {
    // Create a ProxyServer instance
    // let auction = Arc::new(Mutex::new(
    //     galaxy::auction::SimpleAuction::new("http://localhost:50051".to_string())
    //         .await
    //         .unwrap(),
    // ));
    let auction = Box::new(EchoAuction::new());

    let app = galaxy::create_app(auction);
    // Start the server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
