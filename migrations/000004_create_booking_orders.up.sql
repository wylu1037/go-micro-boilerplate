-- Booking service: orders, order items and payments

-- 订单表
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

COMMENT ON TABLE booking.orders IS '订单表';
COMMENT ON COLUMN booking.orders.id IS '订单唯一标识';
COMMENT ON COLUMN booking.orders.user_id IS '用户ID';
COMMENT ON COLUMN booking.orders.session_id IS '演出场次ID';
COMMENT ON COLUMN booking.orders.show_title IS '演出标题（冗余字段，避免跨服务查询）';
COMMENT ON COLUMN booking.orders.status IS '订单状态：pending-待支付, paid-已支付, cancelled-已取消, expired-已过期';
COMMENT ON COLUMN booking.orders.total_amount_cents IS '订单总金额（单位：分）';
COMMENT ON COLUMN booking.orders.expires_at IS '订单过期时间';
COMMENT ON COLUMN booking.orders.created_at IS '创建时间';
COMMENT ON COLUMN booking.orders.updated_at IS '更新时间';

CREATE INDEX idx_orders_user_id ON booking.orders(user_id);
CREATE INDEX idx_orders_session_id ON booking.orders(session_id);
CREATE INDEX idx_orders_status ON booking.orders(status);
CREATE INDEX idx_orders_expires_at ON booking.orders(expires_at) WHERE status = 'pending';

-- 订单明细表
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

COMMENT ON TABLE booking.order_items IS '订单明细表';
COMMENT ON COLUMN booking.order_items.id IS '订单明细唯一标识';
COMMENT ON COLUMN booking.order_items.order_id IS '关联的订单ID';
COMMENT ON COLUMN booking.order_items.seat_area_id IS '座位区域ID';
COMMENT ON COLUMN booking.order_items.seat_area_name IS '座位区域名称（冗余字段）';
COMMENT ON COLUMN booking.order_items.quantity IS '购买数量';
COMMENT ON COLUMN booking.order_items.unit_price_cents IS '单价（单位：分）';
COMMENT ON COLUMN booking.order_items.subtotal_cents IS '小计金额（单位：分）';
COMMENT ON COLUMN booking.order_items.created_at IS '创建时间';

CREATE INDEX idx_order_items_order_id ON booking.order_items(order_id);

-- 支付记录表
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

COMMENT ON TABLE booking.payments IS '支付记录表';
COMMENT ON COLUMN booking.payments.id IS '支付记录唯一标识';
COMMENT ON COLUMN booking.payments.order_id IS '关联的订单ID';
COMMENT ON COLUMN booking.payments.method IS '支付方式：alipay-支付宝, wechat-微信, credit_card-信用卡';
COMMENT ON COLUMN booking.payments.status IS '支付状态：pending-待支付, success-成功, failed-失败, refunded-已退款';
COMMENT ON COLUMN booking.payments.amount_cents IS '支付金额（单位：分）';
COMMENT ON COLUMN booking.payments.transaction_id IS '第三方支付平台交易ID';
COMMENT ON COLUMN booking.payments.raw_callback IS '支付回调原始数据';
COMMENT ON COLUMN booking.payments.paid_at IS '支付完成时间';
COMMENT ON COLUMN booking.payments.created_at IS '创建时间';
COMMENT ON COLUMN booking.payments.updated_at IS '更新时间';

CREATE INDEX idx_payments_order_id ON booking.payments(order_id);
CREATE INDEX idx_payments_transaction_id ON booking.payments(transaction_id);
CREATE INDEX idx_payments_status ON booking.payments(status);
