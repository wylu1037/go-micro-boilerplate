-- Create schemas for each service

CREATE SCHEMA IF NOT EXISTS identity;
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS booking;
CREATE SCHEMA IF NOT EXISTS notification;

COMMENT ON SCHEMA identity IS '身份认证服务 Schema - 管理用户、令牌等';
COMMENT ON SCHEMA catalog IS '演出目录服务 Schema - 管理演出、场次、座位';
COMMENT ON SCHEMA booking IS '订单预订服务 Schema - 管理订单、支付';
COMMENT ON SCHEMA notification IS '通知服务 Schema - 管理消息模板、发送日志';

-- Grant permissions (use PUBLIC for development compatibility)
GRANT ALL ON SCHEMA identity TO PUBLIC;
GRANT ALL ON SCHEMA catalog TO PUBLIC;
GRANT ALL ON SCHEMA booking TO PUBLIC;
GRANT ALL ON SCHEMA notification TO PUBLIC;
