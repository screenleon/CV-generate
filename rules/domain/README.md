# Domain Rules

Domain-specific constraints used to adapt the template for different technical areas.

Examples:

- Backend API contracts and compatibility rules
- Cloud infrastructure practices (AWS/Azure/GCP)
- Frontend component and state conventions
- Data platform and analytics pipeline conventions

When Domain Rules conflict with Project Context, Project Context wins.

## Starter templates

- `rules/domain/backend-api.md`
- `rules/domain/frontend-components.md`
- `rules/domain/cloud-infra.md`

## Recommended rule schema

For consistency and machine-readability, each rule entry should include:

- Rule ID
- Owner layer
- Domain
- Status (`active` or `superseded`)
- Scope
- Statement
- Rationale
- Verification
- Supersedes
- Superseded by
