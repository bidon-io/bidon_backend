use crate::config::settings;
use bidon::bidding::ProxyBiddingService;

pub async fn run(listener: tokio::net::TcpListener) -> Result<(), Box<dyn std::error::Error>> {
    tracing::info!("Connecting to the gRPC server at {}", settings().grpc_url());
    let bidding_service = ProxyBiddingService::new(settings().grpc_url()).await?;

    let app = bidon::create_app(Box::new(bidding_service));

    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .await
    .unwrap();

    Ok(())
}
