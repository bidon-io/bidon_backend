use anyhow::Result;
use dotenvy::{EnvLoader, EnvMap};
use std::env;
use std::sync::OnceLock;

#[derive(Debug)]
pub struct Config {
    grpc_url: String,
    env: String,
    log_level: String,
    port: String,
    metrics_port: u16,
    sentry_dsn: String,
}

impl Config {
    pub fn new() -> Result<Self> {
        let env_map = Self::load_env_map();
        let grpc_url = env_map
            .var("PROXY_GRPC_URL")
            .unwrap_or("http://0.0.0.0:50051".to_string());
        let env = env_map
            .var("ENVIRONMENT")
            .unwrap_or("development".to_string());
        let metrics_port = env_map
            .var("PROXY_METRICS_PORT")
            .unwrap_or("9095".to_string());

        Ok(Config {
            grpc_url,
            env,
            log_level: env_map.var("PROXY_LOG_LEVEL").unwrap_or("info".to_string()),
            port: env_map.var("PROXY_PORT").unwrap_or("3000".to_string()),
            metrics_port: metrics_port.parse::<u16>().unwrap(),
            sentry_dsn: env_map.var("SENTRY_DSN").unwrap_or("".to_string()),
        })
    }

    fn load_env_map() -> EnvMap {
        let path = ".env";

        if std::path::Path::new(path).exists() {
            EnvLoader::new().load().unwrap()
        } else {
            env::vars().collect()
        }
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

    pub fn metrics_port(&self) -> u16 {
        self.metrics_port
    }

    pub fn sentry_dsn(&self) -> &str {
        self.sentry_dsn.as_str()
    }
}

pub fn settings() -> &'static Config {
    static CONFIG: OnceLock<Config> = OnceLock::new();
    CONFIG.get_or_init(|| Config::new().unwrap())
}
