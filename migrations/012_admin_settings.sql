CREATE TABLE IF NOT EXISTS admin_settings (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO admin_settings (key, value) VALUES
    ('delivery_charge_flat', '{"amount": 50, "max_order_amount": 1000}'::jsonb),
    ('delivery_charge_per_item', '{"amount": 10}'::jsonb),
    ('service_hours', '{"open": "07:00", "close": "21:00"}'::jsonb)
ON CONFLICT (key) DO NOTHING;
