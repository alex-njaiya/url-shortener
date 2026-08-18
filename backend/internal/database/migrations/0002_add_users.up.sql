CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    first_name  VARCHAR(100),
    last_name   VARCHAR(100),
    email       VARCHAR(100) UNIQUE NOT NULL,
    password    TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)
