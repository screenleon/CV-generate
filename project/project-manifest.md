# Project Manifest

Use this file to declare project-local boundaries that override generic guidance. Fill in each section when adopting the template. Leave blank fields as placeholders until confirmed during `on-project-start` initialization.

## Project identity

- Name:
- Repository type:
- Primary language(s):
- Runtime framework(s):

## Non-negotiable constraints

- Constraint 1:
- Constraint 2:
- Constraint 3:

## Build and validation commands

- Build:
- Unit tests:
- Integration tests:
- Lint/static analysis:

## Deployment and operations boundaries

- Environments:
- Release process:
- Incident/rollback rule:

## Security and compliance boundaries

- Secret handling:
- Auth/permission model:
- Data classification:

## Architecture context

- System style (monolith, modular monolith, microservices, etc.):
- Critical integration dependencies:
- Known technical debt:

## Override notes

- Any project rule that should override domain/global guidance:

## Override annotations

Use this format when project rules override base rules:

`Overrides: <base-rule-id> -> <project-rule-id>`

Example:

`Overrides: API-002 -> PROJECT-API-001`

## Override registry

| Base Rule ID | Project Rule ID | Reason | Status |
|---|---|---|---|
|  |  |  | active |
