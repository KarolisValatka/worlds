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
./worlds doctor --vault ./example-vault
./worlds import --from cursor --out ./imported-vault
```

Outputs land in `./out/<target>/` by default. Pass `--install` to write into the usual home paths for that harness (see SPEC.md).

## Commands

| Command | Purpose |
|---------|---------|
| `init [dir]` | Create a starter vault |
| `validate [dir]` | Check layout + wiki frontmatter |
| `sync --vault --target [--out] [--install]` | Emit / install harness skill or rule |
| `mcp --vault <dir>` | MCP server over stdio |
| `doctor [--vault] [--fix]` | Report harness paths; optionally install missing |
| `import --from pi\|claude\|cursor [--out]` | Copy harness skills/rules into a new vault |

## MCP server

`worlds mcp --vault <dir>` speaks [Model Context Protocol](https://modelcontextprotocol.io/) over stdio (official Go SDK). Tools:

- `list_worlds`
- `read_brain` — root `BRAIN.md`
- `read_world` — world `BRAIN.md`, optional `list_wiki`
- `read_wiki` — world + page
- `write_wiki` — world + page + content/title/tags (keeps `created` / `updated`)

### Cursor / Claude mcp.json snippet

```json
{
  "mcpServers": {
    "worlds": {
      "command": "worlds",
      "args": ["mcp", "--vault", "/absolute/path/to/your-vault"]
    }
  }
}
```

Use the absolute path to the `worlds` binary if it is not on `PATH`. For Cursor, put this in `.cursor/mcp.json` (project) or your user MCP config. For Claude Desktop / Claude Code, merge into the corresponding MCP settings file.

## Doctor

```bash
./worlds doctor --vault ./example-vault
./worlds doctor --vault ./example-vault --fix
```

Best-effort checks:

- `~/.pi/agent/` and `skills/worlds-brain`
- `~/.claude/skills/` (and worlds-brain skill)
- `.cursor/rules/` in the current directory
- `AGENTS.worlds.md` / `AGENTS.md` in the current directory
- optional Ollama at `http://127.0.0.1:11434`

`--fix` runs `sync --install` for missing/relevant targets using `--vault`.

## Import

```bash
./worlds import --from pi --out ./imported-vault
./worlds import --from claude
./worlds import --from cursor --out ./from-cursor
```

Reads harness skill/rule locations (does not modify them), creates a vault at `--out` (default `./imported-vault`), and stores content under world `imported` as wiki pages with `source: import`.

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

Requires a recent Go toolchain (module declares Go 1.25; older Go 1.22+ with toolchain download also works). Until published modules cache, build from this repo:

```bash
go build -o worlds ./cmd/worlds
```

## Roadmap

- Watch mode / git hooks
- More targets (OpenCode, Aider, custom templates)
- Embedding / search helpers

## License

MIT
