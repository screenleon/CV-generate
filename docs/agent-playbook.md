# Agent Playbook

## Three-layer architecture

All agent work follows three layers:

1. **Rules** (`docs/operating-rules.md`) — hard constraints: safety, scope, agent-deference, trust level, codebase discovery, validation loop, error recovery, project-specific constraints, decision log.
2. **Skills** (`skills/*/SKILL.md`) — reusable capabilities: repo exploration, test-and-fix loop, error recovery, memory management, prompt cache optimization, plus domain skills (planning, backend, frontend, design, docs).
3. **Loop** — every implementation follows: Discover → **Triage** → Plan → **Critique** → **Approve** → Implement → Test → Fix → Repeat → Record → **Summarize**. Steps in **bold** are trust-level-gated; see `docs/operating-rules.md` → Trust level for activation rules.

## Layered configuration model

In addition to the execution architecture above, repository constraints should be organized into:

1. **Global Rules** — `rules/global/`
2. **Domain Rules** — `rules/domain/`
3. **Project Context** — `project/project-manifest.md`

Precedence follows `docs/operating-rules.md` → Layered configuration and Conflict resolution principle.

## Repository asset map

- Global entrypoint: `AGENTS.md`
- Project subagents: `.claude/agents/*.md`
- Reusable templates: `docs/agent-templates.md`
- Reusable skills: `skills/*/SKILL.md`
- Repo-wide Copilot instructions: `.github/copilot-instructions.md`
- Decision log: `DECISIONS.md` (created per repo)
- Architecture overview: `ARCHITECTURE.md` (created per repo)

## Source of truth and precedence

Use this precedence order when documents overlap:

1. `docs/operating-rules.md` for safety, scope control, validation, and destructive-action rules
2. `docs/agent-playbook.md` for routing, role definitions, and workflow ownership
3. `AGENTS.md` as the short root entrypoint into those two files
4. `docs/agent-templates.md` as reusable prompt scaffolds
5. `.claude/agents/`, `skills/`, and `.github/copilot-instructions.md` as tool-specific implementations of the same role model

If a tool-specific file drifts from the source-of-truth docs (`docs/operating-rules.md` and this playbook), update the tool-specific file to match them.

## Tool portability

The role names in this template are conceptual. Different tools expose them differently:

- Claude-style tooling can map them into `.claude/agents/*`
- Copilot-style tooling can reference them through repository instructions and prompt files
- Codex-style or generic chat tooling can use the same role names through prompt templates and local repo docs

Do not assume every tool supports named subagents. Keep the role model stable even when the implementation surface changes.

## Default routing

### Use the planning agent first when

- a request impacts more than one module
- a request changes API contracts, schemas, migrations, events, or background jobs
- a request touches auth, permissions, audit, uploads, security, or notifications
- a request is still ambiguous and needs scope, order, or risk clarification
- a request is driven by screenshots or mockups and also changes flow or state

### Use specialist agents directly when

- backend contract and domain work is isolated
- general application or frontend implementation is isolated and does not need a planning-first phase
- image-led UI implementation is isolated
- integration work is mostly wiring existing pieces together
- documentation is the primary deliverable
- final review is focused on bugs, security, and regressions

## Role definitions

### `feature-planner`

- defines scope, non-goals, impacted modules, dependencies, order, and validation
- owns ambiguity reduction before implementation starts

### `backend-architect`

- owns contract-first backend design, schema changes, permissions, audit, and high-risk backend behavior

### `application-implementer`

- owns general product implementation that is neither pure backend architecture nor mostly integration wiring
- covers ordinary frontend, service-layer, or app behavior work where a dedicated image-led flow is unnecessary

### `ui-image-implementer`

- owns design-to-code tasks driven by screenshots, mockups, or visual specs

### `integration-engineer`

- owns wiring across API, state, navigation, side effects, caching, and complete user journeys

### `documentation-architect`

- owns repository instructions, onboarding docs, ADRs, runbooks, process docs, and architecture explanations
- optimizes for long-term maintainability and future agent readability
- responsible for automatic maintenance of `DECISIONS.md`, `ARCHITECTURE.md`, and project-specific constraints as a side effect of code changes

### `risk-reviewer`

- owns bug finding, regression detection, permission review, security review, and testing gaps
- also provides **early risk assessment during planning** for high-risk work (schema migrations, auth changes, payment logic, public API changes, cross-service changes)

### `critic`

