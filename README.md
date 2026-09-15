# worlds

One markdown vault. Every coding agent gets the same brain.

Pi, Claude Code, Cursor, and Codex each reinvent skills and project memory in different folders. **worlds** is a tiny open format plus a CLI that syncs one vault into those harnesses.

Status: early / experimental.

## Why

Local and cloud coding agents are multiplying. Context is the bottleneck. Copy-pasting half-broken skill packs per tool is the current state of the art. Worlds fixes the pipe, not the model.

## Quickstart

```bash
go build -o worlds ./cmd/worlds
./worlds init ./my-vault
./worlds validate ./example-vault
./worlds sync --vault ./example-vault --target pi
./worlds sync --vault ./example-vault --target claude
./worlds sync --vault ./example-vault --target cursor
./worlds sync --vault ./example-vault --target codex
```

Outputs land in `./out/<target>/` by default. Pass `--install` to write into the usual home paths for that harness (see SPEC.md).

## Vault shape

```
BRAIN.md                 # hot index
<world>/BRAIN.md         # world index
<world>/wiki/*.md        # durable pages + frontmatter
<world>/labs/            # runnable what|command notes
<world>/raw/             # immutable drops YYYY-MM-DD-what.md
```

Full schema: [SPEC.md](SPEC.md). Demo vault: [example-vault/](example-vault/).

## Install

```bash
go install github.com/KarolisValatka/worlds/cmd/worlds@latest
```

(Until published, build from this repo.)

## Roadmap

- MCP server: `read_brain`, `read_world`, `write_wiki`
- Watch mode / git hooks
- More targets (OpenCode, Aider, custom templates)

## License

MIT
