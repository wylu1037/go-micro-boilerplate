-- Create schemas for each service

CREATE SCHEMA IF NOT EXISTS identity;
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS booking;
CREATE SCHEMA IF NOT EXISTS notification;

COMMENT ON SCHEMA identity IS '身份认证服务 Schema - 管理用户、令牌等';
COMMENT ON SCHEMA catalog IS '演出目录服务 Schema - 管理演出、场次、座位';
COMMENT ON SCHEMA booking IS '订单预订服务 Schema - 管理订单、支付';
COMMENT ON SCHEMA notification IS '通知服务 Schema - 管理消息模板、发送日志';

-- Grant permissions (adjust as needed for production)
GRANT ALL ON SCHEMA identity TO postgres;
GRANT ALL ON SCHEMA catalog TO postgres;
GRANT ALL ON SCHEMA booking TO postgres;
GRANT ALL ON SCHEMA notification TO postgres;
