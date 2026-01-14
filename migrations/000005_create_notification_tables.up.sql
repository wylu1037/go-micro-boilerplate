-- Notification service: templates, logs

-- 消息模板表
CREATE TABLE IF NOT EXISTS notification.templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE notification.templates IS '消息模板表';
COMMENT ON COLUMN notification.templates.id IS '模板唯一标识';
COMMENT ON COLUMN notification.templates.name IS '模板名称';
COMMENT ON COLUMN notification.templates.type IS '类型 (email/sms)';
COMMENT ON COLUMN notification.templates.subject IS '邮件主题';
COMMENT ON COLUMN notification.templates.content IS '模板内容';
COMMENT ON COLUMN notification.templates.created_at IS '创建时间';
COMMENT ON COLUMN notification.templates.updated_at IS '更新时间';

CREATE INDEX idx_templates_type ON notification.templates(type);

-- 发送日志表
CREATE TABLE IF NOT EXISTS notification.logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID REFERENCES notification.templates(id) ON DELETE SET NULL,
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    content TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE notification.logs IS '发送日志表';
COMMENT ON COLUMN notification.logs.id IS '日志唯一标识';
COMMENT ON COLUMN notification.logs.template_id IS '使用的模板';
COMMENT ON COLUMN notification.logs.user_id IS '接收用户';
COMMENT ON COLUMN notification.logs.type IS '类型 (email/sms)';
COMMENT ON COLUMN notification.logs.recipient IS '接收地址 (邮箱/手机号)';
COMMENT ON COLUMN notification.logs.content IS '实际发送内容';
COMMENT ON COLUMN notification.logs.status IS '状态 (pending/sent/failed)';
COMMENT ON COLUMN notification.logs.sent_at IS '发送时间';
COMMENT ON COLUMN notification.logs.error_message IS '失败原因';
COMMENT ON COLUMN notification.logs.created_at IS '创建时间';

CREATE INDEX idx_logs_user_id ON notification.logs(user_id);
CREATE INDEX idx_logs_type ON notification.logs(type);
CREATE INDEX idx_logs_status ON notification.logs(status);
CREATE INDEX idx_logs_created_at ON notification.logs(created_at);

-- 插入默认模板
INSERT INTO notification.templates (name, type, subject, content) VALUES
('order_confirmation', 'email', '订单确认 - {{.OrderNo}}', '尊敬的 {{.UserName}}，您的订单 {{.OrderNo}} 已创建成功，请在 {{.ExpireTime}} 前完成支付。'),
('payment_success', 'email', '支付成功 - {{.OrderNo}}', '尊敬的 {{.UserName}}，您的订单 {{.OrderNo}} 已支付成功，电子票已发送至您的邮箱。'),
('order_cancelled', 'email', '订单取消 - {{.OrderNo}}', '尊敬的 {{.UserName}}，您的订单 {{.OrderNo}} 已取消。'),
('order_confirmation_sms', 'sms', NULL, '【票务系统】您的订单 {{.OrderNo}} 已创建，请在 {{.ExpireTime}} 前完成支付。'),
('payment_success_sms', 'sms', NULL, '【票务系统】订单 {{.OrderNo}} 支付成功，电子票已发送至邮箱。');
