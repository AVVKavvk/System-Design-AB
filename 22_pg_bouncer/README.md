# PgBouncer + PostgreSQL + Go — Connection Pooling

A minimal demo showing how **PgBouncer** sits between a Go web server and PostgreSQL to pool connections efficiently.

## Stack

| Service    | Image                | Port |
| ---------- | -------------------- | ---- |
| PostgreSQL | `postgres:15-alpine` | 5432 |
| PgBouncer  | `edoburu/pgbouncer`  | 6432 |
| Go App     | `echo` + `pgx`       | 8080 |

## How It Works

```
Go App (echo)
     │
     │  postgres://user:password@pgbouncer:6432/mydatabase
     ▼
 PgBouncer  ◄──── pool of N server connections
     │
     │  postgres://user:password@db:5432/mydatabase
     ▼
 PostgreSQL
```

Instead of the app opening a new Postgres connection per request (expensive — ~10MB RAM each), it connects to PgBouncer, which maintains a small pool of real backend connections and reuses them across many clients.

## Quickstart

```bash
docker compose up --build
```

Test the endpoint:

```bash
curl http://localhost:8080/ping
# {"message":"Hello"}
```

## Endpoints

| Method | Path    | Description                                                  |
| ------ | ------- | ------------------------------------------------------------ |
| GET    | `/ping` | Runs `SELECT pg_sleep(2)` via PgBouncer + returns a greeting |
| GET    | `/test` | Additional test handler (fire multiple request concurrent)   |

## Pool Configuration

| Parameter           | Value         | Description                                |
| ------------------- | ------------- | ------------------------------------------ |
| `POOL_MODE`         | `transaction` | Connection released after each transaction |
| `MAX_CLIENT_CONN`   | `100`         | Max simultaneous client connections        |
| `DEFAULT_POOL_SIZE` | `5`           | Max real Postgres connections per pool     |

## Observing the Pool

Connect to the PgBouncer admin console:

```bash
psql "postgresql://user:password@localhost:6432/pgbouncer"
```

Then:

```sql
SHOW POOLS;   -- see cl_active, sv_active, cl_waiting
SHOW STATS;   -- query counts, data transferred
SHOW CLIENTS; -- connected clients
SHOW SERVERS; -- backend connections
```

Watch live while load testing:

```bash
watch -n 1 'psql "postgresql://user:password@localhost:6432/pgbouncer" -c "SHOW POOLS;"'
```

Drive load from another terminal:

```bash
seq 100 | xargs -P 20 -I{} curl -s http://localhost:8080/ping > /dev/null
```

You'll see `sv_active` cap at `DEFAULT_POOL_SIZE` (5) while `cl_active` and `cl_waiting` grow — PgBouncer queuing clients rather than opening new backend connections.

## Key Insight

| Without PgBouncer                       | With PgBouncer                        |
| --------------------------------------- | ------------------------------------- |
| 100 requests = 100 Postgres connections | 100 requests = 5 Postgres connections |
| ~1GB RAM for connections                | ~50MB RAM for connections             |
| Connection setup overhead per request   | Connections reused across requests    |

## Environment Variables

| Variable                  | Description                           |
| ------------------------- | ------------------------------------- |
| `PG_BOUNCER_DATABASE_URL` | Full DSN the Go app uses to connect   |
| `DB_HOST`                 | Postgres hostname (seen by PgBouncer) |
| `DB_PORT`                 | Postgres port                         |
| `DB_USER`                 | Database user                         |
| `DB_PASSWORD`             | Database password                     |
| `DB_NAME`                 | Database name                         |

## Pool Modes

PgBouncer supports three pool modes, each controlling **when a server connection is released back to the pool**.

### `session` (safest)

A server connection is assigned to a client **for the entire duration of the session** (until the client disconnects).

- Behaves exactly like a direct Postgres connection
- All Postgres features work (prepared statements, `SET`, advisory locks, `LISTEN/NOTIFY`)
- **Least efficient** — connection held even when client is idle
- **Use when:** migrating an existing app to PgBouncer with zero code changes, or when your app relies on session-level state

---

### `transaction` (recommended for most apps)

A server connection is assigned only **for the duration of a single transaction**, then released back to the pool.

- Most efficient for typical web apps (request → query → release)
- Breaks features that rely on session state:
  - Prepared statements (`PREPARE / EXECUTE`)
  - `SET` statements (e.g. `SET search_path`)
  - Advisory locks
  - `LISTEN / NOTIFY`
  - `WITH HOLD` cursors
- **Use when:** your app uses short, discrete transactions (REST APIs, microservices)

---

### `statement` (most aggressive, rarely used)

A server connection is released after **every single statement**, even mid-transaction.

- Multi-statement transactions (`BEGIN ... COMMIT`) will **break**
- Only safe for single-statement, auto-commit workloads
- Rarely used in practice
- **Use when:** read-only analytics queries, or simple single-statement fire-and-forget workloads

---

### Comparison

| Feature                      | `session` | `transaction` | `statement` |
| ---------------------------- | :-------: | :-----------: | :---------: |
| Connection reuse             |    Low    |     High      |   Highest   |
| Prepared statements          |    ✅     |      ❌       |     ❌      |
| Multi-statement transactions |    ✅     |      ✅       |     ❌      |
| Session-level `SET`          |    ✅     |      ❌       |     ❌      |
| `LISTEN / NOTIFY`            |    ✅     |      ❌       |     ❌      |
| Advisory locks               |    ✅     |      ❌       |     ❌      |
| Typical web API              | ⚠️ works  |   ✅ ideal    |  ❌ avoid   |

> **Rule of thumb:** Start with `transaction`. Drop to `session` only if your ORM or app breaks. Never use `statement` unless you have a very specific use case.

