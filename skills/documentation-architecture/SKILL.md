---
name: documentation-architecture
description: Use when the main deliverable is maintainable documentation such as repository rules, onboarding guides, runbooks, ADRs, or architecture notes.
---

# Documentation Architecture

Use this skill when documentation is the primary output, or when code changes require documentation to stay in sync.

## Define before writing

1. audience — who will read this? (developers, agents, ops, new hires)
2. source of truth — which file is canonical for this topic?
3. mandatory versus optional guidance — what must be followed vs. what is a recommendation?
4. what stays short versus what moves into focused docs
5. tool-specific files that need alignment

## Automatic documentation maintenance

Documentation rots when it is only updated manually. These rules ensure agents keep docs current as a side effect of normal work, not as a separate task.

### `DECISIONS.md` — auto-append after decisions

Trigger: any agent makes an architectural or behavioral decision during implementation.

Action: append a new entry using the format in `docs/operating-rules.md` → Decision log.

Agents must not skip this step. If a decision was made but not recorded, the risk-reviewer should flag it.

### `ARCHITECTURE.md` — auto-update after structural changes

Trigger: any of these events during implementation:
- A new module or directory is created
- A module is moved, renamed, or deleted
- A module boundary changes (e.g., a service is split or merged)
- A new external dependency or integration is added
- Data flow between modules changes

Action: update the relevant section of `ARCHITECTURE.md`. If the file does not exist, create it with this structure:

```markdown
# Architecture

## Module map

| Directory          | Purpose                          |
|--------------------|----------------------------------|
| src/api/           | HTTP handlers and route definitions |
| src/services/      | Business logic                   |
| ...                | ...                              |

## Key interfaces and contracts

- [Interface name] — [file path] — [what it defines]

## Data flow

[Brief description or diagram of how data moves through the system]

## External dependencies

| Dependency   | Purpose        | Notes                    |
|--------------|----------------|--------------------------|
| PostgreSQL   | Primary store  | Managed via migrations   |
| ...          | ...            | ...                      |

## Known technical debt

- [Description] — [file path or module] — [why it exists]
```

### `Project-specific constraints` — auto-append after constraint discovery

Trigger: during implementation, an agent discovers an unwritten rule that is enforced by the codebase (e.g., "all handlers use middleware X", "dates are always UTC").

Action: add it to the `Project-specific constraints` section in `docs/operating-rules.md`.

### Documentation sync check

After any code change that affects architecture, contracts, or decisions, agents must verify:

1. `DECISIONS.md` has entries for all decisions made in this task
2. `ARCHITECTURE.md` reflects any structural changes
3. `docs/operating-rules.md` project-specific constraints include any newly discovered rules
4. Tool-specific files (`.claude/agents/`, `.github/copilot-instructions.md`) are still aligned with the source-of-truth docs

If any are stale, update them before marking the task complete.

## Writing guidelines

- Keep each doc focused on one topic. Split rather than append endlessly.
- Use tables for structured data (module maps, decisions, risks).
- Prefer concrete examples over abstract principles.
- Date all decision entries.
- Write for the next agent session, not just the current one — assume no prior context.

## Use this skill when

- writing repository instructions
- updating onboarding or process docs
- generating ADRs, runbooks, or architecture notes
- keeping agent-facing and human-facing docs aligned
- code changes require documentation sync (structural changes, new decisions, new constraints)