- adversarial design reviewer invoked **after** a planner or architect produces a proposal and **before** the user decides
- challenges proposals for over-engineering, hidden coupling, missing edge cases, constraint violations, scope creep, and unstated assumptions
- does not rewrite proposals — states what is wrong and lets the proposer fix it
- separate from `risk-reviewer`: critic challenges design quality; risk-reviewer checks implementation safety

## Suggested workflow

### Mandatory steps for all workflows

Every workflow below implicitly includes these steps:

1. **Discover** — run the `repo-exploration` skill before coding
2. **Initialize (new repo entry)** — run `skills/on-project-start/SKILL.md` to scan stack signals and collect missing boundary constraints before implementation
3. **Triage** — run the `demand-triage` skill to classify task scale (Small / Medium / Large) based on evidence from discovery. This determines which subsequent steps are mandatory vs. optional. See `skills/demand-triage/SKILL.md` for classification criteria and workflow adaptation rules
4. **Structured preamble** — state assumptions, constraints, and proposed approach before producing output (see `docs/operating-rules.md` structured output rules). For Small tasks, this may be inline (1–2 sentences)
5. **Test-first for new behavior** — follow TDAI in `docs/operating-rules.md` by defining test cases before implementing behavior-changing work
6. **Validate** — run the `test-and-fix-loop` skill after every code change. For Small tasks, run only targeted tests for the changed file
7. **Recover** — use the `error-recovery` skill when anything fails
8. **Record** — use the `memory-and-state` skill to log decisions, update architecture docs, and check whether memory lifecycle maintenance is needed (see `skills/memory-and-state/SKILL.md` → Memory lifecycle management)
9. **ADR sync** — for architecture changes, update ADRs or decision records in the same task (`docs/operating-rules.md` → ADR automatic update)
10. **Cache-aware loading** — follow the instruction loading order in `skills/prompt-cache-optimization/SKILL.md` to maximize prefix cache hits
11. **Isolate** — each role runs in a separate context. Pass structured handoff artifacts between roles, not raw conversation history (see Context isolation section below). Small tasks need only one agent. Medium tasks at `semi-auto` or `autonomous` may relax isolation per `docs/operating-rules.md` → Task boundary rule
12. **Deliver** — produce output using the mandatory deliverable structure (see `docs/operating-rules.md` → Mandatory deliverable structure). For Small tasks, keep the required structure concise rather than replacing it
13. **Summarize** — after completing any task, produce a brief task completion summary for memory (see `docs/agent-templates.md` → Task completion summary). This summary is additional to the required deliverable structure and enables future pattern reuse and prevents context loss across sessions
14. **Feedback loop** — include a mini retrospective and quality-signal update as defined in `docs/operating-rules.md` → Feedback loop and quality signals

### First-response compliance block

Publish the compliance block defined in `docs/operating-rules.md` → Mandatory first-response compliance block. Required at `supervised` trust level for all tasks and at `semi-auto` for Medium/Large tasks. Optional at `autonomous` trust level.

### Checkpoint gates

Checkpoint activation depends on trust level. See `docs/operating-rules.md` → Checkpoint activation matrix for the full table. Key gates: destructive actions (always), scope expansion (`supervised` and `semi-auto`), plan approval (`supervised` always, `semi-auto` Large only).

### New feature

`feature-planner` → `critic` → **user decision** → `backend-architect`, `application-implementer`, and/or `ui-image-implementer` → `integration-engineer` → `documentation-architect` as needed → `risk-reviewer`

### High-risk backend change

`feature-planner` → `critic` → `risk-reviewer` (plan assessment) → **user decision** → `backend-architect` → `risk-reviewer` (final review)

### Small change

If the `demand-triage` skill classifies the task as Small:

`application-implementer` (with inline 1–2 sentence plan) → targeted validation only

No planning agent, critic, or risk-reviewer required. The implementer reads the file, states the change in 1–2 sentences, implements, and runs targeted tests.

At `semi-auto` and `autonomous` trust levels, the agent proceeds directly without waiting for approval. The validation loop runs autonomously.

At `supervised` trust level, the following remain explicit and mandatory:

1. First-response compliance block
2. Structured preamble (inline 1–2 sentences is acceptable)
3. DECISIONS.md contradiction check outcome
4. Validation plan and targeted verification result
5. Mandatory deliverable structure (concise is allowed; omission is not)

### General application change

