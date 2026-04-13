---
name: memory-and-state
description: Use to maintain persistent context across agent sessions — decision logs, architecture memory, and constraint tracking.
---

# Memory and State

Use this skill to prevent agents from forgetting prior decisions and repeating mistakes.

## Three memory layers

### 1. Decision log (`DECISIONS.md`)

A version-controlled file at the repo root. Records architectural and behavioral decisions.

Format:

```markdown
## YYYY-MM-DD: [Decision title]
- **Context**: Why this decision was needed
- **Decision**: What was decided
- **Alternatives considered**: What was rejected and why
- **Constraints introduced**: What future work must respect
```

Rules:
- Read `DECISIONS.md` at the start of every planning or implementation task. This is not optional.
- Append to it when making a decision that future work depends on.
- Never silently contradict an existing decision — either follow it or propose a reversal with user approval.
- When appending, verify the new entry does not conflict with existing entries.

### 2. Architecture memory (`ARCHITECTURE.md`)

A version-controlled file describing the codebase structure.

Should contain:
- Module map (directory → purpose)
- Key interfaces and their contracts
- Data flow overview
- External service dependencies
- Known technical debt

Rules:
- Read before working on unfamiliar modules.
- Update when structural changes are made (new modules, moved files, changed boundaries).

### 3. Session-scoped working memory

For in-progress tasks that span multiple agent interactions:

- Use session notes or scratch files to track:
  - What has been done so far
  - What remains
  - Errors encountered and how they were resolved
  - Open questions

This prevents re-discovery and reduces repeated mistakes within a session.

## Conversation memory tiering

Within a single conversation or agent session, manage history in three tiers to prevent context window exhaustion while preserving relevant information.

### Tier 1 — Short-term (raw recent turns)

Keep the most recent **3–5 turns** of conversation as raw, unmodified content. This preserves full fidelity for the active working context.

- Adjust the window size based on turn complexity: if turns are long (>500 tokens each), reduce to 2–3 turns.
- This tier occupies Layer 4 (volatile context) in the prompt cache loading order.

### Tier 2 — Mid-term (compressed summaries)

When turns age out of the short-term window, compress them into structured summaries:

```markdown
## Conversation summary (turns 1–N)
- **Decisions made**: [list]
- **Files changed**: [list]
- **Errors encountered and resolved**: [list]
- **Open questions**: [list]
- **Current plan state**: [brief]
```

Rules:
- Produce a summary when the short-term window shifts (i.e., every time a turn exits the window).
- For batch efficiency, summarize in groups of 3–5 turns rather than one at a time.
- Store summaries in session memory. The most recent summary replaces (not appends to) older summaries.
- This tier integrates with the existing context compaction protocol in `docs/operating-rules.md`.

### Tier 3 — Long-term (persistent retrieval)

For knowledge that persists beyond a single session, use the existing persistent stores:

| Store | Content | Retrieval |
|-------|---------|-----------|
| `DECISIONS.md` | Architectural and behavioral decisions | Read at task start |
| `ARCHITECTURE.md` | Module map, interfaces, data flow | Read when working on unfamiliar modules |
| Repo memory files | Reusable patterns, component-level notes | Search by module name or task type keywords |
| `DECISIONS_ARCHIVE.md` | Inactive past decisions | Search only for legacy module work |

**Advanced (optional):** For teams with vector database or embedding infrastructure, long-term memory can be augmented with semantic retrieval (RAG):

- Index `DECISIONS.md`, `ARCHITECTURE.md`, session summaries, and past task completion summaries as embeddings.
- At task start, retrieve the top-K most relevant entries by query similarity instead of reading full files.
- This is an optimization for large codebases where file-based reads exceed practical token budgets. It is not required for the playbook to function.

### Token budget guideline

