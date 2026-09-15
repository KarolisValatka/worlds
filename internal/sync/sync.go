package sync

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KarolisValatka/worlds/internal/vault"
)

func Sync(vaultDir, target, outDir string, install bool) error {
	abs, err := vault.Abs(vaultDir)
	if err != nil {
		return err
	}
	if err := vault.Validate(abs); err != nil {
		return err
	}
	body := renderSkill(abs)
	switch target {
	case "pi", "claude", "cursor", "codex":
	default:
		return fmt.Errorf("unknown target %q", target)
	}
	if install {
		return installTarget(target, body, abs)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	name := map[string]string{
		"pi":     "SKILL.md",
		"claude": "SKILL.md",
		"cursor": "worlds-brain.mdc",
		"codex":  "AGENTS.md",
	}[target]
	path := filepath.Join(outDir, name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", path)
	return nil
}

func renderSkill(vaultAbs string) string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: worlds-brain\n")
	b.WriteString("description: Shared durable brain vault for cross-session agent context. Use when you need prior decisions, project status, or to save something that should survive across sessions and harnesses.\n")
	b.WriteString("---\n\n")
	b.WriteString("# worlds brain\n\n")
	b.WriteString("Vault path: `")
	b.WriteString(vaultAbs)
	b.WriteString("`\n\n")
	b.WriteString("## Session start\n")
	b.WriteString("1. Read `BRAIN.md` at the vault root when context matters.\n")
	b.WriteString("2. Open only the relevant world `BRAIN.md`.\n")
	b.WriteString("3. Check wiki page `updated` dates before treating facts as current.\n\n")
	b.WriteString("## Layout\n")
	b.WriteString("- `<world>/wiki/` durable pages with frontmatter\n")
	b.WriteString("- `<world>/labs/` runnable notes\n")
	b.WriteString("- `<world>/raw/` immutable drops `YYYY-MM-DD-what.md`\n\n")
	b.WriteString("## Write rules\n")
	b.WriteString("- One wiki page per thing; update in place\n")
	b.WriteString("- Never store secrets\n")
	b.WriteString("- Commit or push only when the human asks\n")
	return b.String()
}

func installTarget(target, body, vaultAbs string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	var path string
	switch target {
	case "pi":
		path = filepath.Join(home, ".pi", "agent", "skills", "worlds-brain", "SKILL.md")
	case "claude":
		path = filepath.Join(home, ".claude", "skills", "worlds-brain", "SKILL.md")
	case "cursor":
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, ".cursor", "rules", "worlds-brain.mdc")
	case "codex":
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, "AGENTS.worlds.md")
	}
	_ = vaultAbs
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return err
	}
	fmt.Printf("installed %s\n", path)
	return nil
}
