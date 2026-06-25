# STOPS

Smart Transport Optimization and Planning System for Addis Ababa.

STOPS is a responsive full-stack web application that helps commuters locate nearby bus stops and taxi stands, plan routes, estimate fares and travel times, receive AI-assisted explanations, and contribute missing transit data.

## Project Status

Initial documentation and planning only. No application code has been implemented yet.

Implementation stack has been confirmed. The repository is scaffolded for Next.js 15 and Go/Gin.

## Authoritative Documents

Priority order:

1. Software Requirements Specification (SRS)
2. Software Design Specification (SDS)
3. `AGENTS.md`

Attached source documents read during initialization:

- `C:\Users\PC\Downloads\SRS_Final.pdf`
- `C:\Users\PC\Downloads\SDS_Final.pdf`

## Key MVP Features

- User registration and JWT authentication
- Interactive OpenStreetMap map
- Nearby bus stop and taxi stand discovery
- Route planning with walking directions and transfers
- Fare estimation
- Travel time estimation
- Heuristic prediction engine for bus ETA, taxi availability, and congestion
- Gemini-powered transit assistant
- Crowdsourced transit stop mapping and voting
- User profile, saved places, preferences, and trip history
- Admin dashboard for transport data and submission moderation
- Empty-region experience for unmapped areas

## Folder Structure

```text
/
  frontend/
    app/
    components/
    services/
    hooks/
    types/
    utils/
  backend/
    cmd/
    internal/
      controllers/
      services/
      repositories/
      middleware/
      routes/
      config/
      utils/
  docker-compose.yml
```

## Implementation Roadmap

1. Scaffold frontend, backend, database, and local development tooling.
2. Implement authentication and RBAC.
3. Implement PostgreSQL/PostGIS schema and seed data.
4. Implement interactive map and nearby stop discovery.
5. Implement route planning, fare estimation, and travel time estimation.
6. Implement heuristic prediction engine.
7. Integrate Gemini AI assistant.
8. Implement crowdsourced stop reporting and voting.
9. Implement user profile and personalization.
10. Implement admin dashboard and audit logging.
11. Add tests, performance checks, security checks, and deployment configuration.

## Module Breakdown

- Authentication
- User Profile
- Location and Nearby Search
- Route Planning
- Fare Estimation
- Travel Time Estimation
- Crowdsourcing
- Bus Arrival Prediction
- Taxi Availability Prediction
- Congestion Prediction
- AI Orchestration with Gemini
- Recommendation Engine
- Admin Dashboard

## Development Milestones

- Milestone 0: Documentation and planning
- Milestone 1: Repository scaffold
- Milestone 2: Authentication and database foundation
- Milestone 3: Map, stops, and route planning
- Milestone 4: Fare, travel time, and prediction engine
- Milestone 5: Gemini assistant and recommendations
- Milestone 6: Crowdsourcing and admin dashboard
- Milestone 7: Testing, hardening, and deployment

## Local Development

Frontend:

```bash
cd frontend
pnpm install
pnpm dev
```

Backend:

```bash
cd backend
go mod tidy
go run ./cmd/api
```

Database:

```bash
docker compose up db
```

Migrations:

```bash
cd backend
go run ./cmd/migrate
```

Checks:

```bash
cd frontend
pnpm run typecheck
pnpm run lint
pnpm run build

cd ../backend
go test ./...
```

## Confirmed Stack

The developer confirmed the SDS stack on 2026-06-25: Next.js 15 with Go/Gin.
