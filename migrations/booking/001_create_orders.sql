-- booking/001_create_orders.sql

-- Orders table
CREATE TABLE IF NOT EXISTS booking.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    session_id UUID NOT NULL,
    show_title VARCHAR(500) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    total_amount_cents BIGINT NOT NULL,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON booking.orders(user_id);
CREATE INDEX idx_orders_session_id ON booking.orders(session_id);
CREATE INDEX idx_orders_status ON booking.orders(status);
CREATE INDEX idx_orders_expires_at ON booking.orders(expires_at) WHERE status = 'pending';

-- Order items table
CREATE TABLE IF NOT EXISTS booking.order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES booking.orders(id) ON DELETE CASCADE,
    seat_area_id UUID NOT NULL,
    seat_area_name VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price_cents BIGINT NOT NULL,
    subtotal_cents BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_order_items_order_id ON booking.order_items(order_id);

-- Payments table
CREATE TABLE IF NOT EXISTS booking.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES booking.orders(id) ON DELETE CASCADE,
    method VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    amount_cents BIGINT NOT NULL,
    transaction_id VARCHAR(255),
    raw_callback TEXT,
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order_id ON booking.payments(order_id);
CREATE INDEX idx_payments_transaction_id ON booking.payments(transaction_id);
CREATE INDEX idx_payments_status ON booking.payments(status);
