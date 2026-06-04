pub mod clickhouse_route;

pub use clickhouse_route::get_clickhouse_route;

use axum::Router;
use tower_http::{
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};

use crate::state::app_state::AppState;

/// Build the complete application router.
pub fn create_router(state: AppState) -> Router {
    let ch_router = get_clickhouse_route();

    Router::new()
        .nest("/api/v1/clickhouse", ch_router)
        // Inject shared state into all handlers
        .with_state(state)
        // ── Middleware stack (outermost = first to run) ──────────────────────
        // Automatic HTTP tracing (logs method, path, status, latency)
        .layer(TraceLayer::new_for_http())
        // CORS — open for development; tighten for production
        .layer(
            CorsLayer::new()
                .allow_origin(Any)
                .allow_methods(Any)
                .allow_headers(Any),
        )
}
