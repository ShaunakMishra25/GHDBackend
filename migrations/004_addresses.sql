CREATE TABLE IF NOT EXISTS addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label VARCHAR(50) NOT NULL DEFAULT 'Home',
    full_address TEXT NOT NULL,
    landmark VARCHAR(200) NOT NULL DEFAULT '',
    latitude DECIMAL(10,7) NOT NULL DEFAULT 0,
    longitude DECIMAL(10,7) NOT NULL DEFAULT 0,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_addresses_user ON addresses(user_id);
