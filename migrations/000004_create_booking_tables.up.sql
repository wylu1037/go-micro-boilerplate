-- Booking service: orders, tickets, payments

-- 订单表
CREATE TABLE IF NOT EXISTS booking.orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_no VARCHAR(32) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    session_id UUID NOT NULL,
    seat_area_id UUID NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending_payment',
    expires_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE booking.orders IS '订单表';
COMMENT ON COLUMN booking.orders.id IS '订单唯一标识';
COMMENT ON COLUMN booking.orders.order_no IS '订单号 (业务编号)';
COMMENT ON COLUMN booking.orders.user_id IS '下单用户ID';
COMMENT ON COLUMN booking.orders.session_id IS '场次ID';
COMMENT ON COLUMN booking.orders.seat_area_id IS '座位区域ID';
COMMENT ON COLUMN booking.orders.quantity IS '购买数量';
COMMENT ON COLUMN booking.orders.unit_price IS '单价';
COMMENT ON COLUMN booking.orders.total_amount IS '总金额';
COMMENT ON COLUMN booking.orders.status IS '状态 (pending_payment/paid/cancelled/refunded/completed)';
COMMENT ON COLUMN booking.orders.expires_at IS '支付截止时间';
COMMENT ON COLUMN booking.orders.paid_at IS '支付时间';
COMMENT ON COLUMN booking.orders.cancelled_at IS '取消时间';
COMMENT ON COLUMN booking.orders.created_at IS '创建时间';
COMMENT ON COLUMN booking.orders.updated_at IS '更新时间';

CREATE INDEX idx_orders_user_id ON booking.orders(user_id);
CREATE INDEX idx_orders_session_id ON booking.orders(session_id);
CREATE INDEX idx_orders_status ON booking.orders(status);
CREATE INDEX idx_orders_created_at ON booking.orders(created_at);
CREATE INDEX idx_orders_expires_at ON booking.orders(expires_at) WHERE status = 'pending_payment';

-- 电子票表
CREATE TABLE IF NOT EXISTS booking.tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES booking.orders(id) ON DELETE CASCADE,
    ticket_no VARCHAR(32) NOT NULL UNIQUE,
    qr_code TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'valid',
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE booking.tickets IS '电子票表';
COMMENT ON COLUMN booking.tickets.id IS '票唯一标识';
COMMENT ON COLUMN booking.tickets.order_id IS '所属订单';
COMMENT ON COLUMN booking.tickets.ticket_no IS '票号';
COMMENT ON COLUMN booking.tickets.qr_code IS '二维码内容';
COMMENT ON COLUMN booking.tickets.status IS '状态 (valid/used/refunded)';
COMMENT ON COLUMN booking.tickets.used_at IS '使用时间';
COMMENT ON COLUMN booking.tickets.created_at IS '创建时间';

CREATE INDEX idx_tickets_order_id ON booking.tickets(order_id);
CREATE INDEX idx_tickets_status ON booking.tickets(status);

-- 支付记录表
CREATE TABLE IF NOT EXISTS booking.payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES booking.orders(id) ON DELETE CASCADE,
    payment_method VARCHAR(20) NOT NULL,
    transaction_id VARCHAR(100) UNIQUE,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE booking.payments IS '支付记录表';
COMMENT ON COLUMN booking.payments.id IS '支付记录唯一标识';
COMMENT ON COLUMN booking.payments.order_id IS '关联订单';
COMMENT ON COLUMN booking.payments.payment_method IS '支付方式 (alipay/wechat/card)';
COMMENT ON COLUMN booking.payments.transaction_id IS '第三方交易号';
COMMENT ON COLUMN booking.payments.amount IS '支付金额';
COMMENT ON COLUMN booking.payments.status IS '状态 (pending/success/failed/refunded)';
COMMENT ON COLUMN booking.payments.paid_at IS '支付成功时间';
COMMENT ON COLUMN booking.payments.created_at IS '创建时间';

CREATE INDEX idx_payments_order_id ON booking.payments(order_id);
CREATE INDEX idx_payments_status ON booking.payments(status);
