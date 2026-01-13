-- Catalog service: shows, sessions and seat areas

-- 演出表
CREATE TABLE IF NOT EXISTS catalog.shows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    artist VARCHAR(255) NOT NULL,
    description TEXT,
    venue VARCHAR(500) NOT NULL,
    city VARCHAR(100) NOT NULL,
    poster_url TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.shows IS '演出表';
COMMENT ON COLUMN catalog.shows.id IS '演出唯一标识';
COMMENT ON COLUMN catalog.shows.title IS '演出标题';
COMMENT ON COLUMN catalog.shows.artist IS '演出艺术家/演员';
COMMENT ON COLUMN catalog.shows.description IS '演出描述';
COMMENT ON COLUMN catalog.shows.venue IS '演出场馆';
COMMENT ON COLUMN catalog.shows.city IS '演出城市';
COMMENT ON COLUMN catalog.shows.poster_url IS '海报URL';
COMMENT ON COLUMN catalog.shows.status IS '状态：draft-草稿, published-已发布, cancelled-已取消';
COMMENT ON COLUMN catalog.shows.created_at IS '创建时间';
COMMENT ON COLUMN catalog.shows.updated_at IS '更新时间';

CREATE INDEX idx_shows_city ON catalog.shows(city);
CREATE INDEX idx_shows_artist ON catalog.shows(artist);
CREATE INDEX idx_shows_status ON catalog.shows(status);

-- 演出场次表
CREATE TABLE IF NOT EXISTS catalog.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES catalog.shows(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.sessions IS '演出场次表';
COMMENT ON COLUMN catalog.sessions.id IS '场次唯一标识';
COMMENT ON COLUMN catalog.sessions.show_id IS '关联的演出ID';
COMMENT ON COLUMN catalog.sessions.start_time IS '开始时间';
COMMENT ON COLUMN catalog.sessions.end_time IS '结束时间';
COMMENT ON COLUMN catalog.sessions.status IS '状态：scheduled-已排期, completed-已完成, cancelled-已取消';
COMMENT ON COLUMN catalog.sessions.created_at IS '创建时间';
COMMENT ON COLUMN catalog.sessions.updated_at IS '更新时间';

CREATE INDEX idx_sessions_show_id ON catalog.sessions(show_id);
CREATE INDEX idx_sessions_start_time ON catalog.sessions(start_time);
CREATE INDEX idx_sessions_status ON catalog.sessions(status);

-- 座位区域表
CREATE TABLE IF NOT EXISTS catalog.seat_areas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES catalog.sessions(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price_cents BIGINT NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_available_seats CHECK (available_seats >= 0 AND available_seats <= total_seats)
);

COMMENT ON TABLE catalog.seat_areas IS '座位区域表';
COMMENT ON COLUMN catalog.seat_areas.id IS '座位区域唯一标识';
COMMENT ON COLUMN catalog.seat_areas.session_id IS '关联的场次ID';
COMMENT ON COLUMN catalog.seat_areas.name IS '座位区域名称（如VIP区、普通区）';
COMMENT ON COLUMN catalog.seat_areas.price_cents IS '价格（单位：分）';
COMMENT ON COLUMN catalog.seat_areas.total_seats IS '总座位数';
COMMENT ON COLUMN catalog.seat_areas.available_seats IS '可用座位数';
COMMENT ON COLUMN catalog.seat_areas.version IS '版本号（用于乐观锁）';
COMMENT ON COLUMN catalog.seat_areas.created_at IS '创建时间';
COMMENT ON COLUMN catalog.seat_areas.updated_at IS '更新时间';

CREATE INDEX idx_seat_areas_session_id ON catalog.seat_areas(session_id);
