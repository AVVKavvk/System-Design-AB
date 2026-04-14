use axum::{
    Router,
    routing::{get, post},
};

use crate::{controllers::clickhouse_controller, state::app_state::AppState};

pub fn get_clickhouse_route() -> Router<AppState> {
    let ch_route = Router::new()
        .route("/ingest-event", post(clickhouse_controller::ingest_events))
        .route("/traffic", get(clickhouse_controller::get_traffic))
        .route("/health", get(clickhouse_controller::health));

    ch_route
}
