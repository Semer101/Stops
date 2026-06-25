# AGENTS.md

# AI Smart Transport for Addis Ababa
## AI Development Rulebook

Version: 1.0

## Purpose

This repository contains the implementation of STOPS (Smart Transport Optimization and Planning System), also described in the attached documentation as AI Smart Transport for Addis Ababa.

This file defines the rules every AI coding assistant must follow before generating code.

The Software Requirements Specification (SRS) and Software Design Specification (SDS) are the single sources of truth.

If this document conflicts with the SRS or SDS, the SRS/SDS always take precedence.

## Authoritative Documents

Priority order:

1. Software Requirements Specification (SRS)
2. Software Design Specification (SDS)
3. AGENTS.md

Never invent requirements not found in these documents.

If information is missing or contradictory, ask the developer instead of making assumptions.

## Project Type

This project is a responsive full-stack web application.

It is not React Native, Flutter, Android, or iOS.

The application must run entirely in a modern web browser.

## Technology Stack

The approved implementation stack is the SDS stack:

- Frontend: Next.js 15, React, TypeScript, Tailwind CSS, Leaflet
- Backend: Go with Gin
- Database: PostgreSQL with PostGIS
- Prediction: lightweight heuristic engine
- AI/NLP: Google Gemini API
- Mapping: OpenStreetMap and OSRM
- Deployment: Docker/Docker Compose locally, Vercel frontend, Railway/Render/Koyeb backend, Neon PostgreSQL/PostGIS database

## Architecture

The architecture is strictly layered.

Frontend:

Pages -> Components -> Services -> API

Backend:

Controllers -> Services -> Repositories -> Database

Business logic belongs only inside the service layer.

Database access belongs only inside the repository layer.

Controllers must remain thin.

Never place business logic inside React components.

## Folder Structure

Approved frontend structure:

```text
frontend/
  app/
  components/
  services/
  hooks/
  types/
  utils/
  public/
```

Approved backend structure:

```text
backend/
  cmd/api/
  internal/
    config/
    controllers/
    middleware/
    repositories/
    routes/
    services/
    utils/
  migrations/
```

## Coding Standards

- Use TypeScript wherever the selected stack supports TypeScript.
- Enable strict typing.
- Avoid `any`.
- Prefer interfaces for API models.
- Use async/await.
- Handle all Promise errors.
- Never suppress TypeScript errors.

## Naming Conventions

- Variables: camelCase
- Functions: camelCase
- React components: PascalCase
- Interfaces: PascalCase
- Enums: PascalCase
- React component files: PascalCase.tsx
- Utility files: camelCase.ts

## UI Rules

- Use Tailwind CSS.
- Create reusable components.
- Avoid duplicated UI.
- Responsive design is mandatory for desktop, tablet, and mobile.
- Keep the UI map-centered and mobile-first as required by the SRS/SDS.

## API Rules

- REST only.
- Use JSON.
- Status codes must follow standards: 200, 201, 400, 401, 403, 404, 500.
- Never change an endpoint once defined.

Documented endpoints from the SDS:

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/stops/nearby`
- `GET /api/routes`
- `GET /api/fare`
- `POST /api/stops/report`
- `POST /api/stops/vote`
- `GET /api/predictions`
- `POST /api/chat`

## Database Rules

- Use PostgreSQL with PostGIS.
- Use migrations.
- Never modify tables manually.
- Never delete data during migrations.
- Use foreign keys.
- Use spatial indexes for location queries.

## Authentication Rules

- JWT authentication only.
- Passwords must be hashed using bcrypt.
- Never store plaintext passwords.
- Protect all private endpoints.
- Validate JWT before protected requests.
- Enforce role-based access control for admin endpoints.

## AI Features

Google Gemini is used only for natural-language explanations, travel recommendations, route comparisons, and chatbot responses.

Gemini must never directly modify database records or perform numerical prediction.

Numerical predictions belong to the heuristic prediction engine.

## Mapping

- Use OpenStreetMap.
- Use Leaflet.
- Use OSRM for routing.
- Do not replace mapping providers.

## Error Handling

Every API must return structured errors:

```json
{
  "success": false,
  "message": "...",
  "error": "..."
}
```

Never expose stack traces.

## Logging

Log server errors.

Do not log passwords, JWT tokens, secrets, or API keys.

## Security Rules

- Validate all input.
- Prevent SQL injection.
- Prevent XSS.
- Prevent CSRF where applicable.
- Sanitize user input.
- Use HTTPS in production.
- Treat user location and trip history as sensitive data.

## Performance

- Target API response time below 300 ms under normal conditions.
- Target route calculation below 2 seconds.
- Target map rendering below 2 seconds.
- Target Gemini response below 5 seconds.
- Target nearby stop search below 1 second with indexed spatial data.
- Lazy-load pages.
- Paginate large datasets.
- Avoid duplicate API requests.
- Optimize database queries.

## Git Rules

- Small commits.
- One feature per commit.
- Meaningful commit messages.
- Never rewrite unrelated files.

## Testing

Each feature should include practical unit tests, API/integration tests, and manual UI verification.

Fix failing tests before adding features.

## Hard Constraints

AI assistants must not:

- Introduce new frameworks without approval.
- Replace the approved frontend framework.
- Replace the approved backend framework.
- Replace PostgreSQL/PostGIS.
- Replace Tailwind.
- Replace Leaflet.
- Replace JWT.
- Rename database tables or API endpoints after they are defined.
- Change authentication flow.
- Change project architecture.
- Introduce breaking changes.

## Scope Rules

When implementing a feature:

- Modify only requested files.
- Do not refactor unrelated modules.
- Do not optimize unrelated code.
- Do not change APIs unless explicitly instructed.

## Definition of Done

A task is complete only if:

- It compiles successfully.
- It has no TypeScript errors where TypeScript applies.
- It has no lint errors.
- Responsive UI is verified.
- API behavior is documented.
- Error handling is implemented.
- It uses the existing architecture.
- It matches the SRS.
- It matches the SDS.

## Session Documentation

Before ending every work session, update:

- `PROJECT_CONTEXT.md`
- `TASKS.md`
- `CHANGELOG_AI.md`
- `DECISIONS.md` if decisions changed
- `SESSION_HANDOFF.md`

End of AGENTS.md
