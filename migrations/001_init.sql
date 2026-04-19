CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    customer TEXT NOT NULL,
    email TEXT NOT NULL,
    amount NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
    status TEXT NOT NULL CHECK (status IN ('new', 'paid', 'shipped', 'cancelled')),
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
