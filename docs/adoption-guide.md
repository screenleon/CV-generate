# Adoption Guide

Use this guide when adapting the template into a new repository.

## Minimum installation

Copy these files first:

- `AGENTS.md`
- `docs/operating-rules.md`
- `docs/agent-playbook.md`

Then choose any of the following that match your toolchain:

- `.claude/agents/` for Claude-compatible project subagents
- `.github/copilot-instructions.md` for GitHub Copilot repository instructions
- `docs/agent-templates.md` for generic reusable prompts
- `skills/` for reusable workflow packaging

## Critical first customization (do not skip)

The template will produce vague, unusable agent behavior until you fill in these items:

### 0. Layered configuration scaffolding (mandatory)

Create and maintain these paths:

- `rules/global/` for cross-project core rules
- `rules/domain/` for domain-specific rules
- `project/project-manifest.md` for project-local boundaries

Precedence is: Project Context -> Domain Rules -> Global Rules.

For placement criteria, conflict resolution details, and anti-patterns, follow `docs/layered-configuration.md`.

### 1. Project-specific constraints (mandatory)

Open `docs/operating-rules.md` and fill the `Project-specific constraints` section with your actual rules. Be concrete and specific — agents cannot follow "best practices", they can only follow explicit instructions.

Good examples:

```markdown
- Use raw SQL with sqlc; no ORM
- Do not modify the DB schema without a migration file in db/migrations/
- Pricing logic must use JSONB rule definitions in the pricing_rules table
- All HTTP handlers must use the shared middleware stack in internal/middleware/
- Authentication uses JWT with RS256; do not switch algorithms
- Frontend state management uses Zustand; do not introduce Redux
- All API responses follow the envelope format in docs/api-envelope.md
- Error codes must be registered in internal/errors/codes.go
```

Bad examples (too vague for agents):

```markdown
- Follow best practices
- Write clean code
- Use appropriate patterns
```

### 2. Validation commands (mandatory)

Agents need to know exactly which commands to run. Add these to `docs/operating-rules.md` or a `CONTRIBUTING.md`:

```markdown
## Validation commands
- Unit tests: `go test ./...`
- Lint: `golangci-lint run`
- Type check: N/A (Go)
- Integration tests: `make test-integration`
- Build: `go build ./cmd/...`
```

### 3. Decision log (recommended)

Create a `DECISIONS.md` at the repo root:

```markdown
# Decisions

## 2024-01-15: Use sqlc instead of GORM
- Why: Type safety, predictable queries, no magic
- Constraint: All queries must be in .sql files under db/queries/

## 2024-02-01: JSONB for pricing rules
- Why: Business needs frequent rule changes without migrations
- Constraint: Pricing logic reads from pricing_rules table, never hardcoded
```

### 4. Architecture overview (recommended)

Create an `ARCHITECTURE.md` or fill it into your README to help agents understand the codebase structure:

```markdown
## Module map
- cmd/api/         → HTTP server entrypoint
- internal/handler/ → HTTP handlers
- internal/service/ → Business logic
- internal/repo/    → Database access (sqlc)
- db/migrations/    → SQL migrations
- db/queries/       → sqlc query files
```

### 5. Initialization protocol (recommended)

Add the `skills/on-project-start/SKILL.md` workflow to your first-session process so agents dynamically discover missing boundaries.

Expected first-session questions include:

- Framework-specific naming or package conventions
- Infra tooling choices (for example Terraform vs CDK)
- Deployment/runtime constraints that are not documented yet

## First customization pass

Edit these items immediately:

- repository module names
- safety rails
- validation commands and expectations
- project-specific constraints (see above)
- role names you do or do not want to keep
- review and merge expectations

## Optional simplification

If your team is small or your project is simple, you can remove roles you do not need.

Good candidates to remove early:

- `backend-architect` in a frontend-only repo
- `ui-image-implementer` in a backend-only repo
- `documentation-architect` if documentation is not a repeated workflow yet

