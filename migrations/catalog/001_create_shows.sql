-- catalog/001_create_shows.sql

-- Shows table
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

CREATE INDEX idx_shows_city ON catalog.shows(city);
CREATE INDEX idx_shows_artist ON catalog.shows(artist);
CREATE INDEX idx_shows_status ON catalog.shows(status);

-- Sessions table (performance times)
CREATE TABLE IF NOT EXISTS catalog.sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    show_id UUID NOT NULL REFERENCES catalog.shows(id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_show_id ON catalog.sessions(show_id);
CREATE INDEX idx_sessions_start_time ON catalog.sessions(start_time);
CREATE INDEX idx_sessions_status ON catalog.sessions(status);

-- Seat areas table
CREATE TABLE IF NOT EXISTS catalog.seat_areas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES catalog.sessions(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    price_cents BIGINT NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_available_seats CHECK (available_seats >= 0 AND available_seats <= total_seats)
);

CREATE INDEX idx_seat_areas_session_id ON catalog.seat_areas(session_id);

-- Version for optimistic locking on inventory
ALTER TABLE catalog.seat_areas ADD COLUMN version INTEGER DEFAULT 1;
