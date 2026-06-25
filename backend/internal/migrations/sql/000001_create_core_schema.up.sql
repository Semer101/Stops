CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(120) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(32) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    home_location GEOMETRY(Point, 4326),
    work_location GEOMETRY(Point, 4326),
    contribution_score INT NOT NULL DEFAULT 0,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_contact_required CHECK (email IS NOT NULL OR phone IS NOT NULL),
    CONSTRAINT users_role_check CHECK (role IN ('user', 'admin'))
);

CREATE TABLE IF NOT EXISTS transit_stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(160) NOT NULL,
    type VARCHAR(20) NOT NULL,
    location GEOMETRY(Point, 4326) NOT NULL,
    area VARCHAR(120),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    votes INT NOT NULL DEFAULT 0,
    availability_status VARCHAR(30) NOT NULL DEFAULT 'Unknown',
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT transit_stops_type_check CHECK (type IN ('bus_stop', 'taxi_stand')),
    CONSTRAINT transit_stops_availability_check CHECK (availability_status IN ('Available', 'Busy', 'Unknown'))
);

CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_name VARCHAR(160) NOT NULL,
    start_point VARCHAR(160) NOT NULL,
    destination_point VARCHAR(160) NOT NULL,
    fare DECIMAL(10, 2),
    city VARCHAR(120) NOT NULL DEFAULT 'Addis Ababa',
    estimated_time INT,
    status VARCHAR(20) NOT NULL DEFAULT 'Active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT routes_status_check CHECK (status IN ('Active', 'Inactive'))
);

CREATE TABLE IF NOT EXISTS fares (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    transport_type VARCHAR(20) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fares_transport_type_check CHECK (transport_type IN ('bus', 'taxi'))
);

CREATE TABLE IF NOT EXISTS trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    route_id UUID REFERENCES routes(id) ON DELETE SET NULL,
    origin VARCHAR(160) NOT NULL,
    destination VARCHAR(160) NOT NULL,
    duration INT,
    fare DECIMAL(10, 2),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS traffic_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    traffic_level VARCHAR(20) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT traffic_data_level_check CHECK (traffic_level IN ('Low', 'Medium', 'High'))
);

CREATE TABLE IF NOT EXISTS predictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    route_id UUID REFERENCES routes(id) ON DELETE SET NULL,
    prediction_type VARCHAR(30) NOT NULL,
    predicted_value VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT predictions_type_check CHECK (prediction_type IN ('ETA', 'Taxi', 'Congestion'))
);

CREATE TABLE IF NOT EXISTS crowdsource_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stop_id UUID NOT NULL REFERENCES transit_stops(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    net_votes INT NOT NULL DEFAULT 0,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT crowdsource_reports_status_check CHECK (status IN ('pending', 'verified', 'rejected'))
);

CREATE TABLE IF NOT EXISTS crowdsource_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES crowdsource_reports(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote_value INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT crowdsource_votes_value_check CHECK (vote_value IN (-1, 1)),
    CONSTRAINT crowdsource_votes_unique_user_report UNIQUE (report_id, user_id)
);

CREATE TABLE IF NOT EXISTS admin_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(120) NOT NULL,
    entity_type VARCHAR(80) NOT NULL,
    entity_id UUID,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS users_home_location_idx ON users USING GIST (home_location);
CREATE INDEX IF NOT EXISTS users_work_location_idx ON users USING GIST (work_location);
CREATE INDEX IF NOT EXISTS transit_stops_location_idx ON transit_stops USING GIST (location);
CREATE INDEX IF NOT EXISTS transit_stops_type_verified_idx ON transit_stops (type, verified);
CREATE INDEX IF NOT EXISTS routes_city_status_idx ON routes (city, status);
CREATE INDEX IF NOT EXISTS trips_user_timestamp_idx ON trips (user_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS traffic_data_route_recorded_idx ON traffic_data (route_id, recorded_at DESC);
CREATE INDEX IF NOT EXISTS predictions_user_created_idx ON predictions (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS crowdsource_reports_status_idx ON crowdsource_reports (status);
CREATE INDEX IF NOT EXISTS admin_audit_logs_admin_created_idx ON admin_audit_logs (admin_user_id, created_at DESC);
