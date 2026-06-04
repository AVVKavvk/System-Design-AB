-- Run automatically on first container start via docker-entrypoint-initdb.d

CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.page_events
(
    -- Stored as Unix epoch seconds (UInt32).
    -- ClickHouse DateTime columns are perfect for range scans.
    ts        DateTime,

    -- LowCardinality compresses repeated strings (page, browser, country)
    -- into a dictionary — typically 10x smaller on disk and faster to GROUP BY.
    page      LowCardinality(String),
    browser   LowCardinality(String),
    country   LowCardinality(String),
    referrer  LowCardinality(String),

    duration  UInt16
)
ENGINE = MergeTree()

-- PARTITION BY day so old data can be dropped cheaply with
-- ALTER TABLE DROP PARTITION '20240101' instead of expensive deletes.
PARTITION BY toYYYYMMDD(ts)

-- PRIMARY KEY / ORDER BY on ts means the 1-minute window query
-- resolves to a single index lookup + small sequential read.
ORDER BY ts

-- Keeps 1 year of data before TTL auto-deletes old partitions.
TTL ts + INTERVAL 1 YEAR DELETE

SETTINGS
    -- index_granularity controls how many rows per primary key mark.
    -- 8192 is the default; lower = more precise seeks, higher memory.
    index_granularity = 8192;