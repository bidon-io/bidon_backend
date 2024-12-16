mod infrastructure;

use bidon::org::bidon::proto::v1::mediation::{
    APP_EXT, AUCTION_RESPONSE_EXT, BID_EXT, DEVICE_EXT, USER_EXT,
};
use infrastructure::{config, server};
use prost::ExtensionRegistry;

#[tokio::main]
async fn main() {
    let registry = {
        let mut registry = ExtensionRegistry::new();
        registry.register(USER_EXT);
        registry.register(APP_EXT);
        registry.register(DEVICE_EXT);
        registry.register(AUCTION_RESPONSE_EXT);
        registry.register(BID_EXT);
        registry
    };
    // Initialize the registry. This should be done once in the application at startup.
    bidon::codec::init_registry(registry).unwrap();

    let listen_addr = format!("0.0.0.0:{}", config::settings().port());

    let welcome_string = format!(
        r#"Starting server at: {}
    -> Log Level:   {}
    -> Environment: {}."#,
        listen_addr,
        config::settings().log_level(),
        config::settings().env()
    );
    println!("{}", welcome_string);

    // Start the server
    let listener = tokio::net::TcpListener::bind(listen_addr).await.unwrap();
    server::run(listener).await.unwrap();
}
