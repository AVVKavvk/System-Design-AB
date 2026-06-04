# Analytics — Rust + Axum + ClickHouse

High-throughput page-event ingestion API with a 1-minute traffic query endpoint.

## Quick start

### 1. Start ClickHouse

```bash
docker compose up -d
# wait ~10s for the health check to pass
docker compose ps
```

### 2. Configure

```bash
cp .env.example .env
# edit .env if your ClickHouse is remote
```

### 3. Run the server

```bash
cargo run --release --bin server
# → listening on 0.0.0.0:8080
```

### 4. Run the load tester (1 million events)

```bash
cargo run --release --bin load_test
```

Or with custom parameters:

```bash
cargo run --release --bin load_test -- \
  --url     http://localhost:8080 \
  --total   1000000 \
  --workers 100 \
  --batch   1000
```

---

## API

### POST api/v1/clickhouse/ingest-event

Insert a batch of page-view events.

**Request body** — JSON array:

```json
[
  {
    "ts": "2024-01-13T10:05:00Z",
    "page": "/pricing",
    "browser": "Chrome",
    "country": "IN",
    "referrer": "google.com",
    "duration": 42
  }
]
```

**Response**:

```json
{ "inserted": 500, "elapsed_ms": 12 }
```

---

### GET api/v1/clickhouse/traffic?ts=2024-01-13T10:05:00Z

Returns event count for the 1-minute window `[ts, ts + 60s)`.

**Response**:

```json
{
  "window_start": "2024-01-13T10:05:00Z",
  "window_end": "2024-01-13T10:06:00Z",
  "count": 48291,
  "elapsed_ms": 3
}
```

---

### GET api/v1/clickhouse/health

Liveness probe — also pings ClickHouse.

---

## Architecture notes

### Why batch inserts?

ClickHouse is optimised for large sequential writes. Row-by-row inserts
create one "part" per insert on disk — at high throughput this
overwhelms the merge background process. Batching 500+ rows per request
keeps part creation sane.

Rule of thumb: **never insert fewer than 1 000 rows at a time** in production.
Use a buffer (in-process queue, Kafka, or ClickHouse's `Buffer` engine)
if your upstream sends individual events.

### Why `LowCardinality(String)`?

For columns with fewer than ~10 000 distinct values (page, browser, country),
ClickHouse stores a per-column dictionary and replaces values with small
integer indices. This:

- Reduces disk usage ~5–10x for those columns
- Speeds up `GROUP BY` (compare integers, not strings)
- Costs nothing extra at write time

### Why `ORDER BY ts`?

The `GET /traffic` query does:

```sql
WHERE ts >= ? AND ts < ?
```

Because `ts` is the primary sort key, ClickHouse resolves this with
a binary search over the sparse primary index to find the first matching
granule, then reads only the granules that overlap the range.
At 1 million rows with default `index_granularity = 8192` that's
~122 index marks — the lookup is essentially free.

### Connection pooling

The `clickhouse::Client` is wrapped in `Arc` and cloned into every Axum
handler. The underlying HTTP client (reqwest) maintains a connection pool
— connections are reused across requests without any extra setup.

### Graceful shutdown

`axum::serve(...).with_graceful_shutdown(shutdown_signal())` drains
in-flight requests before the process exits. Safe to `kill -SIGTERM` in
production / Kubernetes.
