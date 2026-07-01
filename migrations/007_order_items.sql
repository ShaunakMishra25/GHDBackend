CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    product_name VARCHAR(200) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL CHECK (quantity > 0),
    total_price DECIMAL(10,2) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
