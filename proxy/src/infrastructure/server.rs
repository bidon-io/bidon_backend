use crate::config::settings;
use bidon::bidding::ProxyBiddingService;
use tracing_subscriber::{
    filter::LevelFilter, fmt, layer::SubscriberExt, util::SubscriberInitExt, EnvFilter,
};

pub async fn run(listener: tokio::net::TcpListener) -> Result<(), Box<dyn std::error::Error>> {
    init_logger();

    let bidding_service = ProxyBiddingService::new(settings().grpc_url()).await?;

    println!("Connecting to the gRPC server at {}", settings().grpc_url());

    let app = bidon::create_app(Box::new(bidding_service));

    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .await
    .unwrap();

    Ok(())
}

fn init_logger() {
    tracing_subscriber::registry()
        .with(
            EnvFilter::builder()
                .with_default_directive(LevelFilter::ERROR.into())
                .parse_lossy(settings().log_level()),
        )
        .with(fmt::layer())
        .init();
}
