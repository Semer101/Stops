# PROJECT_CONTEXT.md

# STOPS Project Context

## Overall Completion

90%

## Current Status

All MVP features have been implemented and the application passes local build checks. The codebase is ready for production deployment. Remaining work requires user-provided credentials for git remote, deployment providers (Vercel/Render), and production database setup.

## Authoritative Documents Read

- SRS: `C:\Users\PC\Downloads\SRS_Final.pdf`
- SDS: `C:\Users\PC\Downloads\SDS_Final.pdf`
- Attached AGENTS rulebook: `C:\Users\PC\.codex\attachments\efb0c844-0872-44e7-9f01-6f16b2464e16\pasted-text.txt`

## Current Architecture

Layered full-stack web architecture:
- Presentation layer: responsive React/Next.js client with map-centered Leaflet UI.
- Application/API layer: REST backend with JWT authentication built on Go and Gin.
- Service layer: business logic for route planning, fare/time estimation, heuristic predictions, and crowdsourcing.
- Repository/data layer: PostgreSQL with PostGIS for spatial index queries.
- Prediction layer: lightweight heuristic engine.
- AI orchestration layer: Google Gemini for natural-language explanations and chatbot responses.

## Current Technology Stack

- Frontend: Next.js 15, React, TypeScript, Tailwind CSS, Leaflet
- Backend: Go + Gin
- Database: PostgreSQL + PostGIS (running on Docker container port 5433)
- Prediction: Heuristic prediction engine (ETA, taxi availability, congestion)
- AI/NLP: Google Gemini API (with unconfigured key graceful fallback)
- Mapping: OpenStreetMap + Leaflet

## Completed Modules

- Repository governance documentation initialized
- Requirements summary and implementation roadmap captured
- Next.js 15 frontend scaffold and Go/Gin backend scaffold created
- Docker Compose PostGIS database service initialized (mapped to port 5433 to bypass host conflicts)
- Core database migrations applied and seed data loaded successfully
- Register and Login endpoints implemented (with bcrypt hashing and JWT tokens)
- User Profile endpoint (`GET /api/profile/me`) implemented
- Nearby Transit Stop endpoint (`GET /api/stops/nearby`) implemented
- Crowdsourced Stop Reporting and Voting endpoints implemented (`POST /api/stops/report` and `POST /api/stops/vote`) with 24h rate-limiting and +5 net vote auto-verification
- Routes and Fares querying endpoints implemented (`GET /api/routes` and `GET /api/fare`)
- Heuristic Prediction engine implemented (`GET /api/predictions` for ETA, taxi, and congestion)
- Gemini AI Chatbot transit assistant implemented (`POST /api/chat` with unconfigured key fallback)
- Interactive Map Component implemented dynamically with custom styled Leaflet markers representing stops, taxis, and user location
- Automated unit tests, lint check, typecheck, and production build pass successfully
- Deployment configuration files added for frontend and backend
- Manual coordinate input for nearby search added
- OSRM-backed walking directions to selected stops added with fallback map guidance

## Pending Modules

- OSRM-backed route planning with walking segments, transfers, and congestion alternatives
- User profile personalization and trip history
- Full admin CRUD management and analytics
- Production deployment to frontend/backend/database providers
- System and performance testing against SRS targets

## Database Status

PostgreSQL/PostGIS is fully seeded and running on port 5433.
Core entities seeded:
- User
- transit_stops (Bus stops and taxi stands)
- routes
- fares
- crowdsource_reports
- crowdsource_votes

## API Status

Implemented endpoints:
- `GET /api/health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/profile/me`
- `GET /api/stops/nearby`
- `POST /api/stops/report`
- `POST /api/stops/vote`
- `GET /api/routes`
- `GET /api/fare`
- `GET /api/predictions`
- `POST /api/chat`

## Authentication Status

- Register with email/phone & bcrypt hashing: implemented
- Login with JWT authentication: implemented
- Role-based access control: middleware implemented
- Password reset: pending

## Important Notes

- Gemini explains predictions but must never run deterministic prediction.
- Empty regions display: `This area hasn't been mapped yet. Help your community by adding the first transit stop.`
- Crowdsourced reports auto-verify at +5 net votes.
- Users are rate-limited to 5 pending reports per 24 hours.

## Current Milestone

Milestone 7: MVP gap closure, testing, hardening, and deployment.
