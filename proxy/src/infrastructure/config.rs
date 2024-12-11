use dotenvy::EnvLoader;
use std::sync::OnceLock;

#[derive(Debug)]
pub struct Config {
    grpc_url: String,
    env: String,
    log_level: String,
    port: String,
}

impl Config {
    pub fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let env_map = EnvLoader::new().load()?;
        let grpc_url = env_map
            .var("PROXY_GRPC_URL")
            .unwrap_or("http://0.0.0.0:50051".to_string());
        let env = env_map
            .var("ENVIRONMENT")
            .unwrap_or("development".to_string());
        let log_level = env_map.var("PROXY_LOG_LEVEL").unwrap_or("info".to_string());
        let port = env_map.var("PROXY_PORT").unwrap_or("3000".to_string());

        Ok(Config {
            grpc_url,
            env,
            log_level,
            port,
        })
    }

    pub fn grpc_url(&self) -> &str {
        self.grpc_url.as_str()
    }

    pub fn log_level(&self) -> &str {
        self.log_level.as_str()
    }

    pub fn env(&self) -> &str {
        self.env.as_str()
    }

    pub fn port(&self) -> &str {
        self.port.as_str()
    }
}

pub fn settings() -> &'static Config {
    static CONFIG: OnceLock<Config> = OnceLock::new();
    CONFIG.get_or_init(|| Config::new().unwrap())
}
