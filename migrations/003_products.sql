CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name_hi VARCHAR(200) NOT NULL,
    name_en VARCHAR(200) NOT NULL,
    description_hi TEXT NOT NULL DEFAULT '',
    description_en TEXT NOT NULL DEFAULT '',
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    unit VARCHAR(50) NOT NULL DEFAULT 'piece',
    image_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT true,
    stock_qty DECIMAL(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_name_en ON products(name_en);
CREATE INDEX IF NOT EXISTS idx_products_name_hi ON products(name_hi);
CREATE INDEX IF NOT EXISTS idx_products_created ON products(created_at);
