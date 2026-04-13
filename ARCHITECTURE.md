# Architecture Overview

> **Adopter note**: Fill this file in when adopting this template. This file is intentionally blank so teams start with a clean architecture description. Keep it updated as the codebase evolves.
>
> Agents read this file before working on unfamiliar modules (see `skills/memory-and-state/SKILL.md` → Architecture memory). When it is missing or stale, agents lose structural context and may make incorrect assumptions.
>
> Minimum viable content: fill in Module map and Data flow. Add the remaining sections as the project matures.

## Module map

| Directory / module | Purpose |
|-------------------|---------|
| `backend/` | Go HTTP server — entry point (`main.go`) |
| `backend/handlers/` | HTTP request handlers (one file per resource) |
| `backend/generators/` | Document generation logic — PDF and Word (DOCX) |
| `backend/models/` | Shared data models (`CVData`, `Experience`, `Education`, `Language`) |
| `frontend/` | Static web frontend (HTML + CSS + vanilla JS) |
| `frontend/css/` | Stylesheet (`style.css`) |
| `frontend/js/` | Frontend logic (`app.js`) — form collection and API calls |

## Data flow

```text
Browser (frontend/index.html)
  → POST /api/generate  (JSON: CVData)
    → backend/handlers/GenerateCV
      → backend/generators/GeneratePDF  or  GenerateWord
        → returns []byte (PDF / DOCX)
  ← binary file download (Content-Disposition: attachment)
```

Static assets (`/`, `/css/*`, `/js/*`) are served directly by the Go HTTP
server from the `frontend/` directory.

## Key interfaces and contracts

- `models.CVData` (`backend/models/cv.go`) — single source of truth for all CV fields accepted by the API
- `POST /api/generate` — accepts `application/json` body (CVData), returns PDF or DOCX binary
- `generators.GeneratePDF(data *models.CVData) ([]byte, error)` — PDF generation entry point
- `generators.GenerateWord(data *models.CVData) ([]byte, error)` — DOCX generation entry point
- `CVData.Template` — controls layout: `"simple"` (modern Western), `"japan"` (履歴書 Rirekisho), or `"shokumu"` (職務経歴書 Shokumu Keirekisho)
- `CVData.Format` — controls output: `"pdf"` or `"word"`/`"docx"`

## External service dependencies

| Service | Purpose | Notes |
|---------|---------|-------|
| `github.com/go-pdf/fpdf` | PDF generation | Pure-Go library, no external dependencies |
| None (stdlib `archive/zip`) | DOCX generation | DOCX is a ZIP+XML format built with Go stdlib |

## Deployment units

- Single deployable unit: Go HTTP server (`backend/`) serves both the REST API and the frontend static files.
- **Build**: `cd backend && go build -o cv-generator .`
- **Run**: `FRONTEND_DIR=../frontend ./cv-generator` (default port 8080, override with `PORT` env var)

## Known technical debt

<!-- Optional but valuable. Record known shortcuts or deferred work so agents do not mistake
     intentional debt for bugs, and do not introduce more of the same pattern.

     Example:
     - User search uses a full-table ILIKE scan (no index). Known N+1 with >10k users. Deferred until load justifies indexing.
     - Order status transitions are hardcoded strings, not an enum. Refactor tracked in issue #42. -->

_None documented yet._
