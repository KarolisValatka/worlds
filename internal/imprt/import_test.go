package imprt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KarolisValatka/worlds/internal/vault"
)

func TestImportFromPi(t *testing.T) {
	home := t.TempDir()
	skillDir := filepath.Join(home, ".pi", "agent", "skills", "my-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# my skill\n\nDo things.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "imported-vault")
	if err := Run(Options{From: "pi", Out: out, Home: home, CWD: t.TempDir()}); err != nil {
		t.Fatal(err)
	}
	page := filepath.Join(out, "imported", "wiki", "my-skill.md")
	body, err := os.ReadFile(page)
	if err != nil {
		t.Fatal(err)
	}
	s := string(body)
	if !strings.Contains(s, "source: import") {
		t.Fatalf("missing source: %s", s)
	}
	if !strings.Contains(s, "Do things.") {
		t.Fatalf("missing content: %s", s)
	}
	// source file untouched
	orig, _ := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if string(orig) != "# my skill\n\nDo things.\n" {
		t.Fatal("source modified")
	}
	if err := vault.Validate(out); err != nil {
		t.Fatal(err)
	}
}

func TestImportFromCursor(t *testing.T) {
	cwd := t.TempDir()
	rules := filepath.Join(cwd, ".cursor", "rules")
	if err := os.MkdirAll(rules, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rules, "team.mdc"), []byte("always use tests\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(t.TempDir(), "v")
	if err := Run(Options{From: "cursor", Out: out, Home: t.TempDir(), CWD: cwd}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "imported", "wiki", "team.md")); err != nil {
		t.Fatal(err)
	}
}
