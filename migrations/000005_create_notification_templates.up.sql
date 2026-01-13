-- Notification service: templates and send logs

-- 消息模板表
CREATE TABLE IF NOT EXISTS notification.templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    subject VARCHAR(500),
    content TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE notification.templates IS '消息模板表';
COMMENT ON COLUMN notification.templates.id IS '模板唯一标识';
COMMENT ON COLUMN notification.templates.code IS '模板代码（唯一标识符）';
COMMENT ON COLUMN notification.templates.name IS '模板名称';
COMMENT ON COLUMN notification.templates.type IS '消息类型：email-邮件, sms-短信, push-推送';
COMMENT ON COLUMN notification.templates.subject IS '消息主题（仅用于邮件）';
COMMENT ON COLUMN notification.templates.content IS '消息内容模板（支持变量替换，如{{user_name}}）';
COMMENT ON COLUMN notification.templates.active IS '是否激活';
COMMENT ON COLUMN notification.templates.created_at IS '创建时间';
COMMENT ON COLUMN notification.templates.updated_at IS '更新时间';

CREATE INDEX idx_templates_code ON notification.templates(code);
CREATE INDEX idx_templates_type ON notification.templates(type);

-- 发送日志表
CREATE TABLE IF NOT EXISTS notification.send_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_code VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    sent_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE notification.send_logs IS '消息发送日志表';
COMMENT ON COLUMN notification.send_logs.id IS '日志唯一标识';
COMMENT ON COLUMN notification.send_logs.template_code IS '使用的模板代码';
COMMENT ON COLUMN notification.send_logs.type IS '消息类型：email-邮件, sms-短信, push-推送';
COMMENT ON COLUMN notification.send_logs.recipient IS '接收者（邮箱/手机号）';
COMMENT ON COLUMN notification.send_logs.status IS '发送状态：pending-待发送, sent-已发送, failed-失败';
COMMENT ON COLUMN notification.send_logs.error_message IS '错误信息（发送失败时记录）';
COMMENT ON COLUMN notification.send_logs.sent_at IS '发送时间';

CREATE INDEX idx_send_logs_recipient ON notification.send_logs(recipient);
CREATE INDEX idx_send_logs_template_code ON notification.send_logs(template_code);
CREATE INDEX idx_send_logs_sent_at ON notification.send_logs(sent_at);

-- Insert default templates
INSERT INTO notification.templates (code, name, type, subject, content) VALUES
('ORDER_CREATED', '订单创建通知', 'email', '您的订单已创建 - {{show_title}}', '尊敬的{{user_name}}，您的订单{{order_id}}已创建成功，请在{{expire_time}}前完成支付。'),
('PAYMENT_SUCCESS', '支付成功通知', 'email', '支付成功 - {{show_title}}', '尊敬的{{user_name}}，您已成功购买{{show_title}}门票{{quantity}}张，订单号：{{order_id}}。'),
('ORDER_CANCELLED', '订单取消通知', 'email', '订单已取消 - {{show_title}}', '尊敬的{{user_name}}，您的订单{{order_id}}已取消。'),
('ORDER_CREATED_SMS', '订单创建短信', 'sms', NULL, '【票务系统】您的订单{{order_id}}已创建，请在{{expire_minutes}}分钟内完成支付。'),
('PAYMENT_SUCCESS_SMS', '支付成功短信', 'sms', NULL, '【票务系统】支付成功！订单号{{order_id}}，{{show_title}}门票{{quantity}}张。')
ON CONFLICT (code) DO NOTHING;
