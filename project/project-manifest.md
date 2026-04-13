# Project Manifest

Use this file to declare project-local boundaries that override generic guidance. Fill in each section when adopting the template. Leave blank fields as placeholders until confirmed during `on-project-start` initialization.

## Project identity

- Name: CV Generator
- Repository type: Full-stack application
- Primary language(s): Go (backend), HTML/CSS/JavaScript (frontend)
- Runtime framework(s): Go `net/http` stdlib

## Non-negotiable constraints

- All CV generation must happen server-side in Go (no client-side PDF libs).
- `CVData.Template` and `CVData.Format` control the output; frontend must always send both.
- DOCX output must be a valid ZIP archive (PK header) readable by Microsoft Word and LibreOffice.

## Build and validation commands

- Build: `cd backend && go build -o cv-generator .`
- Unit tests: `cd backend && go test ./...`
- Integration tests: none yet
- Lint/static analysis: `cd backend && go vet ./...`

## Deployment and operations boundaries

- Environments: local dev (default port 8080)
- Release process: build binary, set `FRONTEND_DIR` and `PORT` env vars, run binary
- Incident/rollback rule: stateless — restart binary to recover

## Security and compliance boundaries

- Secret handling: no secrets required; no database
- Auth/permission model: none (open — suitable for local or internal use)
- Data classification: CV data is transient — never stored on server

## Architecture context

- System style: single-process monolith
- Critical integration dependencies: `github.com/go-pdf/fpdf` for PDF; Go stdlib for DOCX
- Known technical debt: font support for CJK characters in PDF requires embedding a Unicode font (current implementation uses Helvetica which does not render Japanese kanji)

## Override notes

- N/A

## Override annotations

Use this format when project rules override base rules:

`Overrides: <base-rule-id> -> <project-rule-id>`

Example:

`Overrides: API-002 -> PROJECT-API-001`

## Override registry

| Base Rule ID | Project Rule ID | Reason | Status |
|---|---|---|---|
|  |  |  | active |