---

## SHOW POOLS — Column Reference

```sql
psql "postgresql://user:password@localhost:6432/pgbouncer" -c "SHOW POOLS;"
```

### Client-side columns (connections from our app → PgBouncer)

| Column                  | What it means                                                             |
| ----------------------- | ------------------------------------------------------------------------- |
| `database`              | The PgBouncer database name the pool belongs to                           |
| `user`                  | The Postgres user for this pool                                           |
| `cl_active`             | Clients currently linked to a server connection and executing a query     |
| `cl_waiting`            | Clients waiting for a free server connection (pool is full)               |
| `cl_active_cancel_req`  | Clients that sent a cancel request and are waiting for it to be forwarded |
| `cl_waiting_cancel_req` | Clients waiting to send a cancel request to the server                    |

### Server-side columns (connections from PgBouncer → Postgres)

| Column              | What it means                                                                     |
| ------------------- | --------------------------------------------------------------------------------- |
| `sv_active`         | Server connections currently assigned to a client and doing work                  |
| `sv_active_cancel`  | Server connections currently handling a cancel request                            |
| `sv_being_canceled` | Server connections that received a cancel and are being cleaned up                |
| `sv_idle`           | Server connections open but idle, ready to be assigned to the next client         |
| `sv_used`           | Server connections that were idle and not yet tested for liveness                 |
| `sv_tested`         | Server connections currently being tested with a liveness check                   |
| `sv_login`          | Server connections currently being established (TCP handshake / auth in progress) |

### Timing & config columns

| Column               | What it means                                                          |
| -------------------- | ---------------------------------------------------------------------- |
| `maxwait`            | How long (seconds) the oldest `cl_waiting` client has been waiting     |
| `maxwait_us`         | Same as `maxwait` but in microseconds (higher precision)               |
| `pool_mode`          | The pool mode for this pool (`session`, `transaction`, or `statement`) |
| `load_balance_hosts` | Whether load balancing across multiple backend hosts is enabled        |

### What to watch during load testing

```
cl_active   → should rise as requests come in
sv_active   → caps at DEFAULT_POOL_SIZE (5 in this demo)
cl_waiting  → grows when sv_active hits the cap — clients queuing
maxwait     → if this climbs, your pool size is too small for the load
sv_idle     → healthy sign — connections ready and waiting
sv_login    → should be 0 at rest; spikes mean pool is scaling up
```

### The key relationship

```
cl_active (20 clients)
      │
      ▼
PgBouncer pool
      │
      ├── sv_active (5 real connections doing work)
      ├── sv_idle   (0 free — pool is saturated)
      │
      └── cl_waiting (15 clients queued) ← maxwait starts climbing here
```

> If `cl_waiting` is consistently non-zero and `maxwait` is growing,
> increase `DEFAULT_POOL_SIZE` — but remember each unit is a real Postgres connection.

!["output.png"](./img/image.png)

### `statement` (most aggressive, rarely used)

A server connection is released back to the pool **after every single SQL statement**,
even if you are inside a `BEGIN ... COMMIT` block.

PgBouncer does not wait for your transaction to finish.
The moment one statement completes, the connection is gone.

---

#### Why this breaks transactions

```sql
BEGIN;                        -- connection A assigned
  UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- connection A
  -- PgBouncer releases connection A here ← mid transaction!
  UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- connection B assigned
COMMIT;                       -- connection B commits — but BEGIN was on connection A
```

Connection B has no idea a transaction was started.
It just sees a standalone `UPDATE` and auto-commits it.
**The debit happened. The credit may not.**
Your data is now inconsistent.

---

#### The only safe workload

Single, self-contained, read-only queries that need no transaction context:

```sql
-- Safe ✅ — one statement, no state, no transaction needed
SELECT count(*) FROM events WHERE created_at > now() - interval '1 day';

-- Safe ✅ — fire and forget, result is independent
SELECT AVG(response_time_ms) FROM request_logs WHERE endpoint = '/api/checkout';

-- Unsafe ❌ — two statements that must be atomic
INSERT INTO orders (...) VALUES (...);
INSERT INTO order_items (...) VALUES (...);
```

---

#### Real-world use case (the rare one)

An **internal metrics scraper** that runs isolated `SELECT` queries every 30 seconds
against a read replica to feed a dashboard:

```
Grafana / Prometheus exporter
        │
        │  SELECT sum(revenue) FROM orders WHERE ...
        │  SELECT count(*) FROM active_sessions
        │  SELECT avg(latency_ms) FROM request_logs
        ▼
   PgBouncer (statement mode)
        │
        ▼
   Postgres read replica
```

Each query is independent. No transactions. No shared state.
`statement` mode works perfectly here and squeezes maximum connection reuse.

---

#### Why nobody uses it for application code

Most ORMs and query builders — GORM, SQLAlchemy, ActiveRecord, pgx —
wrap operations in implicit transactions under the hood.

```go
// Looks like one statement. Actually three under the hood.
db.Create(&order)
// BEGIN
// INSERT INTO orders ...
// INSERT INTO order_items ...   ← PgBouncer cuts the connection here
// COMMIT                        ← never reaches Postgres
```

You would not see an error. You would just lose data silently.

---

#### Summary

| Question                                       | Answer                    |
| ---------------------------------------------- | ------------------------- |
| Is your query a single `SELECT` with no state? | ✅ `statement` is fine    |
| Do you use `BEGIN / COMMIT` anywhere?          | ❌ do not use `statement` |
| Does your ORM manage transactions for you?     | ❌ do not use `statement` |
| Are you on a read replica scraping metrics?    | ✅ reasonable use case    |

> **Bottom line:** If you are asking whether you should use `statement` mode,
> the answer is almost certainly no. Use `transaction` instead.
