CREATE TABLE urls (
    id          BIGSERIAL PRIMARY KEY,
    short_code  TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_urls_short_code ON urls (short_code);

CREATE TABLE clicks (
    id          BIGSERIAL PRIMARY KEY,
    url_id      BIGINT NOT NULL REFERENCES urls (id) ON DELETE CASCADE,
    clicked_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    referrer    TEXT,
    user_agent  TEXT
);

CREATE INDEX idx_clicks_url_id ON clicks (url_id);
CREATE INDEX idx_clicks_clicked_at ON clicks (clicked_at);