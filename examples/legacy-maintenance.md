# Example Profile: Legacy Maintenance

Use this profile for older stacks where stability and compatibility are more important than modernization.

## Priorities

1. Minimize regression risk
2. Respect existing conventions
3. Constrain change scope aggressively

## Suggested configuration

```yaml
trust_level: semi-auto
prefer_existing_code_practice: true
require_explicit_approval_for_pattern_replacement: true
legacy_modules_require_archive_decision_search: true
```

## Rules emphasis

- Follow existing codebase practice by default.
- Do not introduce new frameworks or architectural patterns without explicit approval.
- Require contradiction checks against both `DECISIONS.md` and `DECISIONS_ARCHIVE.md` for legacy areas.
- Encourage characterization tests before touching fragile modules.

## Recommended workflow

`repo-exploration -> demand-triage -> implement minimally -> targeted validation -> decision update`
