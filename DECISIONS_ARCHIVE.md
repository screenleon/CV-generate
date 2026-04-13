# Decisions Archive

This file holds decisions that were made during the development of this template itself (not of any adopting project). They are preserved here for historical context and to explain the design rationale behind the template's rules and structure.

**Adopters**: You should not need this file for your day-to-day work. Use the active `DECISIONS.md` for your own project decisions. Only read this archive if you are modifying the template framework itself (e.g., changing `docs/operating-rules.md`, `docs/agent-playbook.md`, or the skill files) and need to understand why a structural choice was made.

Agents should search this archive only when working on the template's own framework files — not for ordinary adopter tasks.

See `skills/memory-and-state/SKILL.md` → Selective read strategy for when to read this file.

---

## 2026-04-09: Memory lifecycle management for persistent memory files

- **Context**: Persistent memory files (especially `DECISIONS.md`) grow unboundedly over time. A real-world instance reached ~20,000 characters, consuming significant context tokens and slowing agent discovery. No archival, pruning, or selective-read mechanism existed.
- **Decision**: Introduced a memory lifecycle management system with four mechanisms: (1) decision archive rules — when entry count exceeds 50 or file size exceeds 30 KB, entries whose constraints are no longer enforced move to `DECISIONS_ARCHIVE.md`; (2) selective read strategy — tiered reading (always read active decisions, search archive only for legacy tasks); (3) session memory hygiene — unpromoted session notes are deleted at session end; (4) memory health indicators — thresholds trigger maintenance actions on demand, with quarterly review as a backstop for low-volume projects.
- **Alternatives considered**: Date-based auto-archive (e.g., >90 days) — rejected because age alone does not indicate irrelevance; high-volume projects can exceed thresholds within days, making fixed time windows inadequate. Multi-file topic-based split — rejected because it fragments the decision log and increases discovery cost.
- **Constraints introduced**: Archive operations require safety checks (search for references before moving); archived decisions retain full content (no lossy compression); agents must search archive before concluding no prior decision exists for legacy module work; memory health checks are triggered by threshold indicators, with quarterly review as a low-volume backstop.
- **Archived on**: 2026-04-11 — Template development history; constraints are structural rules for the template framework itself, not for adopter projects.

## 2026-04-08: Feedback loop governance for process stability and quality

- **Context**: Workflow rules were strengthened, but sustained quality requires a closed loop that continuously captures friction, measures adherence, and triggers wording updates when recurring misses appear.
- **Decision**: Introduced a formal feedback loop requirement with task-end mini retrospectives, rolling quality signals, and a recurrence escalation rule (3 repeated misses triggers source-of-truth update + synchronization + changelog entry).
- **Alternatives considered**: Keep feedback collection informal in ad-hoc comments only. Rejected because informal feedback is hard to aggregate, easy to forget, and does not reliably trigger policy updates.
- **Constraints introduced**: Completed tasks must include feedback mini retrospectives; teams should review quality signals every 10 tasks or weekly; recurring process failures must be corrected through synchronized documentation updates.
- **Archived on**: 2026-04-11 — Template development history; these are process rules built into the template framework.

## 2026-04-08: Explicit workflow declaration and Small-path output hardening

- **Context**: Recent feedback identified three recurring failure modes: (1) the canonical loop looked universally mandatory while Small path text appeared to permit implicit shortcutting, (2) users could not verify whether required discovery/triage steps were performed, and (3) Small-task minimum output expectations were not consistently enforced across source-of-truth and tool-specific files.
- **Decision**: Hardened workflow guidance across canonical and tool-specific docs by introducing a mandatory first-response compliance block, defining the loop as a canonical superset with explicit post-triage path selection, and codifying a non-skippable Small-task minimum output contract.
- **Alternatives considered**: Keep existing wording and rely on reviewer enforcement only. Rejected because non-explicit requirements are difficult to audit and are frequently interpreted as optional under time pressure.
- **Constraints introduced**: Implementation must not begin without a visible compliance block; Small-path usage must remain explicit and verifiable; workflow-related changes must synchronize `AGENTS.md`, `docs/agent-playbook.md`, `docs/operating-rules.md`, `.github/copilot-instructions.md`, `skills/demand-triage/SKILL.md`, and `CHANGELOG.md`.
- **Archived on**: 2026-04-11 — Template development history; these are structural constraints built into the template framework.

