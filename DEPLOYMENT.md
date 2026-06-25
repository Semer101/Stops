# Deployment

## Status

The app is deploy-ready at the configuration level, but actual deployment requires authenticated provider access and production environment variables.

## Backend

Backend target: Render or another Docker-capable host.

Required environment variables:

```text
APP_ENV=production
PORT=8080
DATABASE_URL=<production PostgreSQL/PostGIS connection string>
JWT_SECRET=<strong random secret>
JWT_TTL_MINUTES=60
CORS_ALLOWED_ORIGIN=<deployed frontend origin>
GEMINI_API_KEY=<optional Gemini API key>
```

Render blueprint:

```text
render.yaml
```

The backend Docker image includes:

- `stops-api`
- `stops-migrate`

The Render blueprint runs `stops-migrate` before deploying the web process.

## Frontend

Frontend target: Vercel.

Required environment variable:

```text
NEXT_PUBLIC_API_BASE_URL=<deployed backend origin>
NEXT_PUBLIC_OSRM_BASE_URL=https://router.project-osrm.org
```

Vercel config:

```text
frontend/vercel.json
```

## Local Verification Commands

```bash
cd backend
go test ./...

cd ../frontend
pnpm run typecheck
pnpm run lint
pnpm run build
```
