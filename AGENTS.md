# Agent Playbook

Read these files before starting work:

0. `docs/rules-quickstart.md` — minimal rules for first-pass loading (read this first)
1. `docs/operating-rules.md` — mandatory safety, scope, validation, error recovery, and project-specific constraint rules
2. `docs/agent-playbook.md` — routing rules and role definitions
3. `docs/agent-templates.md` — reusable task and prompt templates as optional helpers

Read `docs/adoption-guide.md` when adapting this template to a new repository.
Read `docs/external-practices-notes.md` only when evolving the framework itself.
Read `docs/layered-configuration.md` when creating or refactoring Global/Domain/Project rules.
Read `docs/rule-optimization-plan.md` for the staged simplification and automation roadmap.

## Three-layer architecture

See `docs/agent-playbook.md` → Three-layer architecture for the full definition. Summary:

1. **Rules** — `docs/operating-rules.md` (safety, scope, validation, constraints)
2. **Skills** — `skills/*/SKILL.md` (12 skills: demand-triage, repo-exploration, test-and-fix-loop, error-recovery, memory-and-state, feature-planning, backend-change-planning, application-implementation, design-to-code, documentation-architecture, prompt-cache-optimization, on-project-start)
3. **Loop** — `Discover → Triage → Plan → Critique → Approve → Implement → Test → Fix → Repeat → Record → Summarize`. Steps **Plan**, **Critique**, **Approve** are trust-level-gated; see `docs/operating-rules.md` → Trust level.

## Configuration layering

Keep constraints in layered configuration form:

1. `rules/global/` — cross-project core rules
2. `rules/domain/` — domain-specific rules
3. `project/project-manifest.md` — project-local boundaries

See `docs/operating-rules.md` for precedence and conflict resolution.

Before implementation begins, the first response must include a compliance block: files/docs read, triage scale, selected path, and checkpoint expectations — except when trust level is `semi-auto`/`autonomous` and the task is Small. See `docs/operating-rules.md` → Mandatory first-response compliance block for the trust-level-specific requirements.

Core rules:

- Use a planning agent first for cross-module, ambiguous, high-risk, API, DB, auth, or security work.
- Use an application implementer for general product or frontend work that is not primarily backend architecture, pure integration, or image-led UI.
- Use a design-focused agent first for image-led UI implementation.
- Use a documentation-focused agent when the main output is repo rules, architectural notes, onboarding docs, ADRs, runbooks, or API/process documentation.
- Keep reusable instructions in version-controlled files, not only in chat history.
- Prefer specialized agents with clear ownership over one general-purpose agent.
- Never treat code as complete until the validation loop passes.
- Each role runs in its own context (separate invocation). Do not chain roles in one conversation. Pass structured handoff artifacts between roles, not raw history.

Source of truth:

- `docs/operating-rules.md` is the source of truth for safety, scope, validation, and review rules.
- `docs/agent-playbook.md` is the source of truth for role routing and role ownership.
- `docs/agent-templates.md`, `.claude/agents/`, `skills/`, and `.github/copilot-instructions.md` must stay aligned with those two files.
- `prompt-budget.yml` (at repo root) declares execution mode (`supervised` or `autonomous`), enabled roles, and token budget. See `docs/adoption-guide.md` → Autonomous execution mode for when and how to use autonomous mode.
