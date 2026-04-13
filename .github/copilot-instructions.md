# Repository instructions for AI coding agents

- Read `AGENTS.md` first.
- **Agent-deference principle**: this template only adds rules the agent tool does not already provide natively. Capabilities the agent handles out of the box (e.g., built-in safety rails, tool routing, output formatting) should not be re-specified here. See `docs/operating-rules.md` → Agent-deference principle.
- **Trust level**: the project defaults to `semi-auto`. Override with `supervised` or `autonomous` as needed. For fully unattended execution (no approval gates at all), set `trust_level: autonomous` + `dangerouslySkipAllCheckpoints: true`. See `docs/operating-rules.md` → Trust level.
- Treat named roles such as `feature-planner` or `risk-reviewer` as conceptual roles. If the tool cannot spawn named subagents, use the matching prompt template or local docs instead.
- Follow `docs/operating-rules.md` for safety, scope, and validation rules.
- Check `prompt-budget.yml` at the repo root for `execution_mode` (`supervised` or `autonomous`) before acting on checkpoint gates. See `docs/operating-rules.md` → Autonomous execution mode for gate behavior per mode.
- Use `feature-planner` for cross-module, ambiguous, contract, database, auth, security, or image-led flow changes.
- Use `backend-architect` for backend contract and domain work.
- Use `application-implementer` for general product or frontend implementation work.
- Use `ui-image-implementer` for screenshot-driven or mockup-driven UI work.
- Use `integration-engineer` to close wiring and flow gaps.
- Use `documentation-architect` for repository rules, ADRs, onboarding docs, and other durable documentation.
- Use `risk-reviewer` before finalizing behavior-changing or high-risk work.
- Use `critic` after a plan or proposal is produced, before user approval, to challenge design quality.
- Keep reusable guidance in version-controlled markdown instead of chat-only instructions.
- Each role should run in a separate context (subagent invocation). Pass structured handoff artifacts between roles, not raw conversation history.
- When loading instruction files, follow the four-layer loading order (static rules → stable skills → project state → volatile context) to maximize prompt cache hits. See `skills/prompt-cache-optimization/SKILL.md` for details.
- Follow layered configuration precedence from `docs/operating-rules.md`: Project Context → Domain Rules → Global Rules.

## Mandatory workflow

Before any implementation:
1. Discover the codebase (`skills/repo-exploration/SKILL.md`).
2. For first entry in a new repository, run `skills/on_project_start/SKILL.md` to perform environment scan and ask boundary questions.
3. Classify the task scale (`skills/demand-triage/SKILL.md`): Small, Medium, or Large.
4. Publish a compliance block (see `docs/operating-rules.md` → Mandatory first-response compliance block). At `semi-auto`/`autonomous` trust level, Small tasks may skip the compliance block.
5. Read `DECISIONS.md` and check for contradictions. For legacy modules, also search `DECISIONS_ARCHIVE.md` if it exists.
6. State assumptions, constraints, and proposed approach before writing code.
7. For behavior-changing work, define tests first per TDAI requirement in `docs/operating-rules.md`.

After planning (for Medium/Large work):
8. Invoke `critic` to challenge the plan; present plan + critique to user. At `supervised`/`semi-auto` trust level, wait for user approval before continuing. At `autonomous` trust level, critic review runs for Large tasks but does not block execution — the agent proceeds after the critique is produced.

After any code change:
9. Run the validation loop (`skills/test-and-fix-loop/SKILL.md`) — this runs autonomously (auto-fix without human approval). For Small tasks, run targeted tests only.
10. Use error recovery (`skills/error-recovery/SKILL.md`) if anything fails. Escalate to human after 3 consecutive failures.
11. Record decisions in `DECISIONS.md` when applicable.
12. If architecture changed, update ADRs (or decision log when ADR directory does not exist) in the same task.
13. If scope expands beyond the plan: at `supervised`/`semi-auto` trust level, STOP and present the expanded scope for approval. At `autonomous` trust level, log an advisory notice and continue.
14. Produce a task completion summary (`docs/agent-templates.md` → Task completion summary). At `semi-auto`/`autonomous` trust level for Small tasks, a brief summary suffices.
