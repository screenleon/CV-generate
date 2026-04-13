# Rule Optimization Plan

Objective: keep rules clear for humans and cheap for agents to load.

## Current assessment

- Core rules are complete but heavy (`docs/operating-rules.md` is large).
- Layered configuration exists and now has governance, but adoption can still drift without automation.
- Key behavior is spread across multiple docs; discoverability is acceptable but can be sharper for first-pass agent loading.

## Improvement roadmap

### Phase 1 (immediate) — completed

1. Keep `docs/rules-quickstart.md` as first-pass load target.
2. Keep source docs as canonical; avoid duplicating long prose in quickstart.
3. Require every new rule to declare owner layer (Global/Domain/Project).

### Phase 2 (short-term) — completed

1. Add domain rule templates:
   - `rules/domain/backend-api.md`
   - `rules/domain/frontend-components.md`
   - `rules/domain/cloud-infra.md`
2. Add an override annotation format in `project/project-manifest.md`:
   - `Overrides: <base-rule-id> -> <project-rule-id>`
3. Add a lightweight checklist section in PR template for rule updates.

### Phase 3 (automation) — in progress

1. Add CI linter for layered rules:
   - detect duplicate rule text across layers
   - detect unresolved same-layer conflicts
   - ensure superseded markers are explicit
2. Add periodic maintenance task (weekly or every 10 tasks):
   - check layer drift
   - prune stale overrides
3. Add quick report artifact in CI output:
   - active overrides count
   - duplicated rule candidates
   - unresolved conflicts

### Next concrete steps

1. ~~Add rule-id uniqueness checks across domain templates.~~ (completed)
2. ~~Add unresolved conflict markers check (for example `Status: superseded` without replacement target).~~ (completed)
3. Extend CI output with a short markdown summary artifact.

## Success criteria

- New agents can follow the repository by reading quickstart + one deep doc.
- Rule conflict resolution is deterministic and auditable.
- Rule drift is caught automatically before merge.
- Prompt load remains stable with minimal always-load content.
