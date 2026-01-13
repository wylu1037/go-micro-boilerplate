-- Identity service: users and related tables

-- 用户表
CREATE TABLE IF NOT EXISTS identity.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    phone VARCHAR(50),
    avatar_url TEXT,
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE identity.users IS '用户表';
COMMENT ON COLUMN identity.users.id IS '用户唯一标识';
COMMENT ON COLUMN identity.users.email IS '用户邮箱地址';
COMMENT ON COLUMN identity.users.password_hash IS '密码哈希值';
COMMENT ON COLUMN identity.users.name IS '用户姓名';
COMMENT ON COLUMN identity.users.phone IS '手机号码';
COMMENT ON COLUMN identity.users.avatar_url IS '头像URL';
COMMENT ON COLUMN identity.users.email_verified IS '邮箱是否已验证';
COMMENT ON COLUMN identity.users.created_at IS '创建时间';
COMMENT ON COLUMN identity.users.updated_at IS '更新时间';

CREATE INDEX idx_users_email ON identity.users(email);

-- 刷新令牌表
CREATE TABLE IF NOT EXISTS identity.refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE identity.refresh_tokens IS '刷新令牌表';
COMMENT ON COLUMN identity.refresh_tokens.id IS '令牌唯一标识';
COMMENT ON COLUMN identity.refresh_tokens.user_id IS '关联的用户ID';
COMMENT ON COLUMN identity.refresh_tokens.token_hash IS '令牌哈希值';
COMMENT ON COLUMN identity.refresh_tokens.expires_at IS '过期时间';
COMMENT ON COLUMN identity.refresh_tokens.created_at IS '创建时间';

CREATE INDEX idx_refresh_tokens_user_id ON identity.refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON identity.refresh_tokens(expires_at);

-- 密码重置令牌表
CREATE TABLE IF NOT EXISTS identity.password_reset_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES identity.users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE identity.password_reset_tokens IS '密码重置令牌表';
COMMENT ON COLUMN identity.password_reset_tokens.id IS '令牌唯一标识';
COMMENT ON COLUMN identity.password_reset_tokens.user_id IS '关联的用户ID';
COMMENT ON COLUMN identity.password_reset_tokens.token_hash IS '令牌哈希值';
COMMENT ON COLUMN identity.password_reset_tokens.expires_at IS '过期时间';
COMMENT ON COLUMN identity.password_reset_tokens.used IS '是否已使用';
COMMENT ON COLUMN identity.password_reset_tokens.created_at IS '创建时间';

CREATE INDEX idx_password_reset_tokens_user_id ON identity.password_reset_tokens(user_id);
