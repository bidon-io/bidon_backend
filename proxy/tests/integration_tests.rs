use std::collections::HashMap;

use axum::http;
use axum_test_helper::TestClient;
use bidon::{
    create_app,
    org::bidon::proto::v1::bidding_service_server::{BiddingService, BiddingServiceServer},
};
use http::status::StatusCode;
mod common;

use bidon::com::iabtechlab::openrtb::v3 as openrtb;
use bidon::com::iabtechlab::openrtb::v3::openrtb::PayloadOneof;
use bidon::com::iabtechlab::openrtb::v3::Openrtb;
use bidon::com::iabtechlab::openrtb::v3::Response as OpenrtbResponse;
use bidon::org::bidon::proto::v1::mediation;
use bidon::org::bidon::proto::v1::mediation::{
    APP_EXT, AUCTION_RESPONSE_EXT, BID_EXT, DEVICE_EXT, USER_EXT,
};
use common::auction_request::get_auction_request;
use prost::ExtensionRegistry;
use tonic::{transport::Server, Request, Response, Status};
// Mock bidding service implementation
#[derive(Default)]
struct MockBiddingService;

#[tonic::async_trait]
impl BiddingService for MockBiddingService {
    async fn bid(&self, _request: Request<Openrtb>) -> Result<Response<Openrtb>, Status> {
        // Create a hardcoded OpenRTB response
        let response = Openrtb {
            ver: Some("1.0".to_string()),
            domainspec: Some("".to_string()),
            domainver: Some("".to_string()),
            payload_oneof: Some(PayloadOneof::Response(OpenrtbResponse {
                id: Some("auction123".to_string()),
                bidid: None,
                nbr: None,
                seatbid: vec![openrtb::SeatBid {
                    bid: vec![openrtb::Bid {
                        id: Some("bid1".to_string()),
                        item: Some("item1".to_string()),
                        price: Some(2.5),
                        cid: Some("demand1".to_string()),
                        extension_set: {
                            let mut ext = prost::ExtensionSet::default();
                            let bid_ext = mediation::BidExt {
                                label: Some("key123".to_string()),
                                bid_type: Some("bid_type".to_string()),
                                ext: HashMap::new(),
                            };
                            ext.set_extension_data(mediation::BID_EXT, bid_ext).unwrap();
                            ext
                        },
                        ..Default::default()
                    }],
                    seat: None,
                    ..Default::default()
                }],
                extension_set: {
                    let mut ext = prost::ExtensionSet::default();
                    let auction_response_ext = mediation::AuctionResponseExt {
                        auction_id: Some("auction123".to_string()),
                        auction_configuration_id: Some(456),
                        auction_configuration_uid: Some("config789".to_string()),
                        token: Some("key123".to_string()),
                        auction_pricefloor: Some(1.0),
                        auction_timeout: Some(500),
                        external_win_notifications: Some(true),
                        segment: Some(mediation::Segment {
                            id: Some("segment_id".to_string()),
                            uid: Some("segment_uid".to_string()),
                            ext: None,
                        }),
                    };
                    ext.set_extension_data(mediation::AUCTION_RESPONSE_EXT, auction_response_ext)
                        .unwrap();
                    ext
                },
                ..Default::default()
            })),
        };

        Ok(Response::new(response))
    }
}

#[tokio::test]
async fn test_auction() {
    // Create a new extension registry
    let registry = {
        let mut registry = ExtensionRegistry::new();
        registry.register(USER_EXT);
        registry.register(APP_EXT);
        registry.register(DEVICE_EXT);
        registry.register(AUCTION_RESPONSE_EXT);
        registry.register(BID_EXT);
        registry
    };
    // Initialize the registry
    bidon::codec::init_registry(registry).unwrap();

    // Start the mock gRPC server
    let addr = "[::1]:50051".parse().unwrap();
    let mock_service = MockBiddingService::default();

    let server_handle = tokio::spawn(async move {
        Server::builder()
            .add_service(BiddingServiceServer::new(mock_service))
            .serve(addr)
            .await
            .unwrap();
    });

    // Wait a bit for the server to start
    tokio::time::sleep(tokio::time::Duration::from_millis(100)).await;

    // Create the app with gRPC client pointing to our mock server
    let grpc_client = bidon::bidding::ProxyBiddingService::new("http://[::1]:50051")
        .await
        .unwrap();
    let app = create_app(Box::new(grpc_client));

    // Create a test client
    let client = TestClient::new(app);

    // Send a request with the header
    let response = client
        .await
        .post("/v2/auction/banner")
        .header("x-bidon-version", "TestValue")
        .header("X-Forwarded-For", "127.0.0.1")
        .json(&get_auction_request())
        .send()
        .await;

    // Get response text once and store it
    // Assert the response
    assert_eq!(response.status(), StatusCode::OK);
    let response_text = response.text().await;
    // println!("Response: {:?}", response_text);
    assert!(!response_text.is_empty(), "Output must not be empty");

    // Optional: Shutdown the mock server
    server_handle.abort();
}