## 2026-04-07: Adaptive workflow with demand triage, context compaction, and skill quality enhancements

- **Context**: The playbook applied the same full pre-change workflow of discovery, planning, critique, approval, implementation, the test/lint/fix/repeat validation loop, and decision recording to all tasks regardless of scale. Small tasks (typo fixes, config changes, single-line bug fixes) consumed excessive tokens and time by going through planning agents, critics, risk-reviewers, full deliverable structures, and context anchors. Additionally, long tasks suffered from context bloat with no compaction mechanism, and skills lacked self-verification checklists.
- **Decision**: Introduced four changes: (1) Demand triage, (2) Context compaction, (3) Categorized memory, (4) Skill quality — conformance self-check sections added to five skills.
- **Alternatives considered**: Two-tier (Small vs. Full) without Medium: rejected. Parallel LIGHTWEIGHT workflow: rejected per critic review. Triage before codebase discovery: rejected. LOC-based threshold: rejected.
- **Constraints introduced**: Triage must run after discovery; uncertainty defaults to Medium; auth/security/schema/breaking always force non-Small; compaction summaries mandatory after each phase; conformance self-checks should be verified before completion.
- **Archived on**: 2026-04-11 — Template development history; these are structural design choices for the template framework.

## 2026-04-10: Prompt cache optimization via four-layer instruction loading order

- **Context**: No strategy existed for maximizing prompt cache hit rates across LLM providers. Without a consistent loading order, the same instruction files could be assembled in different orders across requests, causing cache misses even when the content was identical.
- **Decision**: Introduced a four-layer instruction loading order classified by change frequency. Created `skills/prompt-cache-optimization/SKILL.md` with the full skill definition.
- **Alternatives considered**: Rely solely on Claude Code's built-in cache optimization — rejected (multi-tool environment). Manual per-request cache hint annotations — rejected (not portable). No action — rejected (natural prefix overlap unreliable).
- **Constraints introduced**: Files must load in layer order (1 → 2 → 3 → 4); skills in Layer 2 in alphabetical order for the same task type; per-request content must not be injected into Layer 1 or 2; single instruction files should stay under 8 KB.
- **Archived on**: 2026-04-11 — Template development history; these are structural constraints built into the template framework.

## 2026-04-10: Cross-file deduplication, tool definition stability, and conversation memory tiering

- **Context**: A prompt efficiency audit found ~4,050 tokens of redundant content across source-of-truth files. Additionally, no guidance existed on tool/function schema stability or conversation memory tiering.
- **Decision**: Three changes: (1) Collapsed duplicate definitions to single sources with cross-references; (2) Tool definition stability guidance added to prompt-cache-optimization skill; (3) Conversation memory tiering model added to memory-and-state skill.
- **Alternatives considered**: Separate skills for tool management and conversation memory — rejected (better fit in existing skills). Keep duplicated content for "readability" — rejected (token cost outweighs convenience).
- **Constraints introduced**: Each concept must have exactly one source-of-truth file; tool schemas loaded in alphabetical order for task-type subset; conversation memory respects three-tier model with ~8,500 token budget.
- **Archived on**: 2026-04-11 — Template development history.

## 2026-04-10: Adoption-time prompt budget trimming and format template extraction

- **Context**: The full template consumes ~13,000–17,000 tokens per request. No guidance existed for adopters to self-trim. Verbose format blocks in operating-rules added ~500 tokens to every request even when not needed.
- **Decision**: (1) Moved four format blocks from operating-rules to agent-templates; (2) Added prompt budget trimming guide to adoption-guide; (3) Added `prompt-budget.yml` configuration support.
- **Alternatives considered**: Hardcode trimming into the template — rejected (different projects need different subsets). Auto-detect via static analysis — rejected (skill usage is runtime-dependent). Token budget enforcement at API layer — rejected (most providers don't support it).
- **Constraints introduced**: Format templates in operating-rules are now one-line references; full formats live in agent-templates.md. Adopting projects should create prompt-budget.yml if they want to declare a token budget.
- **Archived on**: 2026-04-11 — Template development history.
