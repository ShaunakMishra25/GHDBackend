CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fcm_token TEXT NOT NULL,
    device_info JSONB NOT NULL DEFAULT '{}',
    last_active_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, fcm_token)
);

CREATE INDEX IF NOT EXISTS idx_devices_user ON devices(user_id);
