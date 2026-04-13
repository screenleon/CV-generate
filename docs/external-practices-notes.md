# External Practices Notes

## OpenAI: use repository-local instruction files and layered documentation

OpenAI guidance around agentic coding emphasizes repository-local instruction files and structured documentation instead of relying on long, one-off chat prompts.

Useful takeaway:

- keep the root instruction file short
- put detailed rules in adjacent docs
- let the agent load only the guidance it needs

Sources:
- https://openai.com/index/harness-engineering/
- https://openai.com/pt-PT/index/introducing-codex/
- https://agents.md/

## Anthropic: use project-level subagents for specialization and context isolation

Anthropic documents project-level subagents under `.claude/agents/`, with filesystem-based Markdown definitions and separate context.

Useful takeaway:

- create small, single-purpose agents
- describe when each agent should be used
- avoid one overloaded universal agent

Sources:
- https://docs.anthropic.com/en/docs/claude-code/tutorials
- https://docs.claude.com/en/api/agent-sdk/subagents

## GitHub Copilot: keep reusable instructions and prompt files inside the repo

GitHub supports repository custom instructions and reusable prompt files.

Useful takeaway:

- store stable repo-wide guidance in `.github/copilot-instructions.md`
- store repeatable prompt assets as version-controlled files

Sources:
- https://docs.github.com/en/copilot/how-tos/custom-instructions/adding-repository-custom-instructions-for-github-copilot
- https://docs.github.com/en/copilot/tutorials/customization-library/prompt-files
