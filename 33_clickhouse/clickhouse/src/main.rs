mod config;
mod controllers;
mod errors;
mod models;
mod repositories;
mod routes;
mod services;
mod state;

use anyhow::Context;
use std::sync::Arc;

use clickhouse::Client;

use crate::{
    config::Config, repositories::clickhouse::CHRespository,
    services::clickhouse::ClickhouseServiceImpl, state::app_state::AppState,
};
use tracing_subscriber::{EnvFilter, fmt, prelude::*};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    tracing_subscriber::registry()
        .with(fmt::layer())
        .with(EnvFilter::try_from_default_env().unwrap_or_else(|_| "info".into()))
        .init();

    tracing::info!("Starting server");

    // ── 1. Load .env file
    dotenvy::dotenv().ok();

    // Load typed config
    let cfg = Config::from_env().context("Failed to load configuration")?;

    // Build ClickHouse client — connection is reused across all requests
    let ch_client = Client::default()
        .with_url(&cfg.clickhouse_url)
        .with_database(&cfg.clickhouse_db)
        .with_user(&cfg.clickhouse_user)
        .with_password(&cfg.clickhouse_password);

    //  Dependency Injection

    let ch_repo = Arc::new(CHRespository::new(ch_client));
    let ch_service = Arc::new(ClickhouseServiceImpl::new(ch_repo));
    let state = AppState::new(ch_service);

    // ── 7. Build router
    let app = routes::create_router(state);
    // ── 8. Start server
    let addr = format!("{}:{}", cfg.host, cfg.port);
    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .context(format!("Failed to bind to {}", addr))?;

    tracing::info!(address = %addr, "Server listening");

    axum::serve(listener, app)
        .with_graceful_shutdown(shutdown_signal())
        .await
        .context("Server error")?;

    Ok(())
}

/// Listens for Ctrl-C / SIGTERM and triggers graceful shutdown.
async fn shutdown_signal() {
    use tokio::signal;

    let ctrl_c = async {
        signal::ctrl_c()
            .await
            .expect("failed to install Ctrl+C handler");
    };

    #[cfg(unix)]
    let terminate = async {
        signal::unix::signal(signal::unix::SignalKind::terminate())
            .expect("failed to install SIGTERM handler")
            .recv()
            .await;
    };

    #[cfg(not(unix))]
    let terminate = std::future::pending::<()>();

    tokio::select! {
        _ = ctrl_c    => tracing::info!("received Ctrl-C, shutting down"),
        _ = terminate => tracing::info!("received SIGTERM, shutting down"),
    }
}
