# Example Profile: MVP / Rapid Mode

Use this profile for early-stage products prioritizing delivery speed.

## Priorities

1. Fast iteration
2. Tight scope and rapid feedback
3. Lightweight process with basic safety

## Suggested configuration

```yaml
trust_level: semi-auto
small_tasks_skip_compliance_block: true
targeted_tests_for_small_tasks: true
critic_required_only_for_large_changes: true
```

## Rules emphasis

- Keep changes small and reversible.
- Preserve mandatory validation loop, but prefer targeted tests first.
- Require architecture decisions only when boundaries or contracts change.
- Use concise summaries instead of verbose deliverables for small tasks.

## Recommended workflow

`application-implementer -> targeted validation -> brief summary`
