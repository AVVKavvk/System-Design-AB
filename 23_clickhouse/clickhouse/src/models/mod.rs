use chrono::DateTime;
use chrono::Utc;
use clickhouse::Row;
use serde::{Deserialize, Serialize};

/// One page-view event as it arrives from the client.
#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct EventInput {
    /// ISO-8601 timestamp, e.g. "2024-01-13T10:05:00Z"
    pub ts: DateTime<Utc>,
    pub page: String,
    pub browser: String,
    pub country: String,
    pub referrer: String,
    /// Seconds the user spent on the page
    pub duration: u16,
}

/// Row shape that ClickHouse expects — Unix timestamp as u32.
/// The `clickhouse::Row` derive handles serialization to the native format.
#[derive(Debug, Serialize, Row)]
pub struct EventRow {
    /// ClickHouse DateTime is seconds since epoch (u32).
    ts: u32,
    page: String,
    browser: String,
    country: String,
    referrer: String,
    duration: u16,
}

impl From<EventInput> for EventRow {
    fn from(input: EventInput) -> Self {
        EventRow {
            ts: input.ts.timestamp() as u32,
            page: input.page,
            browser: input.browser,
            country: input.country,
            referrer: input.referrer,
            duration: input.duration,
        }
    }
}

#[derive(Debug, Deserialize)]
pub struct TrafficQuery {
    pub ts: DateTime<Utc>,
}

#[derive(Debug, Deserialize, Row)]
pub struct TrafficRow {
    pub count: u64,
}
