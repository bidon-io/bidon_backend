use crate::config::settings;
use axum::{http::StatusCode, response::IntoResponse, routing::get, Router};
use axum_prometheus::PrometheusMetricLayerBuilder;
use bidon::bidding::ProxyBiddingService;

pub async fn run(listener: tokio::net::TcpListener) -> Result<(), Box<dyn std::error::Error>> {
    let bidding_service = ProxyBiddingService::new(settings().grpc_url()).await?;

    println!("Connecting to the gRPC server at {}", settings().grpc_url());

    let mut app = bidon::create_app(Box::new(bidding_service));

    let (prometheus_layer, metric_handle) = PrometheusMetricLayerBuilder::new()
        .with_ignore_patterns(&["/metrics"])
        .with_prefix("bidon_proxy")
        .with_default_metrics()
        .build_pair();

    // Add the prometheus layer to the app
    app = app.layer(prometheus_layer);

    // Create the metrics app
    let metrics_app =
        Router::new().route("/metrics", get(|| async move { metric_handle.render() }));

    // Add the health check route to the app
    app = app.route("/health", get(health));

    // Spawn the metrics server on port 9095
    tokio::spawn(async move {
        let metrics_listener = tokio::net::TcpListener::bind((
            std::net::Ipv4Addr::new(0, 0, 0, 0),
            settings().metrics_port(),
        ))
        .await
        .expect("Failed to bind metrics port");
        axum::serve(metrics_listener, metrics_app)
            .await
            .expect("Metrics server error");
    });

    axum::serve(
        listener,
        app.into_make_service_with_connect_info::<std::net::SocketAddr>(),
    )
    .await
    .unwrap();

    Ok(())
}

async fn health() -> impl IntoResponse {
    match tonic::transport::Endpoint::new(settings().grpc_url()) {
        Ok(endpoint) => match endpoint.connect().await {
            Ok(_) => StatusCode::OK,
            Err(_) => StatusCode::SERVICE_UNAVAILABLE,
        },
        Err(_) => StatusCode::SERVICE_UNAVAILABLE,
    }
}
