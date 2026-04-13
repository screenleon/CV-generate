---
name: on-project-start
description: Use on first entry to a new repository to run environment scanning and ask targeted boundary questions before implementation.
---

# On Project Start

Use this skill only when entering a repository for the first time in a session or when project boundaries are still unclear.

## Goal

Convert unknown project boundaries into explicit, confirmed constraints before implementation starts.

## Protocol

1. Scan repository signals
   - Inspect top-level files and dependency manifests (e.g., `package.json`, `pom.xml`, `go.mod`, `requirements.txt`, `Cargo.toml`, `Gemfile`).
   - Detect likely stack and runtime (for example Spring Boot, Django, React, Terraform, CDK).
   - Detect validation commands from CI config (`.github/workflows/`, `Makefile`, `Justfile`, scripts).
   - Check for existing conventions docs (`CONVENTIONS.md`, `ARCHITECTURE.md`, `.editorconfig`).

2. Ask targeted boundary questions
   - Ask only high-value questions that can change implementation decisions.
   - Prefer either/or questions when possible.
   - Skip questions whose answers are already documented in existing project files.
   - Examples:
     - "I detected Spring Boot. Is there a package naming standard I must follow?"
     - "I detected AWS IaC files. Should infrastructure changes use Terraform or CDK?"
     - "I see no test framework configured. Which testing approach should I use?"

3. Confirm and persist
   - Summarize confirmed boundaries in concise bullets.
   - Add or update `project/project-manifest.md` with the confirmed constraints.
   - If a new durable architectural constraint is discovered, append to `DECISIONS.md`.

4. Continue with normal workflow
   - Run `repo-exploration` then `demand-triage`.
   - Do not start implementation until initialization questions are resolved or explicitly deferred by user instruction.

## Output format

```text
Initialization summary:
- Detected stack:
- Boundary questions asked:
- User confirmations:
- Constraints recorded in project manifest:
- Open items (if any):
```

## Conformance self-check

- [ ] Stack signals were discovered from real files, not guessed
- [ ] At least one boundary question was asked when constraints were unclear
- [ ] Confirmed constraints were written to project context
- [ ] Workflow proceeded to discovery/triage before implementation
