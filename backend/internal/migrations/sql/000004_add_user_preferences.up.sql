-- Add user preferences and saved places
ALTER TABLE users ADD COLUMN IF NOT EXISTS preferred_transport_type VARCHAR(20) DEFAULT 'bus';
ALTER TABLE users ADD CONSTRAINT users_preferred_transport_check CHECK (preferred_transport_type IN ('bus', 'taxi'));

-- Create saved places table
CREATE TABLE IF NOT EXISTS saved_places (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(120) NOT NULL,
    location GEOMETRY(Point, 4326) NOT NULL,
    place_type VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT saved_places_type_check CHECK (place_type IN ('home', 'work', 'custom'))
);

CREATE INDEX IF NOT EXISTS saved_places_user_idx ON saved_places (user_id, place_type);