## Prompt budget trimming

The full template loads all skills, roles, and format templates in every request. When adopted by a real project, trimming unused content reduces per-request token cost significantly.

### Step 1: Remove unused role templates

Delete agent definitions and templates for roles your project does not use. Common removals:

| If your project is… | Safe to remove |
|----------------------|----------------|
| Frontend-only | `backend-architect` agent + template |
| Backend-only | `ui-image-implementer` agent + template, `skills/design-to-code/` |
| Small team (1–2 devs) | `critic`, `risk-reviewer` (inline reviews instead) |
| No documentation workflow | `documentation-architect` agent + template, `skills/documentation-architecture/` |

Update `docs/agent-playbook.md` routing table and `.github/copilot-instructions.md` after removing roles.

### Step 2: Remove unused skills

Each skill adds ~1,500–3,500 tokens when loaded. Delete skill folders that do not match your workflow:

- `skills/design-to-code/` — only needed for screenshot-driven UI work
- `skills/documentation-architecture/` — only needed for docs-as-deliverable workflows
- `skills/feature-planning/` — can be removed if all work is Small/Medium scale
- `skills/backend-change-planning/` — can be removed in frontend-only projects

### Step 3: Simplify format templates

After your team is comfortable with the workflow, consider condensing verbose format templates.

If you change or remove any template structure here, also update the corresponding requirements in `docs/operating-rules.md` (and any other repo rules documents) so the documented rules stay consistent with the templates your team is expected to use.

- **Deliverable structure** — if your team always produces the same 3 sections, reduce the 5-section template to match
- **Checkpoint format** — if your tool enforces approval natively, the checkpoint template can be shortened
- **Handoff artifact** — if you rarely chain agents, this can be removed entirely

### Step 4: Configure Layer 2 loading strategy

Use the prompt cache optimization skill (`skills/prompt-cache-optimization/SKILL.md`) to choose which skills load per task type. For small projects:

- **Minimal set**: `demand-triage` + `repo-exploration` only (all task types)
- **Standard set**: add `test-and-fix-loop` for implementation tasks
- **Full set**: use the canonical skill sets from the prompt cache skill (recommended for teams > 3)

### Step 5: Set a token budget with `prompt-budget.yml`

Create a `prompt-budget.yml` at the repo root to declare your project's prompt configuration:

```yaml
# prompt-budget.yml — Prompt budget configuration for this project
# Agents read this file to determine which skills and roles to load.

budget:
  layer1_target_tokens: 3000    # Target for static rules (operating-rules + agent-playbook)
  layer2_max_tokens: 6000       # Max for skills loaded per request
  layer3_max_tokens: 3000       # Max for DECISIONS.md + ARCHITECTURE.md

roles:
  enabled:
    - feature-planner
    - application-implementer
    - risk-reviewer
    - critic
  disabled:
    - backend-architect          # Not needed: frontend-only project
    - ui-image-implementer       # Not needed: no design-to-code workflow
    - documentation-architect    # Not needed: docs are informal

skills:
  always_load:
    - demand-triage
    - repo-exploration
  on_demand:
    - test-and-fix-loop
    - error-recovery
    - memory-and-state
    - prompt-cache-optimization
  disabled:
    - design-to-code             # Not needed
    - documentation-architecture # Not needed
    - feature-planning           # All tasks are Small/Medium
    - backend-change-planning    # Frontend-only

trimming:
  decisions_archive_threshold_kb: 30
  decisions_archive_threshold_entries: 50
  session_memory_max_files: 10
```

This file is informational — agents use it as guidance to select which skills and role templates to load. It does not enforce hard limits but makes the intended budget visible and auditable.

### Trimming impact estimate

