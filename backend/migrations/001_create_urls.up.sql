CREATE TABLE urls (
    id           BIGSERIAL PRIMARY KEY,
    short_code   VARCHAR(10) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    click_count  BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX idx_urls_short_code ON urls (short_code);
