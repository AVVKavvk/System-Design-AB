CREATE TABLE IF NOT EXISTS urls (
    short_url TEXT PRIMARY KEY,
    long_url TEXT NOT NULL,
    user_id uuid NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);