---
name: repo-exploration
description: Use before any implementation to understand the codebase structure, existing patterns, dependency graph, and project conventions.
---

# Repo Exploration

Use this skill to build a mental model of the codebase before making changes.

## Required steps

### 1. Structural discovery

- List the top-level directory structure.
- Identify the module map: where handlers, services, repositories, models, and tests live.
- Read `ARCHITECTURE.md`, `CONTRIBUTING.md`, or equivalent docs if they exist.

### 2. Pattern identification

Read at least one existing file in each layer you will modify and note:

- Naming convention (files, functions, variables, types)
- Error handling pattern (custom error types? wrapped errors? error codes?)
- Logging pattern (structured? which library?)
- Test convention (table-driven? test fixtures? mocks or fakes?)
- Import organization (stdlib first? grouped? aliased?)

### 3. Dependency graph

Before cross-file changes:

- Trace imports from the file you will change.
- Identify shared types, interfaces, and constants.
- Check for circular dependency risks.

### 4. Constraint check

- Read `docs/operating-rules.md` → `Project-specific constraints` section.
- Read `DECISIONS.md` if it exists.
- Do not introduce patterns that contradict recorded decisions.

## Output

After discovery, produce a brief summary:

```text
Modules involved: [list]
Existing patterns: [key observations]
Constraints: [project-specific rules that apply]
Files to read before coding: [list]
```

## Use this skill when

- Starting work on an unfamiliar module
- The task touches more than 2 files
- You are unsure about the project's conventions
- A previous attempt produced code that did not match the project style
