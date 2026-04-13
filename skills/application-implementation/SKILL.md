---
name: application-implementation
description: Use for general product implementation work that is not primarily backend architecture, pure integration wiring, or screenshot-driven design-to-code.
---

# Application Implementation

Use this skill for ordinary application changes — frontend, service-layer, or app behavior work that does not require a full architecture pass or image-led flow.

## Check in this order

1. **Target behavior** — what user-visible behavior will change? State it in one sentence.
2. **Bounded edit scope** — list the files that need edits. If the list exceeds the approved plan, stop and request approval for the expanded scope.
3. **State coverage** — for each UI or API path affected, confirm you handle:
   - loading state
   - empty state (no data)
   - error state (network failure, validation error, permission denied)
   - success state
4. **Whether planning or integration help is needed** — if you discover cross-module dependencies or contract changes that were not in the plan, escalate.
5. **Verification path** — identify the specific test commands to run after changes.

## Implementation guidelines

- **Match existing patterns** — use the same component structure, naming conventions, state management approach, and error handling style already in the codebase. Do not introduce new patterns without justification.
- **Minimize blast radius** — change only what is needed. If a refactor would help, propose it separately rather than bundling it with the feature.
- **Handle errors at every boundary** — network calls, file I/O, user input parsing, and type conversions are all boundaries. Each needs explicit error handling.
- **Write tests before marking complete** — at minimum, test the happy path and one error path. Use the project's existing test patterns.
- **Update documentation** — if the change introduces a new pattern or modifies a documented behavior, update `ARCHITECTURE.md` or `DECISIONS.md` as described in the `documentation-architecture` skill.

## Common mistakes to avoid

- Forgetting to handle the loading state (UI flickers or shows stale data)
- Swallowing errors silently (empty catch blocks, ignoring promise rejections)
- Expanding scope without checking with the user
- Skipping the validation loop after changes
- Introducing a new state management pattern when one already exists

## Anti-patterns and negative examples

These are concrete examples of what **not** to do. Use them as a mental checklist.

### Do not guess file structure

Bad: Assuming `src/components/Button.tsx` exists because it sounds right.
Good: Use `repo-exploration` to discover the actual file paths before making changes.

### Do not add abstractions for single use

Bad: Creating a `FormValidationHelper` class to validate one field in one form.
Good: Inline the validation logic. Extract only when a second consumer appears.

### Do not change error handling style mid-codebase

Bad: Introducing `try/catch` with custom error classes when the codebase uses `Result<T, E>` pattern.
Good: Follow the existing error handling convention, even if you prefer a different style.

### Do not bundle refactoring with feature work

Bad: Renaming variables and restructuring imports while adding a new feature.
Good: Ship the feature first. Propose refactoring as a separate task.

### Do not ignore the demand triage classification

Bad: Running the full planning → critic → review workflow for a one-line copy change.
Good: Check the scale classification. For Small tasks, use the lightweight path.

## Use this skill when

- the task is ordinary product implementation
- the task is frontend or app behavior work without a dedicated visual source
- the task needs code changes but not a new architecture pass

## Conformance self-check

Before marking implementation as complete, verify:

- [ ] The change matches the stated target behavior (or inline plan for Small tasks)
- [ ] All affected states are handled (loading, empty, error, success) where applicable
- [ ] No scope expansion occurred without approval
- [ ] The validation loop was run and passed
- [ ] Existing patterns were followed (no new patterns introduced without justification)
- [ ] DECISIONS.md was checked for contradictions
- [ ] A task completion summary was produced (see `docs/agent-templates.md`)
