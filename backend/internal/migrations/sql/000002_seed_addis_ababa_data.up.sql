INSERT INTO transit_stops (name, type, location, area, verified, votes, availability_status)
VALUES
    ('Mexico Square Bus Stop', 'bus_stop', ST_SetSRID(ST_MakePoint(38.7469, 9.0105), 4326), 'Mexico Square', TRUE, 5, 'Unknown'),
    ('Meskel Square Bus Stop', 'bus_stop', ST_SetSRID(ST_MakePoint(38.7612, 9.0108), 4326), 'Meskel Square', TRUE, 5, 'Unknown'),
    ('Bole Medhanialem Taxi Stand', 'taxi_stand', ST_SetSRID(ST_MakePoint(38.7890, 8.9984), 4326), 'Bole', TRUE, 5, 'Available'),
    ('Piazza Taxi Stand', 'taxi_stand', ST_SetSRID(ST_MakePoint(38.7578, 9.0367), 4326), 'Piazza', TRUE, 5, 'Available')
ON CONFLICT DO NOTHING;

INSERT INTO routes (route_name, start_point, destination_point, fare, city, estimated_time, status)
VALUES
    ('Bole to Piazza', 'Bole', 'Piazza', 25.00, 'Addis Ababa', 35, 'Active'),
    ('Mexico to Meskel Square', 'Mexico Square', 'Meskel Square', 15.00, 'Addis Ababa', 18, 'Active'),
    ('Piazza to Mexico', 'Piazza', 'Mexico Square', 18.00, 'Addis Ababa', 24, 'Active')
ON CONFLICT DO NOTHING;

INSERT INTO fares (route_id, transport_type, amount)
SELECT id, 'bus', fare
FROM routes
WHERE city = 'Addis Ababa' AND fare IS NOT NULL
ON CONFLICT DO NOTHING;
