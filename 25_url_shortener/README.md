# URL Shortener

A fast URL shortening service built with **Go**, **Echo**, **PostgreSQL**, and **Redis**.

## Tech Stack

- **Go** — Core language
- **Echo** — HTTP framework
- **PostgreSQL** — Persistent storage
- **Redis** — Caching layer for fast redirects

## How It Works

1. Client sends a `POST` request with a long URL
2. Service generates a unique short key using **SHA-256 hashing**
3. Short URL is stored in PostgreSQL
4. On redirect, Redis is checked first (cache hit), falling back to PostgreSQL (cache miss)
5. Client is redirected to the original URL via `307 Temporary Redirect`

## API

### Shorten a URL

```http
POST /api/urls/shorten
```

**Request**

```json
{
  "long_url": "https://portfolio.vipinkumawat.xyz",
  "user_id": "3a3dfc36-8a62-4ad5-bcdd-e75affb1002c"
}
```

**Response**

```json
{
  "short_url": "p8cLE3",
  "long_url": "https://portfolio.vipinkumawat.xyz"
}
```

---

### Redirect to Long URL

```http
GET /api/urls/:short_url
```

Redirects to the original long URL with `307 Temporary Redirect`.

**Example**

```bash
curl -L http://localhost:8888/api/urls/p8cLE3
```

---

## Caching Strategy

```
Request → Redis (cache hit?)
├── YES → Redirect immediately
└── NO → Query PostgreSQL → Cache in Redis → Redirect
```

TTL is configurable via `constants.REDIS_TTL_SECONDS`.

## Short Key Generation

- Algorithm: **SHA-256** + nanosecond timestamp salt
- Charset: `a-z`, `A-Z`, `0-9`, `-` (63 characters)
- Configurable length via `constants.SHORT_URL_LENGTH`
- Collision-safe: loops and regenerates on the rare chance of a duplicate

## Running Locally

```bash
# Start dependencies
docker-compose up -d
```

### Migrations

```bash
migrate create -ext sql -dir migrations -seq url

migrate -path migrations -database "postgresql://myuser:mypassword@localhost:5432/url_shortener?sslmode=disable" up
```

### Run the server

```go
go run main.go
```
