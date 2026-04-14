use std::sync::Arc;

use async_trait::async_trait;
use axum::response::{IntoResponse, Response};

use crate::{
    errors::ApiError,
    models::{EventInput, TrafficQuery},
    repositories::clickhouse::ClickhouseRespository,
};

#[async_trait]
pub trait ClickhouseService: Send + Sync {
    async fn ingest_events(&self, events: Vec<EventInput>) -> Result<Response, ApiError>;
    async fn get_traffic(&self, query: TrafficQuery) -> Result<Response, ApiError>;
    async fn health_check(&self) -> Result<Response, ApiError>;
}

pub struct ClickhouseServiceImpl {
    pub repo: Arc<dyn ClickhouseRespository>,
}

impl ClickhouseServiceImpl {
    pub fn new(repo: Arc<dyn ClickhouseRespository>) -> Self {
        Self { repo }
    }
}

#[async_trait]
impl ClickhouseService for ClickhouseServiceImpl {
    async fn ingest_events(&self, events: Vec<EventInput>) -> Result<Response, ApiError> {
        self.repo.ingest_events(events).await?;

        // Return a simple 200 OK Json response
        Ok(axum::Json(serde_json::json!({ "status": "success" })).into_response())
    }

    async fn get_traffic(&self, query: TrafficQuery) -> Result<Response, ApiError> {
        // This works because repo.get_traffic already returns Result<Response, ApiError>
        self.repo.get_traffic(query).await
    }

    async fn health_check(&self) -> Result<Response, ApiError> {
        // This works because repo.health_check already returns Result<Response, ApiError>
        self.repo.health_check().await
    }
}
