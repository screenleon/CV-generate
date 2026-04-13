# Domain Template: Frontend Component Rules

Use this template for frontend component constraints reusable across repositories.

## Rule entries

Repeat this block for each rule.

```markdown
### Rule: <RULE_ID>
- Owner layer: Domain
- Domain: frontend-components
- Status: active
- Scope: <module or surface>
- Statement: <clear non-negotiable rule>
- Rationale: <why>
- Verification: <how to verify>
- Supersedes: <RULE_ID or N/A>
- Superseded by: <RULE_ID or N/A>
```

## Starter examples

### Rule: UI-001

- Owner layer: Domain
- Domain: frontend-components
- Status: active
- Scope: shared UI components
- Statement: Shared components must be stateless unless local UI state is required.
- Rationale: Predictable rendering and easier testing.
- Verification: Component tests and lint rules for prop-driven behavior.
- Supersedes: N/A
- Superseded by: N/A

### Rule: UI-002

- Owner layer: Domain
- Domain: frontend-components
- Status: active
- Scope: form interactions
- Statement: All form components must expose loading, error, and disabled states.
- Rationale: Consistent UX and accessibility.
- Verification: Visual or interaction tests cover each state.
- Supersedes: N/A
- Superseded by: N/A