| Tier | Target budget | Enforcement |
|------|--------------|-------------|
| Short-term (raw turns) | ≤ 4,000 tokens | Trim oldest turn when exceeded |
| Mid-term (summary) | ≤ 1,500 tokens | Regenerate summary with tighter compression |
| Long-term (persistent reads) | ≤ 3,000 tokens per task | Use selective read strategy (see below) |
| **Total conversation memory** | **≤ 8,500 tokens** | Roughly 6–7% of a 128K context window |

### Interaction with prompt cache optimization

Conversation memory is entirely within Layer 4 (volatile context). It does not affect the cached prefix in Layers 1–3. However, keeping conversation memory compact:
- Leaves more context window for actual code and tool outputs.
- Reduces per-request cost even when cache misses occur.

## Context anchor protocol

For any task spanning more than one step or more than one file, maintain a context anchor using the canonical template in `docs/agent-templates.md` → Context anchor template.

Do not duplicate or redefine the template here; treat `docs/agent-templates.md` as the single source of truth for the anchor format.

Update this anchor before each major step. This prevents drift by forcing the agent to re-read the plan and current state.

## Contradiction detection

Before making any decision, check `DECISIONS.md` for conflicts:

1. Read the full decision log
2. Compare each existing entry against the proposed change
3. If a conflict exists, state: the existing decision (date + title), the proposed change, why they conflict, and options (follow existing or reverse with justification)

STOP and wait for user decision. Do not resolve contradictions autonomously.

## When to write memory

| Event | Action |
|-------|--------|
| Architectural decision made | Append to `DECISIONS.md` |
| New module or structural change | Update `ARCHITECTURE.md` |
| Constraint discovered during work | Add to `Project-specific constraints` in `docs/operating-rules.md` |
| Error pattern found | Note in session memory to avoid repeating |
| Task partially complete | Write progress to session notes |
| Technology or library introduced | Append to `DECISIONS.md` |
| Schema or contract changed | Append to `DECISIONS.md` |
| Tradeoff made | Append to `DECISIONS.md` |

## When to read memory

| Event | Action |
|-------|--------|
| Starting any implementation | Read `DECISIONS.md` and `ARCHITECTURE.md` |
| Encountering unfamiliar module | Read `ARCHITECTURE.md` |
| Making a decision | Check `DECISIONS.md` for prior related decisions |
| Resuming interrupted work | Read session notes |
| Starting a long task | Produce a context anchor |
| Starting a Small task | Query categorized memory for similar past patterns (see below) |

## Categorized memory structure

| Category | Content | Primary store | Query when |
|----------|---------|---------------|------------|
| Project-level | Architectural decisions, global conventions, tech choices | `DECISIONS.md`, `ARCHITECTURE.md` | Starting any task, making architectural choices |
| Component-level | Per-module patterns, quirks, module-specific constraints | Module READMEs, session/repo memory files | Working on a specific or unfamiliar module (search by module name/path) |
| Change-pattern | Recurring fix patterns, validated approaches for similar tasks | Session/repo memory files | Starting a Small task (search by task-type keywords, e.g., "validation", "config update") |

**Small task retrieval**: Before implementing a Small task, search session/repo memory for the affected module or task-type keywords. If a matching pattern exists, follow it. If not, proceed normally and capture the pattern in the task completion summary.

## Memory lifecycle management

Persistent memory files grow over time. Without active lifecycle management, they consume excessive tokens on every read and eventually exceed context windows.

### Decision archive (cold storage)

#### When to archive

- `DECISIONS.md` exceeds **50 entries** or **30 KB**, OR
- Any memory health indicator (see below) reaches the "needs attention" threshold, OR
- During periodic maintenance review (quarterly for low-volume projects; more frequently when thresholds are hit)

#### Archive procedure

1. Review each entry in `DECISIONS.md`
2. For each entry, check: **are the constraints introduced still actively referenced by current code or docs?**
3. If the constraints are **no longer active** (the code has moved on, the pattern was replaced, etc.), move the entry to `DECISIONS_ARCHIVE.md`
4. If the constraints are **still enforced**, keep the entry in `DECISIONS.md`

