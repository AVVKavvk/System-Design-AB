use crate::services::clickhouse::ClickhouseService;
use std::sync::Arc;

/// Shared application state.
/// `Clone` is required because Axum clones State for each request.
/// Since Arc is cheap to clone (just increments a counter), this is fine.
#[derive(Clone)]
pub struct AppState {
    pub ch_service: Arc<dyn ClickhouseService>,
}

impl AppState {
    pub fn new(ch_service: Arc<dyn ClickhouseService>) -> Self {
        Self { ch_service }
    }
}
