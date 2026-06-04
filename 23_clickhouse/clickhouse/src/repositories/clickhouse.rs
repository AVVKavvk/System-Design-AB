use async_trait::async_trait;
use axum::{
    Json,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use clickhouse::Client;
use std::time::Instant;
use tracing::{error, info};

use crate::{
    errors::ApiError,
    models::{EventInput, EventRow, TrafficQuery, TrafficRow},
};

#[async_trait]
pub trait ClickhouseRespository: Send + Sync {
    async fn ingest_events(&self, events: Vec<EventInput>) -> Result<(), ApiError>;
    async fn get_traffic(&self, query: TrafficQuery) -> Result<Response, ApiError>;
    async fn health_check(&self) -> Result<Response, ApiError>;
}

#[derive(Clone)]
pub struct CHRespository {
    ch: Client,
}

impl CHRespository {
    pub fn new(ch: Client) -> Self {
        Self { ch }
    }
}

#[async_trait]
impl ClickhouseRespository for CHRespository {
    async fn ingest_events(&self, events: Vec<EventInput>) -> Result<(), ApiError> {
        let count = events.len();
        let started = Instant::now();

        let mut page_events_inserter = self
            .ch
            .inserter("page_events")?
            .with_max_rows(count as u64 + 1);

        for event in events {
            page_events_inserter.write(&EventRow::from(event))?;
        }

        page_events_inserter.end().await?;

        let elapsed_ms = started.elapsed().as_millis();
        info!(rows = count, elapsed = elapsed_ms, "inserted batch");

        Ok(())
    }

    // Change 'impl IntoResponse' to 'Response' to match the trait
    async fn get_traffic(&self, query: TrafficQuery) -> Result<Response, ApiError> {
        let window_start = query.ts.timestamp() as u32;
        let window_end = window_start + 60;

        let started = Instant::now();

        let rows: Vec<TrafficRow> = self
            .ch
            .query("SELECT count() AS count FROM page_events WHERE ts >= ? AND ts < ?")
            .bind(window_start)
            .bind(window_end)
            .fetch_all()
            .await?;

        let count = rows.first().map(|r| r.count).unwrap_or(0);
        let elapsed_ms = started.elapsed().as_millis();

        info!(
            window_start = query.ts.to_rfc3339(),
            count = count,
            elapsed = elapsed_ms,
            "traffic query"
        );

        // Added .into_response()
        Ok(Json(serde_json::json!({
            "count": count,
            "elapsed_ms": elapsed_ms,
        }))
        .into_response())
    }

    // Change 'impl IntoResponse' to 'Response' to match the trait
    async fn health_check(&self) -> Result<Response, ApiError> {
        match self.ch.query("SELECT 1").fetch_one::<u8>().await {
            // FIX: Added .into_response() to both arms
            Ok(_) => {
                Ok((StatusCode::OK, Json(serde_json::json!({ "status": "ok" }))).into_response())
            }
            Err(e) => {
                error!("health check failed: {e}");
                Ok((
                    StatusCode::SERVICE_UNAVAILABLE,
                    Json(serde_json::json!({ "status": "degraded" })),
                )
                    .into_response())
            }
        }
    }
}
