# worlds vault specification

## Goals

- One durable markdown vault shared across agent harnesses
- Progressive disclosure: hot index first, world pages on demand
- Safe defaults: no secrets, immutable raw drops, explicit writes

## Layout

```
<vault>/
  BRAIN.md
  <world>/
    BRAIN.md
    wiki/
    labs/
    raw/
```

A **world** is a subject of work (not an arbitrary disk folder). World directory names are lowercase kebab-case: `learn`, `studio`, `hackme`.

## Root BRAIN.md

Plain markdown. Lists worlds and hot pages. Keep it short (under ~45 lines recommended).

## World BRAIN.md

Index for that world. One line per hot page, preferably wikilink style `[[page-name]]`.

## Wiki pages

Markdown files under `<world>/wiki/`. Required YAML frontmatter:

```yaml
---
title: string
created: YYYY-MM-DD
updated: YYYY-MM-DD
source: string
tags: [string]
---
```

Rules:
- One page per thing; update in place
- Date facts that go stale
- Wikilinks `[[page-name]]` encouraged; missing targets are fine
- Never store passwords, tokens, or other people's personal data

## Labs

Optional markdown under `<world>/labs/`. Prefer two-column what|command tables.

## Raw

Immutable drops under `<world>/raw/`, named `YYYY-MM-DD-<what>.md`. Never edit; only add.

## Sync targets

| target | Default out dir | `--install` destination |
|--------|-----------------|-------------------------|
| pi | `out/pi/` | `~/.pi/agent/skills/worlds-brain/SKILL.md` (+ optional extension note) |
| claude | `out/claude/` | `~/.claude/skills/worlds-brain/SKILL.md` |
| cursor | `out/cursor/` | `<cwd>/.cursor/rules/worlds-brain.mdc` if in a project, else `out/cursor/` |
| codex | `out/codex/` | `./AGENTS.worlds.md` fragment in cwd or `out/codex/AGENTS.md` |

Generated skill content must instruct the agent to:
1. Read `<vault>/BRAIN.md` at session start when context matters
2. Open only the relevant world BRAIN.md
3. Check `updated` before treating wiki facts as current
4. Write back with frontmatter; do not commit unless the human asks

## Non-goals (MVP)

- Real-time multiplayer editing
- Embedding search
- Hosted sync SaaS
