# Rules Quickstart

This is the minimal rule set for agents. Read this first, then expand into source docs only when needed.

## Source of truth

1. `docs/operating-rules.md` for safety, scope, validation, conflict handling
2. `docs/agent-playbook.md` for routing and role ownership
3. `DECISIONS.md` for active architectural constraints

## Trust level

Default: `semi-auto`. Override per project or session.

- `supervised` — all checkpoints require human approval
- `semi-auto` — Small/low-risk tasks run autonomously; checkpoints for Large/destructive work
- `autonomous` — proceeds without approval except destructive actions (unless `dangerouslySkipAllCheckpoints: true`)

See `docs/operating-rules.md` → Trust level for the full activation matrix.

## Layered configuration

Use three layers:

1. Global: `rules/global/`
2. Domain: `rules/domain/`
3. Project: `project/project-manifest.md`

Precedence: Project > Domain > Global.

If same-layer conflicts remain:

1. narrower scope wins
2. latest dated rule wins
3. record tie-break in `DECISIONS.md`

## Mandatory workflow (compact)

1. Discover (`skills/repo-exploration/SKILL.md`)
2. Initialize on first repo entry (`skills/on-project-start/SKILL.md`)
3. Triage (`skills/demand-triage/SKILL.md`)
4. Check contradictions in `DECISIONS.md` (and archive for legacy)
5. For behavior changes: define tests before implementation (TDAI)
6. Implement with minimal scope
7. Validate (test -> lint/typecheck -> fix -> repeat)
8. Record decisions and ADR updates when architecture changes

## Hard constraints

- Never expose credentials or secrets.
- Never do destructive actions without approval unless bypass mode is explicitly enabled.
- Follow existing repository practice unless user explicitly asks for refactor.

## Escalation points

- Contradiction with existing decision in `DECISIONS.md`
- Scope expansion beyond approved plan
- Same error persists after 3 fix attempts
- Architecture change without ADR/decision update

## When to open full docs

- Need trust-level gate details -> `docs/operating-rules.md`
- Need role routing details -> `docs/agent-playbook.md`
- Need layered governance details -> `docs/layered-configuration.md`
