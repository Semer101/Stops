# DECISIONS.md

# Architecture Decision Records

## ADR-0001: Treat SRS and SDS as Authoritative Sources

### Decision

All project implementation decisions must be checked against the SRS first, then the SDS, then `AGENTS.md`.

### Reason

The attached instructions define the SRS and SDS as the highest-priority sources of truth.

### Alternatives Considered

- Use `AGENTS.md` as the only implementation guide.
- Infer missing requirements from common smart-transport application patterns.

### Impact

Implementation must pause when requirements are missing or contradictory. This prevents invented behavior and keeps the project aligned with the documented academic deliverables.

## ADR-0002: Use the SDS Stack

### Decision

Use the SDS stack for implementation: Next.js 15, React, TypeScript, Tailwind CSS, Leaflet, Go/Gin, PostgreSQL/PostGIS, Google Gemini, OpenStreetMap, and OSRM.

### Reason

The SDS has priority over AGENTS.md, and the developer explicitly confirmed this stack on 2026-06-25 because it is faster, resource efficient, and easy to set up.

### Alternatives Considered

- React/Vite with Node.js/Express/Prisma from the attached AGENTS rulebook.
- A hybrid stack combining Next.js and Express.

### Impact

Frontend work will use the Next.js App Router. Backend work will use Go with Gin and a layered `cmd/api` plus `internal` package structure. Database access patterns and migration tooling must be chosen for Go rather than Prisma.

## ADR-0003: Use Plain SQL Migrations Compatible With golang-migrate

### Decision

Use versioned `.up.sql` and `.down.sql` migration files under `backend/migrations/`, compatible with `golang-migrate`.

### Reason

The confirmed backend stack is Go/Gin, and the SDS requires PostgreSQL with PostGIS. Plain SQL keeps PostGIS geometry columns, constraints, and indexes explicit while avoiding a new ORM requirement that is not present in the SDS.

### Alternatives Considered

- Prisma migrations from the lower-priority AGENTS stack.
- ORM-first migrations through a Go ORM.
- Manual database changes.

### Impact

Schema changes are reviewable SQL files. Developers must run migrations through a migration tool rather than editing database tables manually.