If it is bounded and low ambiguity (Medium scale):

`application-implementer` → `risk-reviewer`

If it also changes flow, state, or contracts:

`feature-planner` → `critic` → **user decision** → `application-implementer` → `integration-engineer` → `risk-reviewer`

### Image-led UI change

If it is visual only:

`ui-image-implementer` → `risk-reviewer`

If it also changes logic or flow:

`feature-planner` → `critic` → **user decision** → `ui-image-implementer` → `integration-engineer` → `risk-reviewer`

### Documentation-heavy change

`feature-planner` as needed → **user approval** → `documentation-architect` → `risk-reviewer` when technical correctness matters

### Autonomous workflow variants

When `execution_mode: autonomous` is set in `prompt-budget.yml`, replace **user decision** and **user approval** steps with auto-proceed. All other steps remain unchanged.

**Important**: The loop structure (Discover → Triage → Plan → Critique → Implement → Test → Fix → Repeat → Record → Summarize) still executes in full. Only the human wait states are removed.

| Supervised workflow | Autonomous equivalent |
|--------------------|-----------------------|
| `feature-planner` → `critic` → **user decision** → implementers → `risk-reviewer` | `feature-planner` → `critic` (critique embedded in handoff) → _(auto-proceed, logged)_ → implementers → `risk-reviewer` |
| `feature-planner` → `critic` → `risk-reviewer` (plan) → **user decision** → `backend-architect` → `risk-reviewer` (final) | `feature-planner` → `critic` → `risk-reviewer` (plan; stop if severity-high finding) → _(auto-proceed, logged)_ → `backend-architect` → `risk-reviewer` (final) |
| `feature-planner` as needed → **user approval** → `documentation-architect` → `risk-reviewer` | `feature-planner` as needed → _(auto-proceed, logged)_ → `documentation-architect` → `risk-reviewer` |

**Retained hard stops in autonomous mode** (see `docs/operating-rules.md` → Autonomous execution mode):

- Destructive or irreversible actions (gate 2) — always stop
- Stuck escalation after 3 failed attempts (gate 4) — always stop
- Contradiction detected in `DECISIONS.md` — always stop
- Severity-high finding from `risk-reviewer` during plan assessment — always stop

## Feedback loop execution

Use feedback loop outputs to improve process wording and reduce repeat misses.

### Minimum output per completed task

After the task completion summary, include:

1. Friction observed
2. Miss risk
3. Most useful rule
4. Next improvement candidate

### Cadence and ownership

- Review quality signals every 10 tasks (or weekly)
- `documentation-architect` owns wording updates and synchronization when recurring friction is detected
- `risk-reviewer` should flag repeated process misses even when code-level outcomes are correct

## Context isolation

Each agent role must run in its own context (separate invocation, session, or subagent call). Do not chain roles in a single long conversation.

### Why

Role switching within one context causes:
- **Context drift** — the agent forgets which role it is playing
- **Long-task loss of control** — instructions from early in the conversation are ignored
- **Memory contamination** — reasoning from one role leaks into and distorts the next

### How

- Each step in a workflow is a **separate agent invocation**.
- Agents communicate through **handoff artifacts** (see `docs/operating-rules.md` → Context isolation → Handoff artifact), not through shared conversation history.
- If the tool does not support separate sessions, insert a hard context break: summarize the output into a handoff artifact and restart with only that artifact.

### Workflow with context boundaries

```text
[Context 1] feature-planner → produces plan artifact
[Context 2] critic → receives plan artifact → produces critique artifact
[User]      reviews plan + critique → decides
[Context 3] backend-architect → receives approved plan → produces implementation
[Context 4] risk-reviewer → receives implementation summary → produces review
```

Each `[Context N]` is an isolated invocation. No context carries forward except through explicit handoff artifacts.

## Ownership principles

- Planning agents define scope, order, dependencies, and validation.
- Implementation agents stay inside their domain and avoid unnecessary expansion.
- Integration agents close loops across state, navigation, side effects, and data flow.
- Documentation agents keep instructions, architecture notes, and operational docs aligned with the actual workflow.
- Review agents lead with findings, not summaries.

## Maintenance principles

- Keep root guidance short and stable.
- Put details in focused docs, agents, and skills.
- Promote repeated prompts into reusable templates.
- Keep templates generic unless a repository-specific constraint truly matters.
- Prefer one conceptual role model with many tool-specific implementations, not many unrelated role models.
