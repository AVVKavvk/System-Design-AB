#[derive(Debug, Clone)]
pub struct Config {
    pub host: String,
    pub port: u16,
    pub clickhouse_url: String,
    pub clickhouse_db: String,
    pub clickhouse_user: String,
    pub clickhouse_password: String,
}
use anyhow::{Context, Result}; // Import both

impl Config {
    pub fn from_env() -> Result<Self> {
        // Return a Result
        Ok(Self {
            host: std::env::var("APP_HOST").unwrap_or_else(|_| "0.0.0.0".into()),
            port: std::env::var("APP_PORT")
                .unwrap_or_else(|_| "8080".into())
                .parse::<u16>()
                .context("APP_PORT must be a valid port number")?, // Now this works!

            clickhouse_url: std::env::var("CLICKHOUSE_URL")
                .unwrap_or_else(|_| "http://localhost:8123".into()),
            clickhouse_db: std::env::var("CLICKHOUSE_DB").unwrap_or_else(|_| "analytics".into()),
            clickhouse_user: std::env::var("CLICKHOUSE_USER").unwrap_or_else(|_| "default".into()),
            clickhouse_password: std::env::var("CLICKHOUSE_PASSWORD").unwrap_or_else(|_| "".into()),
        })
    }
}
