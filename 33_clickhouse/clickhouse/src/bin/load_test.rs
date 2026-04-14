/// Load tester — hammers POST /events with 1 million events
/// using concurrent Tokio workers and batched payloads.
///
/// Usage:
///   cargo run --bin load_test -- --url http://localhost:8080 --total 1000000 --workers 50 --batch 500
///
/// Env vars (alternative to flags):
///   LT_URL, LT_TOTAL, LT_WORKERS, LT_BATCH
use std::{
    sync::{
        Arc,
        atomic::{AtomicU64, Ordering},
    },
    time::{Duration, Instant},
};

use chrono::{DateTime, Utc};
use reqwest::Client;
use serde::{Deserialize, Serialize};
use tokio::time::sleep;

// ─────────────────────────────────────────────
// CLI config
// ─────────────────────────────────────────────

#[derive(Debug, Clone)]
struct Config {
    url: String,
    total: u64,
    workers: u64,
    batch: u64,
}

impl Config {
    fn from_env_or_defaults() -> Self {
        let args: Vec<String> = std::env::args().collect();
        let get = |flag: &str, env: &str, default: &str| -> String {
            // check --flag value first
            args.windows(2)
                .find(|w| w[0] == flag)
                .map(|w| w[1].clone())
                .or_else(|| std::env::var(env).ok())
                .unwrap_or_else(|| default.to_string())
        };

        Self {
            url: get("--url", "LT_URL", "http://localhost:8080"),
            total: get("--total", "LT_TOTAL", "1000000").parse().unwrap(),
            workers: get("--workers", "LT_WORKERS", "50").parse().unwrap(),
            batch: get("--batch", "LT_BATCH", "500").parse().unwrap(),
        }
    }
}

// ─────────────────────────────────────────────
// Payload types (must match server's EventInput)
// ─────────────────────────────────────────────

#[derive(Serialize, Clone)]
struct EventInput {
    ts: DateTime<Utc>,
    page: String,
    browser: String,
    country: String,
    referrer: String,
    duration: u16,
}

static PAGES: &[&str] = &[
    "/home",
    "/pricing",
    "/docs",
    "/blog",
    "/signup",
    "/login",
    "/features",
];
static BROWSERS: &[&str] = &["Chrome", "Safari", "Firefox", "Edge"];
static COUNTRIES: &[&str] = &["IN", "US", "DE", "GB", "BR", "FR", "CA"];
static REFERRERS: &[&str] = &["google.com", "direct", "twitter.com", "github.com"];

fn rand_from<T: Copy>(slice: &[T], seed: u64) -> T {
    slice[(seed as usize) % slice.len()]
}

/// Build a batch of events.
/// We use a simple xorshift64 so the load tester itself has zero overhead
/// from a real RNG — keeps the benchmark honest.
fn make_batch(size: u64, worker_id: u64, batch_num: u64) -> Vec<EventInput> {
    let mut seed = worker_id.wrapping_mul(0xdeadbeef) ^ batch_num.wrapping_mul(0x12345678);
    let now = Utc::now();

    (0..size)
        .map(|i| {
            // xorshift64
            seed ^= seed << 13;
            seed ^= seed >> 7;
            seed ^= seed << 17;
            let s2 = seed ^ i;

            EventInput {
                ts: now - chrono::Duration::seconds((s2 % 3600) as i64),
                page: rand_from(PAGES, seed).to_string(),
                browser: rand_from(BROWSERS, s2).to_string(),
                country: rand_from(COUNTRIES, seed ^ s2).to_string(),
                referrer: rand_from(REFERRERS, s2 ^ (seed >> 3)).to_string(),
                duration: ((seed % 235) + 5) as u16,
            }
        })
        .collect()
}

// ─────────────────────────────────────────────
// Shared stats
// ─────────────────────────────────────────────

#[derive(Default)]
struct Stats {
    sent: AtomicU64,
    errors: AtomicU64,
    /// Total HTTP round-trip millis (for avg latency)
    lat_sum: AtomicU64,
    lat_count: AtomicU64,
}

// ─────────────────────────────────────────────
// Worker
// ─────────────────────────────────────────────

async fn worker(
    id: u64,
    batches: u64, // how many batches this worker should send
    batch_sz: u64,
    url: String,
    client: Client,
    stats: Arc<Stats>,
) {
    let endpoint = format!("{url}/api/v1/clickhouse/ingest-event");

    for batch_num in 0..batches {
        let payload = make_batch(batch_sz, id, batch_num);
        let t0 = Instant::now();

        match client.post(&endpoint).json(&payload).send().await {
            Ok(resp) => {
                let lat = t0.elapsed().as_millis() as u64;
                stats.lat_sum.fetch_add(lat, Ordering::Relaxed);
                stats.lat_count.fetch_add(1, Ordering::Relaxed);

                if resp.status().is_success() {
                    stats.sent.fetch_add(batch_sz, Ordering::Relaxed);
                } else {
                    let status = resp.status();
                    let body = resp.text().await.unwrap_or_default();
                    eprintln!("worker {id} got {status}: {body}");
                    stats.errors.fetch_add(batch_sz, Ordering::Relaxed);
                }
            }
            Err(e) => {
                eprintln!("worker {id} request failed: {e}");
                stats.errors.fetch_add(batch_sz, Ordering::Relaxed);
            }
        }
    }
}

