use axum::http;
use axum_test_helper::TestClient;
use bidon::create_app;
use http::status::StatusCode;
mod common;

use common::auction_request::get_auction_request;

#[tokio::test]
async fn test_auction() {
    // Create the app
    let app = create_app(Box::new(bidon::bidding::EchoBiddingService::new()));

    // Create a test client
    let client = TestClient::new(app);

    // Send a request with the header
    let response = client
        .await
        .post("/v2/auction/banner")
        .header("x-bidon-version", "TestValue")
        .json(&get_auction_request())
        .send()
        .await;

    // Assert the response
    assert_eq!(response.status(), StatusCode::OK);
    assert!(
        !response.text().await.is_empty(),
        "Output must not be empty"
    );
}
