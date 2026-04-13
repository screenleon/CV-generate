---
name: test-and-fix-loop
description: Use after writing or modifying code to enforce the mandatory write → test → fix → repeat validation cycle.
---

# Test and Fix Loop

Use this skill to enforce iterative validation after every code change.

## The loop

```text
┌─────────────┐
│  Write code  │
└──────┬──────┘
       ▼
┌─────────────┐
│  Run tests   │ ← project test command (e.g., go test ./..., npm test, pytest)
└──────┬──────┘
       ▼
┌─────────────┐
│  Run lint    │ ← project lint command (e.g., go vet, eslint, mypy)
└──────┬──────┘
       ▼
   Pass? ──No──► Fix errors (minimal change) ──► Re-run from top
       │
      Yes
       ▼
┌─────────────┐
│    Done      │
└─────────────┘
```

## Rules

1. **Always run tests** — never skip this step, even for "obvious" changes.
2. **Run the narrowest scope first** — test only the affected package/module, then broaden if needed.
3. **Read full error output** — do not guess from truncated messages.
4. **Fix minimally** — change only what the error requires. Do not refactor during the fix cycle.
5. **Max 3 attempts per error** — if the same error persists after 3 fix attempts, escalate to the user.
6. **Never delete tests to pass** — if a test fails, fix the code or update the test expectation with justification.
7. **Report the final state** — after the loop converges, confirm: which tests ran, how many passed, any remaining warnings.

## Scale-adapted testing strategy

Adapt testing intensity based on the task scale from the `demand-triage` skill:

### Small tasks

- Run **only the tests for the changed file** and its direct dependents
- If no tests exist for the file, run the nearest module-level test suite
- A manual smoke check description is acceptable if no automated tests cover the change
- Still follow the fix loop if any test fails

### Medium tasks

- Run the **module-level test suite** covering all changed files
- Run lint/static analysis on changed files
- If the change adds new behavior, at least one new test is expected

### Large tasks

- Run the **full project test suite**
- Run full lint/static analysis
- New tests are mandatory for all new behavior
- Integration tests are expected if the change crosses module boundaries
- Consider running performance-sensitive tests if the change affects hot paths

## Test-first guidance

When adding new behavior (not fixing a bug in existing code), consider writing test expectations before implementation:

1. **Define expected behavior** — write a test that describes what the new code should do
2. **Run the test** — verify it fails for the right reason (missing implementation, not a test bug)
3. **Implement** — write the minimal code to make the test pass
4. **Refine** — add edge case tests, then adjust implementation as needed

This is a **recommendation, not a mandate**. Use test-first when:
- The expected behavior is clear and can be expressed as assertions before coding
- The change is a new function, endpoint, or feature (not a refactor or bug fix)
- The project has good test infrastructure that makes writing tests easy

Do not force test-first when:
- The change is exploratory and the API shape is still being discovered
- The test infrastructure would require significant setup to write a meaningful test
- The task is classified as Small and the change is trivial (e.g., config value update)

## Identifying the test command

Check in order:

1. `Makefile` or `Taskfile.yml` → look for `test`, `lint`, `check` targets
2. `package.json` → `scripts.test`, `scripts.lint`
3. `pyproject.toml` / `setup.cfg` → pytest configuration
4. `go.mod` → `go test ./...`
5. `Cargo.toml` → `cargo test`
6. `pom.xml` / `build.gradle` → `mvn test` / `gradle test`
7. Repository documentation (for example `README.md`, `CONTRIBUTING.md`, or the adoption guide) → validation or verification commands

If no test command is identifiable, state that explicitly.

## Use this skill when

- You have just written or modified code
- A previous code generation did not include a verification step
- The user reports that generated code does not compile or pass tests

## Conformance self-check

Before marking the validation loop as complete, verify:

- [ ] Tests were actually executed (not just planned)
- [ ] The test scope matches the task scale (Small: targeted, Medium: module, Large: full suite)
- [ ] All failures were addressed (fixed or escalated), not silently ignored
- [ ] The final test state is reported (tests run, pass count, warnings)
- [ ] No tests were deleted or disabled to achieve a passing state