| Action | Estimated savings per request |
|--------|-------------------------------|
| Remove 1 unused role template | ~200–400 tokens |
| Remove 1 unused skill | ~1,500–3,500 tokens |
| Condense format templates | ~300–500 tokens |
| Configure Layer 2 minimal set | ~3,000–6,000 tokens vs full set |
| Total (aggressive trim) | ~5,000–10,000 tokens |

For a project running 50 agent requests/day, aggressive trimming can save 250K–500K tokens/day.

## Recommended maintenance loop

1. Update the source-of-truth docs first.
2. Sync tool-specific files second.
3. Review whether repeated prompts should graduate into templates or skills.
4. Remove stale roles rather than letting them drift.
5. Update `DECISIONS.md` when architectural decisions are made or reversed.
6. Review `Project-specific constraints` quarterly — remove stale rules, add new ones.
7. Run a **memory health check** whenever a health indicator triggers (>50 entries, >30 KB, >10 session files), or at least quarterly for low-volume projects:
   - If `DECISIONS.md` exceeds 50 entries or 30 KB, archive inactive decisions to `DECISIONS_ARCHIVE.md` (see `skills/memory-and-state/SKILL.md` → Memory lifecycle management)
   - Purge session memory files that were not promoted to repo memory
   - Verify no archived constraint is still referenced by current code

## Autonomous execution mode

By default, agents pause at checkpoint gates and wait for human approval ("PROCEED"). For teams that want agents to operate without human confirmation — for example in CI/CD pipelines or prototyping environments — the template supports an autonomous mode that replaces human wait states with automatic logging and proceeding.

### What autonomous mode changes

Autonomous mode removes the **human wait states**, not the **work steps**. The agent still discovers the codebase, classifies scale, plans, critiques, implements, validates, and records decisions. The only change is that the agent proceeds automatically instead of pausing for your "PROCEED" signal.

| Step | Supervised mode | Autonomous mode |
|------|----------------|-----------------|
| Plan approval | Agent stops and waits | Agent logs plan to `DECISIONS.md`, proceeds |
| Scope expansion | Agent stops and presents | Agent logs expansion, proceeds if within original intent |
| Mid-implementation review | Agent pauses (recommended) | Agent skips |
| Before-merge review | Agent pauses (recommended) | Agent skips |

### What autonomous mode does NOT change

These remain default/recommended stop conditions in autonomous mode unless you explicitly relax the corresponding `prompt-budget.yml` flags:

- **Destructive or irreversible actions** — file deletion, table drops, force-push, branch reset. By default these require your approval when `halt_on_destructive_actions: true`. Only set `false` in fully isolated sandbox environments.
- **Stuck escalation** — if an error persists after 3 fix attempts, the agent stops and reports when `halt_on_stuck_escalation: true`. Do not disable this unless you have an external timeout mechanism in place.
- **DECISIONS.md contradictions** — if a proposed change conflicts with an existing decision, the agent stops and presents both sides. Auto-resolution is never allowed. This rule has no configuration override.
- **Severity-high risk findings** — if the risk-reviewer finds a severity-high issue during plan assessment, the agent stops when `halt_on_high_severity_risk: true`. Only set `false` if the codebase is sandboxed and risk findings are pre-acknowledged.

### How to enable autonomous mode

**Step 1**: In `prompt-budget.yml`, change `execution_mode` from `supervised` to `autonomous`:

```yaml
execution_mode: autonomous

autonomous_mode:
  auto_proceed_on_plan: true
  auto_proceed_on_scope_expansion: true
  halt_on_destructive_actions: true      # Keep true unless fully sandboxed
  halt_on_stuck_escalation: true         # Always keep true
  skip_critic_role: false                # Keep false for better quality
  halt_on_high_severity_risk: true       # Keep true for safety
```

**Step 2**: Ensure `DECISIONS.md` is in a good state before enabling autonomous mode. The agent will auto-log decisions here — a messy decision log will produce noisy entries.

