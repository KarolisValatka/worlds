# worlds

One markdown vault. Point every coding agent at it.

You know the drill: Claude Code has skills in one place, Pi in another, Cursor rules somewhere else, Codex wants `AGENTS.md`. You paste the same "here's how my notes work" speech into a new chat, the model nods, then invents a second note format by Thursday.

**worlds** is a small Go CLI plus a boring markdown layout. Sync the same vault into those harnesses. Optional MCP if you want tools instead of shell-grepping.

Early. Experimental. Useful to me. Maybe useful to you.

## Who this is for

You run more than one coding agent (or you will), and you're tired of re-explaining project memory.

If you only ever use one tool and love pasting a paragraph, you don't need this yet. Just tell the agent. Seriously.

## Install

```bash
git clone https://github.com/KarolisValatka/worlds
cd worlds
go build -o worlds ./cmd/worlds
```

Or:

```bash
go install github.com/KarolisValatka/worlds/cmd/worlds@latest
```

Needs a recent Go (module says 1.25; toolchain download usually handles it).

## 30-second tour

```bash
./worlds validate ./example-vault
./worlds sync --vault ./example-vault --target pi
./worlds doctor --vault ./example-vault
```

`sync` writes under `./out/<target>/` by default. Add `--install` to drop files into the usual home paths (Pi skills, Claude skills, Cursor rules, Codex fragment). Details in [SPEC.md](SPEC.md).

Demo vault lives in [example-vault/](example-vault/). It's fake. Your real notes stay yours.

## The vault

```
BRAIN.md              # hot index (keep it short)
learn/BRAIN.md        # one "world" = one subject
learn/wiki/*.md       # durable pages + frontmatter
learn/labs/           # runnable notes
learn/raw/            # drops you never edit
```

Worlds are subjects (`work`, `studio`, `learn`), not random folders. Wiki pages need `title`, `created`, `updated`, `source`, `tags`. No passwords in the vault. Commit when *you* say so.

Full rules: [SPEC.md](SPEC.md).

## Commands

```text
worlds init [dir]
worlds validate [dir]
worlds sync   --vault <dir> --target pi|claude|cursor|codex [--out dir] [--install]
worlds mcp    --vault <dir>
worlds doctor [--vault dir] [--fix]
worlds import --from pi|claude|cursor [--out dir]
```

**doctor** — looks for Pi / Claude / Cursor / Codex wiring (and optionally Ollama). `--fix` runs `sync --install` for what's missing.

**import** — copy existing skills/rules into a new vault under world `imported`. Does not touch the source files.

**mcp** — stdio MCP server (`list_worlds`, `read_brain`, `read_world`, `read_wiki`, `write_wiki`).

Cursor / Claude style config:

```json
{
  "mcpServers": {
    "worlds": {
      "command": "/absolute/path/to/worlds",
      "args": ["mcp", "--vault", "/absolute/path/to/your-vault"]
    }
  }
}
```

## Why not just tell the agent?

You can. For one session, one harness, that works.

worlds is for the second session, the second tool, and the teammate who wasn't in the chat. The vault is the memory. The CLI stamps a pointer + rules into whatever harness you're using today so you stop retyping the speech.

It does not make the model smarter. It keeps the notes from rotting in six incompatible folders.

## Status / non-goals

Shipped: init, validate, sync, mcp, doctor, import.

Not trying to be Obsidian, a RAG product, or another chat UI. Watch mode and more targets can wait until someone actually wants them.

## License

MIT. Copyright (c) 2026 Karolis Valatka.