// ─────────────────────────────────────────────
// Progress reporter
// ─────────────────────────────────────────────

async fn reporter(total: u64, stats: Arc<Stats>, started: Instant) {
    let mut last_sent = 0u64;
    let mut last_t = Instant::now();

    loop {
        sleep(Duration::from_secs(1)).await;

        let sent = stats.sent.load(Ordering::Relaxed);
        let errors = stats.errors.load(Ordering::Relaxed);
        let lat_c = stats.lat_count.load(Ordering::Relaxed);
        let lat_s = stats.lat_sum.load(Ordering::Relaxed);

        let delta_t = last_t.elapsed().as_secs_f64();
        let delta_sent = sent.saturating_sub(last_sent);
        let rps = (delta_sent as f64 / delta_t) as u64;
        let pct = (sent as f64 / total as f64 * 100.0) as u32;
        let avg_lat = if lat_c > 0 { lat_s / lat_c } else { 0 };
        let elapsed = started.elapsed().as_secs();

        println!(
            "[{elapsed:>4}s] sent={:<10} ({pct:>3}%)  rate={:<8}/s  errors={:<6}  avg_batch_lat={avg_lat}ms",
            sent.to_string(),
            rps.to_string(),
            errors.to_string(),
        );

        last_sent = sent;
        last_t = Instant::now();

        if sent + errors >= total {
            break;
        }
    }
}

// ─────────────────────────────────────────────
// Entry point
// ─────────────────────────────────────────────

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let _ = dotenvy::dotenv();
    let cfg = Config::from_env_or_defaults();

    println!("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    println!("  ClickHouse load tester");
    println!("  target  : {}", cfg.url);
    println!(
        "  total   : {} events",
        cfg.total
            .to_string()
            .as_str()
            .chars()
            .rev()
            .collect::<Vec<_>>()
            .chunks(3)
            .map(|c| c.iter().collect::<String>())
            .collect::<Vec<_>>()
            .join(",")
            .chars()
            .rev()
            .collect::<String>()
    );
    println!("  workers : {}", cfg.workers);
    println!("  batch   : {} events/req", cfg.batch);
    println!("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    // Validate total is divisible cleanly
    let total_batches = cfg.total / cfg.batch;
    let batches_per_worker = total_batches / cfg.workers;
    let actual_total = batches_per_worker * cfg.workers * cfg.batch;

    println!(
        "  actual events to send: {actual_total}  ({batches_per_worker} batches × {workers} workers × {batch} events)",
        workers = cfg.workers,
        batch = cfg.batch,
    );
    println!();

    // Build a shared HTTP client with a large connection pool.
    // One client shared across all workers — reqwest handles the pool internally.
    let client = Client::builder()
        .pool_max_idle_per_host(cfg.workers as usize)
        .timeout(Duration::from_secs(30))
        .build()?;

    let stats = Arc::new(Stats::default());
    let started = Instant::now();

    // Spawn reporter task
    let reporter_handle = tokio::spawn(reporter(actual_total, Arc::clone(&stats), started));

    // Spawn all workers
    let mut handles = Vec::with_capacity(cfg.workers as usize);
    for id in 0..cfg.workers {
        handles.push(tokio::spawn(worker(
            id,
            batches_per_worker,
            cfg.batch,
            cfg.url.clone(),
            client.clone(),
            Arc::clone(&stats),
        )));
    }

    // Wait for all workers to finish
    for h in handles {
        h.await?;
    }

    // Give reporter one last tick
    sleep(Duration::from_millis(1200)).await;
    reporter_handle.abort();

    // Final summary
    let elapsed = started.elapsed();
    let sent = stats.sent.load(Ordering::Relaxed);
    let errors = stats.errors.load(Ordering::Relaxed);
    let lat_c = stats.lat_count.load(Ordering::Relaxed);
    let lat_s = stats.lat_sum.load(Ordering::Relaxed);
    let avg_lat = if lat_c > 0 { lat_s / lat_c } else { 0 };
    let throughput = (sent as f64 / elapsed.as_secs_f64()) as u64;

    println!();
    println!("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
    println!("  DONE");
    println!("  elapsed     : {:.2}s", elapsed.as_secs_f64());
    println!("  sent        : {sent}");
    println!("  errors      : {errors}");
    println!("  throughput  : {throughput} events/sec");
    println!("  avg lat     : {avg_lat}ms per batch request");
    println!("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");

    Ok(())
}