**Step 3**: If your project uses destructive operations as a normal part of its workflow (e.g., a data-migration script that drops temporary tables), document those in `Project-specific constraints` in `docs/operating-rules.md` so the agent knows which destructive operations are pre-approved.

**Step 4**: Review `DECISIONS.md` after the first few autonomous runs to verify the auto-logged entries are sensible.

### Risk tradeoffs

| Risk | Mitigation |
|------|-----------|
| Agent makes a wrong architectural decision without human review | Plan is logged to `DECISIONS.md`; review logs post-hoc and add a correcting entry if needed |
| Scope creeps silently | Gate 3 (scope expansion within original intent) is still logged; unrelated module additions always stop |
| Destructive action executed automatically | Gate 2 is non-bypassable by default (`halt_on_destructive_actions: true`) |
| Critic findings ignored | Critique is embedded in the handoff artifact; implementers must address each point |
| Agent loops on errors | Gate 4 is non-bypassable by default (`halt_on_stuck_escalation: true`) |

### When not to use autonomous mode

- Schema migrations on production data
- Permission or security model changes
- Payment, billing, or financial logic
- Any task that will be run unreviewed in a production environment
- First run on a new codebase (discover patterns in supervised mode first)

## Tool adapter reference

The role model in this template is conceptual. Use the table below to find the right integration surface for each tool.

| Tool | System-level instructions | Per-task instructions | Subagent / role support |
|------|--------------------------|----------------------|------------------------|
| **Claude Code** | `.claude/agents/*.md` auto-loaded per role | `skills/*/SKILL.md` referenced on demand | Named subagents via `.claude/agents/` |
| **GitHub Copilot** | `.github/copilot-instructions.md` auto-injected | `.github/prompts/*.prompt.md` invoked via `#` | No native subagents; use prompt files as role templates |
| **Cursor** | `.cursor/rules/*.mdc` or `.cursorrules` | Inline `@`-mentioned files | No native subagents; use rules files as role templates |
| **Windsurf** | `.windsurfrules` | Referenced files | No native subagents; use rules file per role |
| **Custom OpenAI API** | `system` message (Layer 1+2) | `user` message prefix (Layer 3+4) | No native subagents; spawn separate API calls per role |
| **Codex CLI** | `AGENTS.md` + `docs/operating-rules.md` via repo context | Role templates from `docs/agent-templates.md` | No native subagents; use prompt templates |

### Cursor setup

1. Create `.cursor/rules/` directory (or use `.cursorrules` at the repo root for older versions).
2. Create one `.mdc` file per always-loaded instruction, e.g.:
   - `.cursor/rules/operating-rules.mdc` — paste or reference `docs/operating-rules.md`
   - `.cursor/rules/agent-playbook.mdc` — paste or reference `docs/agent-playbook.md`
3. For role-specific behavior, create a rule file per role and use it in the relevant context.
4. Reference skills by asking the agent to read the relevant `skills/*/SKILL.md` file at the start of a task.

### Windsurf setup

1. Create `.windsurfrules` at the repo root.
2. Include the core instructions from `AGENTS.md`, `docs/operating-rules.md`, and `docs/agent-playbook.md`.
3. For token efficiency, summarize only the most critical rules and link to full files for on-demand reading.
4. Skills and role templates work as referenced files — ask the agent to read them as needed.

### Custom OpenAI API setup

1. Place Layer 1 (`docs/operating-rules.md` + `docs/agent-playbook.md`) in the `system` message.
2. Place Layer 2 (selected skills) at the start of the `user` message, before the task description.
3. Place Layer 3 (`DECISIONS.md`, `ARCHITECTURE.md`) after Layer 2 in the `user` message.
4. Place the actual task query and current file content last (Layer 4).
5. For multi-role workflows, spawn separate API calls for each role using the same Layer 1 `system` message to maximize cache hits.

See `skills/prompt-cache-optimization/SKILL.md` → Tool-specific adaptation for more detail on cache-aware loading per tool.
