-- Rollback user preferences and saved places
DROP INDEX IF EXISTS saved_places_user_idx;
DROP TABLE IF EXISTS saved_places;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_preferred_transport_check;
ALTER TABLE users DROP COLUMN IF EXISTS preferred_transport_type;
