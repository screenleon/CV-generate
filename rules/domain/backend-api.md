# Domain Template: Backend API Rules

Use this template for backend API constraints reusable across repositories.

## Rule entries

Repeat this block for each rule.

```markdown
### Rule: <RULE_ID>
- Owner layer: Domain
- Domain: backend-api
- Status: active
- Scope: <module or surface>
- Statement: <clear non-negotiable rule>
- Rationale: <why>
- Verification: <how to verify>
- Supersedes: <RULE_ID or N/A>
- Superseded by: <RULE_ID or N/A>
```

## Starter examples

### Rule: API-001

- Owner layer: Domain
- Domain: backend-api
- Status: active
- Scope: public HTTP handlers
- Statement: All API responses must follow a single envelope contract.
- Rationale: Consumer compatibility and consistent error handling.
- Verification: Integration tests assert envelope fields for success and error responses.
- Supersedes: N/A
- Superseded by: N/A

### Rule: API-002

- Owner layer: Domain
- Domain: backend-api
- Status: active
- Scope: backward-compatible endpoint evolution
- Statement: Additive changes are allowed; breaking schema changes require versioning.
- Rationale: Prevent client breakage.
- Verification: Contract tests and changelog entry for any API version changes.
- Supersedes: N/A
- Superseded by: N/A
