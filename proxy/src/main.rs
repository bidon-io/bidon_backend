mod infrastructure;

use infrastructure::{config, server};

#[tokio::main]
async fn main() {
    let listen_addr = format!("0.0.0.0:{}", config::settings().port());

    let welcome_string = format!(r#"Starting server at: {}
    -> Log Level:   {}
    -> Environment: {}."#, listen_addr, config::settings().log_level(), config::settings().env());
    println!("{}", welcome_string);

    // Start the server
    let listener = tokio::net::TcpListener::bind(listen_addr).await.unwrap();
    server::run(listener).await.unwrap();
}
