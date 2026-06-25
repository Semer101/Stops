# CHANGELOG_AI.md

# AI Development Changelog

## Session 1 - 2026-06-25

### Summary

Initialized repository governance and planning documentation for STOPS based on the attached SRS, SDS, AGENTS rulebook, and first-task instructions.

No application code was implemented.

### Files Changed

- Created `AGENTS.md`
- Created `PROJECT_CONTEXT.md`
- Created `TASKS.md`
- Created `CHANGELOG_AI.md`
- Created `DECISIONS.md`
- Created `SESSION_HANDOFF.md`
- Created `README.md`

### New Features

- None.

### Bug Fixes

- None.

### Notes

- Recorded a blocking stack conflict: SDS specifies Next.js + Go/Gin, while the attached AGENTS rulebook specifies Vite React + Express/Prisma.
- Implementation is paused until developer approval and stack clarification.

## Session 2 - 2026-06-25

### Summary

Confirmed the SDS technology stack and initialized the application scaffold for a Next.js 15 frontend and Go/Gin backend.

### Files Changed

- Updated `AGENTS.md`
- Updated `PROJECT_CONTEXT.md`
- Updated `TASKS.md`
- Updated `CHANGELOG_AI.md`
- Updated `DECISIONS.md`
- Updated `SESSION_HANDOFF.md`
- Updated `README.md`
- Created frontend scaffold under `frontend/`
- Created backend scaffold under `backend/`

### New Features

- Added frontend app shell structure.
- Added backend health endpoint structure.

### Bug Fixes

- None.

## Session 3 - 2026-06-25

### Summary

Continued Milestone 2 by implementing the backend authentication foundation, core database migration, and frontend login/register screen.

### Files Changed

- Updated backend config, router, and API entry point
- Added PostgreSQL connection setup
- Added user model, repository, auth service, password service, token service, auth controller, and auth middleware
- Added core SQL migration for users, transit stops, routes, fares, trips, traffic data, predictions, crowdsourcing, votes, and admin audit logs
- Added auth service and token service tests
- Added frontend auth types, API client, auth service, and `/auth` page
- Updated project documentation and task checklist

### New Features

- `POST /api/auth/register`
- `POST /api/auth/login`
- JWT token issuing and parsing
- bcrypt password hashing
- Role middleware for protected admin routes
- Frontend login/register form

### Bug Fixes

- None.

## Session 4 - 2026-06-25

### Summary

Added an embedded Go migration runner, initial Addis Ababa seed data, and a protected profile endpoint that exercises JWT middleware.

### Files Changed

- Added `backend/cmd/migrate`
- Added `backend/internal/migrations`
- Added `backend/internal/controllers/profile_controller.go`
- Added `backend/internal/services/profile_service.go`
- Added `backend/internal/services/profile_service_test.go`
- Updated backend routes
- Updated frontend typecheck configuration
- Updated project documentation and task checklist

### New Features

- `go run ./cmd/migrate` applies embedded SQL migrations.
- Initial verified stops, taxi stands, routes, and fares are seeded.
- `GET /api/profile/me` returns the authenticated user's public profile.

### Bug Fixes

- Made frontend `pnpm run typecheck` source-only so it does not depend on generated `.next/types`.

### Known Issues

- Local migration execution could not be verified because Docker Desktop was not running or the Docker Linux engine pipe was unavailable.

## Session 5 - 2026-06-25

### Summary

Implemented the backend nearby transit stop discovery endpoint and added typed frontend service support.

### Files Changed

- Added transit stop model, repository, service, controller, and service tests
- Updated backend routes
- Added frontend transit stop types and service
- Extended frontend API client with GET support
- Updated project documentation and task checklist

### New Features

- `GET /api/stops/nearby` accepts `lat`, `lng`, optional `radius`, and optional `limit`.
- Nearby stops are sorted by PostGIS distance and include walking time estimates.
- Empty-region response uses the exact SRS prompt.

### Bug Fixes

- None.

## Session 6 - 2026-06-25

### Summary

