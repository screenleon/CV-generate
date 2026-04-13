---
agent: agent
description: Review only the current git-status changes inside a specific folder, with findings first.
---

Review the current uncommitted changes only within this folder:

Target folder: ${input:target_folder:Enter the folder path to review, relative to the repository root}

Instructions:

1. Determine which files under the target folder are currently changed according to `git status` or equivalent workspace state.
2. If there are no changed files under that folder, say so explicitly and stop.
3. Review only those changed files. Ignore changes outside the target folder.
4. Use a code review mindset:
   - prioritize bugs
   - behavioral regressions
   - security or permission issues
   - data consistency risks
   - missing validation
   - missing tests
5. Do not give a generic summary first.
6. Present findings first, ordered by severity, with file references when possible.
7. After findings, include:
   - open questions or assumptions
   - a short change summary
8. If there are no findings, say that explicitly and mention any residual risk or testing gaps.

Follow the repository instructions in `AGENTS.md` and the safety/validation guidance in `docs/operating-rules.md`.
