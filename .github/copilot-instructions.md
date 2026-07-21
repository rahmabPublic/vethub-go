# VetHub Go – Project Instructions for AI Coding Agents

## Project overview
- Monorepo with two apps:
  - `server/`: Go REST API for a veterinary clinic domain.
  - `client/`: SvelteKit frontend consuming the API.
- Domain entities: `owners`, `pets`, `visits`, `vets`, `specialties`.

## Backend context (`server/`)
- Entry point: `server/cmd/server/main.go`.
- DB initialization: `server/internal/db/db.go`.
  - Uses in-memory SQLite (`modernc.org/sqlite`).
  - Loads schema and seed from embedded SQL files:
    - `server/internal/db/schema.sql`
    - `server/internal/db/seed.sql`
- API routes are organized by domain:
  - `server/internal/api/owner.go`
  - `server/internal/api/pet.go`
  - `server/internal/api/visit.go`
  - `server/internal/api/vet.go`
- Shared API helpers:
  - `server/internal/api/respond.go` (JSON/error helpers)
  - Prefer existing helpers (`JSON`, `Error`, `BadRequest`, `NotFound`, etc.) instead of ad-hoc responses.
- Middleware:
  - `server/internal/middleware/auth.go` uses Basic Auth (`user` / `password` by default).
  - `server/internal/middleware/cors.go` allows localhost-origin CORS patterns.

## Frontend context (`client/`)
- Framework: SvelteKit + TypeScript + Vite.
- Routes live under `client/src/routes/`.
- API layer:
  - Typed OpenAPI client: `client/src/lib/api/client.ts`
  - Domain controllers: `client/src/lib/api/*/*Controller.ts`
  - Shared API model aliases: `client/src/lib/api/models.ts`
- UI structure:
  - Reusable forms/components in `client/src/lib/components/`.
  - Shared layout in `client/src/lib/components/layout/`.
- Auth is handled by sending Basic Auth headers from `client/src/lib/api/client.ts`.

## API/OpenAPI workflow
- API sync scripts are in `client/package.json`:
  - `download:api`: fetches docs from backend.
  - `generate:api`: regenerates TypeScript types.
  - `sync:api`: project helper script.
- When backend request/response shapes change, update OpenAPI artifacts and keep frontend typings in sync.

## Coding conventions for this repository
- Keep changes surgical and scoped to the task.
- Reuse existing patterns before introducing new abstractions.
- Preserve route naming and domain grouping style used in current API files.
- Keep error responses consistent with existing API error codes/messages.
- Prefer explicit validation and clear error handling (no silent failures).
- For frontend, keep business logic in route pages/controllers and keep form components reusable/presentational.

## Local development assumptions
- Backend runs on `http://localhost:8080`.
- Frontend uses `VITE_SERVER_BASE_URL` fallback to `http://localhost:8080` and appends `/api`.
- Default demo credentials are expected unless overridden by env vars:
  - `VITE_API_USERNAME`
  - `VITE_API_PASSWORD`