Never archive based on date alone — a 2-year-old decision with active constraints stays in `DECISIONS.md`.

#### Archive file format

`DECISIONS_ARCHIVE.md` uses the same entry format as `DECISIONS.md`, with one addition: append `- **Archived on**: YYYY-MM-DD — [reason, e.g., "replaced by decision X"]` to each entry.

#### Safety checks before archiving

- [ ] Archived entry's constraints are verifiably no longer enforced
- [ ] No current code depends on the archived pattern
- [ ] Archive entry includes the reason
- [ ] `DECISIONS.md` retains all entries with active constraints

### Selective read strategy

Agents should not read the full archive on every task. Use a tiered approach. If `DECISIONS_ARCHIVE.md` does not exist yet, treat archive searches as returning no matches and skip any archive-read step that would otherwise require the file.

| Situation | What to read |
|-----------|-------------|
| Normal task | `DECISIONS.md` only (active constraints) |
| Task involves legacy module or old migration | `DECISIONS.md` + search `DECISIONS_ARCHIVE.md` for module name, if the archive file exists |
| Contradiction detection finds no match in active | Search `DECISIONS_ARCHIVE.md` before concluding "no prior decision", if the archive file exists; otherwise treat as no archived match |
| Periodic maintenance review | Read both files in full if `DECISIONS_ARCHIVE.md` exists; otherwise read `DECISIONS.md` only |

### Session memory hygiene

Session-scoped memory (scratch notes, in-progress tracking) should not accumulate without bound.

#### Promotion rule

After task completion, promote session memory to repo-level memory **only if**:

- The pattern was reused 2+ times in different tasks, OR
- The feedback loop mini retrospective flagged it as "most useful"

All other session notes are disposable after the task completion summary is produced.

#### Cleanup cadence

- At the end of each task: review session notes, promote or discard
- At the end of each week (or every 10 tasks): purge session memory that was not promoted

### Memory health indicators

Track during feedback loop quality signal reviews:

| Indicator | Healthy | Needs attention |
|-----------|---------|-----------------|
| `DECISIONS.md` entry count | ≤ 50 | > 50 without recent archive |
| `DECISIONS_ARCHIVE.md` exists | Yes, once any archiving has occurred | No, with 50+ active decisions |
| Session memory files | ≤ 5 active | > 10 without cleanup |
| Stale constraint references | 0 | Any archived constraint still referenced in code |

## Context compaction protocol

Long tasks cause context to grow, increasing cost and reducing model accuracy. Use compaction to prevent this.

### When to compact

- After completing each **phase** of a multi-phase task (e.g., after planning, after each implementation group)
- When the conversation has exceeded **10+ back-and-forth exchanges** on the same task
- Before handing off to a different agent role (produce a handoff artifact — this is also a form of compaction)
- When you notice yourself re-reading earlier messages to remember what was decided

### Compaction procedure

1. **Produce a progress summary** capturing:
   - What has been completed
   - Key decisions made (with `DECISIONS.md` references if applicable)
   - Current state of the work
   - What remains to be done
   - Any errors encountered and how they were resolved

2. **Store the summary** in session memory or working notes

3. **Continue from the summary**, not from the full conversation history

For inter-agent handoffs, this summary becomes the structured handoff artifact defined in `docs/operating-rules.md`.

### Post-task summary

After completing any task (regardless of scale), produce a task completion summary using the template in `docs/agent-templates.md` → Task completion summary. That template is the single source of truth for the summary format.

For Small tasks: if the summary includes a reusable pattern, store it in session or repo memory for future reuse.
For Medium/Large tasks: the summary also feeds into documentation sync checks.

## Use this skill when

- Starting work on a repository for the first time
- Making decisions that affect future work
- Resuming work after a break
- Noticing that an agent is repeating a previously resolved mistake
- Working on a task that spans more than one step or file
