# Domain Template: Cloud Infrastructure Rules

Use this template for cloud infrastructure constraints reusable across repositories.

## Rule entries

Repeat this block for each rule.

```markdown
### Rule: <RULE_ID>
- Owner layer: Domain
- Domain: cloud-infra
- Status: active
- Scope: <module or surface>
- Statement: <clear non-negotiable rule>
- Rationale: <why>
- Verification: <how to verify>
- Supersedes: <RULE_ID or N/A>
- Superseded by: <RULE_ID or N/A>
```

## Starter examples

### Rule: INFRA-001

- Owner layer: Domain
- Domain: cloud-infra
- Status: active
- Scope: infrastructure as code
- Statement: All infrastructure changes must go through version-controlled IaC files.
- Rationale: Repeatability and auditability.
- Verification: CI validates IaC plan/apply checks before merge.
- Supersedes: N/A
- Superseded by: N/A

### Rule: INFRA-002

- Owner layer: Domain
- Domain: cloud-infra
- Status: active
- Scope: secret management
- Statement: Runtime secrets must come from managed secret stores; never from committed files.
- Rationale: Security and compliance.
- Verification: Secret scan and deployment policy checks.
- Supersedes: N/A
- Superseded by: N/A
