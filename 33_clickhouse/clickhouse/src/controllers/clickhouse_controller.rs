use axum::{
    Json,
    extract::{Query, State},
    response::IntoResponse,
};

use crate::{
    errors::ApiError,
    models::{EventInput, TrafficQuery},
    state::app_state::AppState,
};

pub async fn get_traffic(
    State(state): State<AppState>,
    Query(params): Query<TrafficQuery>,
) -> Result<impl IntoResponse, ApiError> {
    let ch_server = state.ch_service;

    let result = ch_server.get_traffic(params).await?;

    Ok(result)
}
pub async fn ingest_events(
    State(state): State<AppState>,
    Json(events): Json<Vec<EventInput>>,
) -> Result<impl IntoResponse, ApiError> {
    let ch_server = state.ch_service;

    let result = ch_server.ingest_events(events).await?;

    Ok(result)
}

pub async fn health(State(state): State<AppState>) -> Result<impl IntoResponse, ApiError> {
    let ch_server = state.ch_service;

    let result = ch_server.health_check().await?;

    Ok(result)
}
