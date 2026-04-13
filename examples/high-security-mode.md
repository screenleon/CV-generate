# Example Profile: High-Security Mode

Use this profile for regulated or security-critical systems.

## Priorities

1. Security and compliance before velocity
2. Strict permission and approval controls
3. Maximum auditability and traceability

## Suggested configuration

```yaml
trust_level: supervised
require_risk_reviewer_for_all_changes: true
require_adr_for_architecture_change: true
block_destructive_actions_without_manual_approval: true
```

## Rules emphasis

- Enforce least-privilege everywhere.
- Require threat-model review for auth, permissions, and external integrations.
- Require test-first for all behavior changes, including medium-risk bug fixes.
- Require decision log entry for every security-relevant tradeoff.

## Recommended workflow

`feature-planner -> critic -> risk-reviewer -> user approval -> implementer -> risk-reviewer`
