-- Catalog service: shows, venues, sessions, seat_areas

-- 场馆表
CREATE TABLE IF NOT EXISTS catalog.venues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    address TEXT,
    capacity INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.venues IS '场馆表';
COMMENT ON COLUMN catalog.venues.id IS '场馆唯一标识';
COMMENT ON COLUMN catalog.venues.name IS '场馆名称';
COMMENT ON COLUMN catalog.venues.city IS '所在城市';
COMMENT ON COLUMN catalog.venues.address IS '详细地址';
COMMENT ON COLUMN catalog.venues.capacity IS '总容量';
COMMENT ON COLUMN catalog.venues.created_at IS '创建时间';

CREATE INDEX idx_venues_city ON catalog.venues(city);

-- 演出表
CREATE TABLE IF NOT EXISTS catalog.shows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    artist VARCHAR(255),
    category VARCHAR(50) NOT NULL,
    poster_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.shows IS '演出表';
COMMENT ON COLUMN catalog.shows.id IS '演出唯一标识';
COMMENT ON COLUMN catalog.shows.title IS '演出名称';
COMMENT ON COLUMN catalog.shows.description IS '演出描述';
COMMENT ON COLUMN catalog.shows.artist IS '艺人/团体名称';
COMMENT ON COLUMN catalog.shows.category IS '演出类型 (concert/musical/sports)';
COMMENT ON COLUMN catalog.shows.poster_url IS '海报图片URL';
COMMENT ON COLUMN catalog.shows.status IS '状态 (draft/published/cancelled)';
COMMENT ON COLUMN catalog.shows.created_at IS '创建时间';
COMMENT ON COLUMN catalog.shows.updated_at IS '更新时间';

CREATE INDEX idx_shows_category ON catalog.shows(category);
CREATE INDEX idx_shows_status ON catalog.shows(status);

-- 场次表
CREATE TABLE IF NOT EXISTS catalog.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES catalog.shows(id) ON DELETE CASCADE,
    venue_id UUID NOT NULL REFERENCES catalog.venues(id) ON DELETE RESTRICT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ,
    sale_start_time TIMESTAMPTZ,
    sale_end_time TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.sessions IS '场次表';
COMMENT ON COLUMN catalog.sessions.id IS '场次唯一标识';
COMMENT ON COLUMN catalog.sessions.show_id IS '所属演出';
COMMENT ON COLUMN catalog.sessions.venue_id IS '所在场馆';
COMMENT ON COLUMN catalog.sessions.start_time IS '开始时间';
COMMENT ON COLUMN catalog.sessions.end_time IS '结束时间';
COMMENT ON COLUMN catalog.sessions.sale_start_time IS '开售时间';
COMMENT ON COLUMN catalog.sessions.sale_end_time IS '停售时间';
COMMENT ON COLUMN catalog.sessions.status IS '状态 (scheduled/on_sale/sold_out/cancelled)';
COMMENT ON COLUMN catalog.sessions.created_at IS '创建时间';

CREATE INDEX idx_sessions_show_id ON catalog.sessions(show_id);
CREATE INDEX idx_sessions_venue_id ON catalog.sessions(venue_id);
CREATE INDEX idx_sessions_start_time ON catalog.sessions(start_time);
CREATE INDEX idx_sessions_status ON catalog.sessions(status);

-- 座位区域表
CREATE TABLE IF NOT EXISTS catalog.seat_areas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES catalog.sessions(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    total_seats INT NOT NULL,
    available_seats INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE catalog.seat_areas IS '座位区域表';
COMMENT ON COLUMN catalog.seat_areas.id IS '区域唯一标识';
COMMENT ON COLUMN catalog.seat_areas.session_id IS '所属场次';
COMMENT ON COLUMN catalog.seat_areas.name IS '区域名称 (VIP A区、内场B区)';
COMMENT ON COLUMN catalog.seat_areas.price IS '票价';
COMMENT ON COLUMN catalog.seat_areas.total_seats IS '总座位数';
COMMENT ON COLUMN catalog.seat_areas.available_seats IS '可用座位数';
COMMENT ON COLUMN catalog.seat_areas.created_at IS '创建时间';

CREATE INDEX idx_seat_areas_session_id ON catalog.seat_areas(session_id);

-- 库存约束: 可用座位数不能超过总座位数
ALTER TABLE catalog.seat_areas
    ADD CONSTRAINT check_available_seats
    CHECK (available_seats >= 0 AND available_seats <= total_seats);
