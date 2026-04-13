# Layered Configuration Guide

Use this guide when defining or reviewing repository constraints across Global, Domain, and Project layers.

## Why this exists

A simple precedence statement is necessary but not sufficient. Teams also need clear placement criteria, deterministic conflict resolution, and maintenance hygiene to keep rules understandable over time.

## Layer responsibilities

| Layer | Owns | Should avoid |
|---|---|---|
| Global (`rules/global/`) | universal guardrails, baseline quality and safety | stack-specific implementation details |
| Domain (`rules/domain/`) | reusable constraints for a technical domain | repo-local deployment or legacy exceptions |
| Project (`project/project-manifest.md`) | repo-specific boundaries, operational constraints, exception handling | restating generic global/domain rules |

## Rule placement decision table

| Question | If yes | If no |
|---|---|---|
| Should this rule apply across nearly all repos? | Global | next question |
| Is this tied to one technical domain and reusable across repos in that domain? | Domain | next question |
| Is this specific to this repo/team/runtime/deployment context? | Project | Domain (fallback) |

## Resolution model

Use this sequence for every conflicting topic:

1. Project overrides Domain and Global.
2. Domain overrides Global.
3. Same-layer conflicts: choose narrower scope first.
4. If still tied: choose latest dated rule.
5. Document tie-break decisions in `DECISIONS.md`.

### Example

Global rule says "use REST for all APIs." Domain rule says "use gRPC for internal service communication." Project manifest says "legacy payment service uses SOAP; do not migrate."

Resolution: the payment service uses SOAP (project override); internal services use gRPC (domain override); all other APIs default to REST (global baseline).

## Governance checks

Run this checklist when updating layered rules:

1. No duplicate full-text rules across layers.
2. Every exception in Project points to the overridden base rule.
3. Superseded rules are explicitly marked as superseded.
4. Tool-specific instructions remain aligned with source-of-truth docs.
5. At least one concrete example exists for each active Domain profile.

## Anti-patterns

- Global layer contains framework-version constraints.
- Project manifest copy-pastes the entire global baseline.
- Domain layer mixes backend and frontend concerns into one broad file.
- Rule conflicts resolved implicitly without a recorded tie-break.

## Suggested maintenance cadence

- Per task: update layered files when constraints change.
- Weekly or every 10 tasks: run a quick conflict and duplication pass.
- Quarterly: prune obsolete domain profiles and superseded project exceptions.
