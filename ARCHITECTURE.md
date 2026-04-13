# Architecture Overview

> **Adopter note**: Fill this file in when adopting this template. This file is intentionally blank so teams start with a clean architecture description. Keep it updated as the codebase evolves.
>
> Agents read this file before working on unfamiliar modules (see `skills/memory-and-state/SKILL.md` → Architecture memory). When it is missing or stale, agents lose structural context and may make incorrect assumptions.
>
> Minimum viable content: fill in Module map and Data flow. Add the remaining sections as the project matures.

## Module map

| Directory / module | Purpose |
|-------------------|---------|
| _example: `src/api/`_ | _HTTP handlers — one file per resource_ |
| _example: `src/services/`_ | _Business logic — stateless service objects_ |
| _example: `src/repos/`_ | _Database access — one repo per aggregate_ |
| _example: `db/migrations/`_ | _SQL migration files (up/down pairs)_ |

<!-- Replace the example rows with your actual directory structure. -->

## Data flow

<!-- Describe how data moves through the system for the primary user actions.
     A simple linear diagram is fine.

     Example:
       Client → HTTP handler (src/api/) → Service (src/services/) → Repository (src/repos/) → PostgreSQL
       Background jobs: Scheduler → Worker (src/workers/) → Service → DB

     Add notes for async paths, event-driven flows, or external service calls. -->

_Not yet documented. Fill in when adopting._

## Key interfaces and contracts

<!-- List the major interfaces, shared types, or public API surfaces that multiple modules depend on.
     Agents use this to understand the blast radius of changes.

     Example format:
     - `UserService` (src/services/user.ts) — owns user creation, lookup, and auth token generation
     - `OrderRepository` (src/repos/order.ts) — owns all SQL for orders; no direct DB access outside this file
     - `POST /api/orders` — external contract; breaking changes require versioning -->

_Not yet documented. Fill in as key interfaces are identified._

## External service dependencies

| Service | Purpose | Notes |
|---------|---------|-------|
| _example: PostgreSQL_ | _Primary data store_ | _Connection via DATABASE_URL env var_ |
| _example: SendGrid_ | _Transactional email_ | _Only used in notification service_ |

<!-- Add any external APIs, message queues, caches, object stores, or third-party services here.
     Include where credentials come from and which internal module owns the integration. -->

## Deployment units

<!-- If the project has multiple deployable units (monorepo with separate services, packages, or apps),
     list them here so agents know which changes are cross-unit (higher risk).

     Example:
     - `apps/api` — Node.js REST API, deployed to Fly.io
     - `apps/worker` — Background job processor, deployed as a separate Fly machine
     - `packages/shared` — Shared TypeScript types; changes here affect both apps -->

_Single deployable unit — not yet documented._

## Known technical debt

<!-- Optional but valuable. Record known shortcuts or deferred work so agents do not mistake
     intentional debt for bugs, and do not introduce more of the same pattern.

     Example:
     - User search uses a full-table ILIKE scan (no index). Known N+1 with >10k users. Deferred until load justifies indexing.
     - Order status transitions are hardcoded strings, not an enum. Refactor tracked in issue #42. -->

_None documented yet._
