-- notification/001_create_templates.sql

-- Message templates
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

CREATE INDEX idx_templates_code ON notification.templates(code);
CREATE INDEX idx_templates_type ON notification.templates(type);

-- Send logs for auditing
CREATE TABLE IF NOT EXISTS notification.send_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_code VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    error_message TEXT,
    sent_at TIMESTAMPTZ DEFAULT NOW()
);

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
