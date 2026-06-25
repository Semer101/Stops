# Migrations

Database migrations are embedded in the backend under `backend/internal/migrations/sql`.

PostgreSQL with PostGIS is required by the SRS and SDS.

Run local migrations with:

```bash
cd backend
go run ./cmd/migrate
```
