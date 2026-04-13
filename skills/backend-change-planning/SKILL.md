---
name: backend-change-planning
description: Use when backend work requires contract-first thinking, schema changes, permission checks, side-effect analysis, or test planning.
---

# Backend Change Planning

Use this skill before implementing non-trivial backend changes.

## Check in this order

1. **Contract changes** — are request/response shapes, status codes, or error formats changing? List every endpoint affected. If the API is public or consumed by other teams, flag it as a breaking change risk.
2. **Schema and migration changes** — are tables, columns, indexes, or constraints being added, modified, or removed? Migrations must be:
   - Reversible (or explicitly flagged as one-way)
   - Safe under concurrent traffic (no locking whole tables in production)
   - Applied before code that depends on them
3. **Query and repository impact** — are queries changing? Check for N+1 patterns, missing indexes on new filter columns, and unbounded result sets.
4. **Validation and error flow** — is input validation complete at the boundary? Are error types and messages consistent with the existing codebase?
5. **Permission and ownership checks** — does the change enforce proper authorization? Can a user access or modify resources they don't own?
6. **Audit and other side effects** — does the change need to emit events, logs, notifications, or update caches? Are these handled transactionally or eventually?
7. **Tests** — for each item above, define the tests that verify correctness:
   - Contract tests (request/response shape validation)
   - Migration tests (up and down)
   - Permission boundary tests (authorized vs. unauthorized)
   - Error path tests (invalid input, dependency failure)

## Expected output

Produce a structured checklist with findings for each item:

1. contract impact — [endpoints affected, breaking changes]
2. schema and migration impact — [tables, columns, reversibility]
3. query impact — [new queries, indexes needed, performance notes]
4. validation and error handling — [boundary checks, error types]
5. permission and ownership impact — [auth rules, resource scoping]
6. audit and side effects — [events, notifications, cache invalidation]
7. required tests — [specific scenarios with expected outcomes]

Write "N/A — [reason]" for items that do not apply.

## Implementation order for backend changes

1. Migrations first — run and verify before deploying code that uses new schema
2. Shared types and interfaces second — DTOs, error types, shared validators
3. Repository/data layer third — new queries, updated repository methods
4. Service/business logic fourth — core behavior
5. Handlers/controllers fifth — wiring to HTTP or event triggers
6. Tests alongside each step — do not batch all tests to the end

## Common backend mistakes to verify against

- Missing database index on a new WHERE clause or JOIN column
- Nullable column added without handling NULL in existing queries
- Migration that locks a large table under traffic
- Error response format inconsistent with existing API conventions
- Transaction boundary too wide (holding locks longer than necessary) or too narrow (partial updates on failure)
- Missing rate limiting or pagination on new list endpoints

Maintain this checklist consistently with the repository guidance in `docs/agent-playbook.md`, and update any related agent/template docs together when intentional workflow changes are made.

## Conformance self-check

Before marking backend planning as complete, verify:

- [ ] All 7 checklist items are addressed (N/A with reason if not applicable)
- [ ] `DECISIONS.md` was read and no contradictions exist
- [ ] Implementation order follows: migrations → shared types → repository → service → handlers → tests
- [ ] Common backend mistakes were checked against
- [ ] For high-risk changes: plan was presented to user for approval before implementing