Audited MVP completion against the SRS/SDS and prepared deployment configuration files.

### Files Changed

- Added `MVP_AUDIT.md`
- Added `DEPLOYMENT.md`
- Added `backend/Dockerfile`
- Added `render.yaml`
- Added `frontend/vercel.json`
- Updated backend CORS configuration
- Updated project context and task checklist

### New Features

- Backend Docker build configuration.
- Render blueprint for backend deployment.
- Vercel config for frontend deployment.
- Deployment environment variable documentation.

### Known Issues

- The SRS/SDS MVP is not yet complete.
- Actual deployment was not performed because provider authentication, git remote, production database, and production environment variables are not configured.
- Backend Docker image build timed out before producing a local image.

## Session 7 - 2026-06-25

### Summary

Closed the nearby-search manual location and walking-directions MVP gaps.

### Files Changed

- Added `frontend/types/directions.ts`
- Added `frontend/services/osrmService.ts`
- Updated `frontend/components/Map.tsx`
- Updated `frontend/app/page.tsx`
- Updated `frontend/.env.example`
- Updated MVP audit and project tracking docs

### New Features

- Manual coordinate input for nearby stop searches.
- OSRM-backed walking directions from current/manual location to a selected stop.
- Direction summary and step list in the sidebar.
- Straight-line fallback guidance when OSRM is unavailable.

### Verification

- `go test ./...`
- `pnpm run typecheck`
- `pnpm run lint`
- `pnpm run build`

## Session 6 - 2026-06-25

### Summary

Successfully connected to local database by mapping container to port `5433` (resolving conflicts with host PostgreSQL service) and completed the core features of the STOPS transit planning application. Implemented crowdsourced stop reporting/voting, route searches, fare calculations, travel predictions, and Gemini chatbot integration. Developed the client-side Leaflet Map component with divIcon custom pin styles and a slide-out assistant chatbot drawer, integrated within a premium Addis Ababa styled homepage. Verified that all backend service tests, frontend typechecks, ESLint syntax checks, and Next.js production builds compile and pass.

### Files Changed

- Modified `docker-compose.yml` (mapped container to host port `5433`)
- Created `backend/.env` and `frontend/.env` configurations
- Created `backend/internal/models/route.go`
- Created `backend/internal/repositories/route_repository.go`
- Created `backend/internal/repositories/crowdsourcing_repository.go`
- Created `backend/internal/services/route_service.go`
- Created `backend/internal/services/crowdsourcing_service.go`
- Created `backend/internal/services/prediction_service.go`
- Created `backend/internal/services/chat_service.go`
- Created `backend/internal/controllers/route_controller.go`
- Created `backend/internal/controllers/crowdsource_controller.go`
- Created `backend/internal/controllers/prediction_controller.go`
- Created `backend/internal/controllers/chat_controller.go`
- Modified `backend/internal/routes/router.go`
- Modified `frontend/app/layout.tsx`
- Created `frontend/components/Map.tsx`
- Created `frontend/components/Chatbot.tsx`
- Modified `frontend/app/page.tsx`
- Updated project context, task checklist, and session handoff

### New Features

- Applied full SQL database schema migrations and seeded Addis Ababa transit stops, routes, and fares database tables.
- `GET /api/routes` & `GET /api/fare` for route queries and fare estimation.
- `GET /api/predictions` yielding heuristic predictions for arrival ETAs, taxi counts, and traffic congestion.
- `POST /api/chat` sending formatted travel histories and route details to Google Gemini AI transit assistant.
- `POST /api/stops/report` & `POST /api/stops/vote` for authenticated stop reporting (limited to 5 pending per 24h) and community voting with auto-verification.
- Map-centered presentation layer showing responsive UI drawer controls, custom styled SVG marker pins, filters, and dynamic map center reloading.

### Bug Fixes

- Resolved container port binding conflict with host PostgreSQL instance by routing to 5433.
- Resolved Go compile errors by cleaning up package-wide duplicate `writeError` declarations.
- Resolved ESLint explicit-any, unused-vars, and hooks dependencies warnings in frontend modules.
