# Decision Log

This file records architectural and behavioral decisions that affect future work.
Agents must read this file before planning or implementation tasks.
See `skills/memory-and-state/SKILL.md` for when to read and write.

> **Adopter note**: This file is intentionally blank in the template. Start adding your own decisions from day one. Template development history has been moved to `DECISIONS_ARCHIVE.md` so you inherit a clean log.
>
> The fastest way to get value: keep one entry per architectural choice so every future agent run can perform contradiction checks.

<!-- Append new decisions below using the format shown. -->
<!-- Do not remove or silently contradict existing entries. -->
<!-- To reverse a decision, add a new entry that explicitly references and supersedes the old one. -->

## YYYY-MM-DD: [Example] Template decision format

- **Context**: Why this decision was needed
- **Decision**: What was decided
- **Alternatives considered**: What was rejected and why
- **Constraints introduced**: What future work must respect

---

<!-- Add your real decisions below this line -->

## 2026-04-13: Go backend + vanilla JS frontend for CV generator

- **Context**: Issue requested a CV generator with Go backend, frontend form, and multi-template PDF/Word output.
- **Decision**: Go standard-library HTTP server (`net/http`) serves both the REST API (`POST /api/generate`) and the static frontend files. No framework dependency. PDF uses `github.com/go-pdf/fpdf`. DOCX is generated as raw ZIP+XML (Go stdlib `archive/zip`) — no Word library required.
- **Alternatives considered**: Echo/Gin frameworks (rejected — adds dependency for a single endpoint). `unioffice` for DOCX (rejected — commercial license, heavyweight). React/Vue for frontend (rejected — overkill for a form with a single API call).
- **Constraints introduced**: Adding new output formats requires a new generator function in `backend/generators/` and a new case in `handlers/cv_handler.go`. Template styles are selected via `CVData.Template` field; new templates extend the existing switch in each generator.

---

## 2026-04-13: Two template styles — "simple" and "japan"

- **Context**: Issue explicitly requested support for Japan (履歴書) and simple/Western CV styles.
- **Decision**: `simple` uses a two-column header + section-separator layout suited for Western CVs. `japan` uses a 2-column table for personal info and section tables styled after the standard Japanese 履歴書 (Rirekisho) format, including Japan-specific fields (birth date, gender, nationality).
- **Alternatives considered**: Storing templates as external files (rejected — increases deployment complexity). HTML-to-PDF conversion (rejected — requires headless browser or wkhtmltopdf, complicates deployment).
- **Constraints introduced**: Japan template exposes `birth_date`, `gender`, and `nationality` fields; these should remain optional so the form degrades gracefully when the simple template is selected.

---

## 2026-04-13: Added 職務経歴書 (Shokumu Keirekisho) template for detailed work history

- **Context**: User feedback requested support for the Japanese 職務経歴書 (Shokumu Keirekisho) format — a detailed professional resume focused on work experience, technical skills, and contributions. This is distinct from the basic 履歴書 (Rirekisho) and is more important for engineers in Japanese job applications.
- **Decision**: Added `"shokumu"` as a third template option. Extended `CVData` model with `TechStack` (categorized skills: languages, frameworks, databases, infrastructure, tools) and extended `Experience` struct with `Project`, `Role`, and per-job `TechStack` fields. Template emphasizes work history details over personal info (only name, email, phone shown in compact header). Skills section uses categorized table format when `TechStack` is provided.
- **Alternatives considered**:
  - Merging shokumu features into the existing `japan` template (rejected — 履歴書 and 職務経歴書 serve different purposes and have different layouts).
  - Using a single `skills` array for all templates (rejected — 職務経歴書 requires categorized skills for readability, per Japanese conventions).
- **Constraints introduced**:
  - `TechStack` and detailed experience fields (`project`, `role`, `tech_stack`) are optional and primarily used by the `shokumu` template.
  - Frontend shows/hides shokumu-specific fields based on template selection.
  - Adding future Japanese document types (e.g., カバーレター cover letter) should follow the same template extension pattern.
