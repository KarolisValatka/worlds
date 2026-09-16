package doctor

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KarolisValatka/worlds/internal/vault"
)

func TestDoctorReport(t *testing.T) {
	home := t.TempDir()
	cwd := t.TempDir()
	v := filepath.Join(t.TempDir(), "v")
	if err := vault.Init(v); err != nil {
		t.Fatal(err)
	}
	// partial install
	if err := os.MkdirAll(filepath.Join(home, ".pi", "agent"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(cwd, ".cursor", "rules"), 0o755); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	err := Run(Options{
		VaultDir: v,
		Home:     home,
		CWD:      cwd,
		Fix:      false,
		Stdout:   &buf,
		HTTPGet: func(url string) (int, error) {
			return 200, nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "OK") || !strings.Contains(out, "missing") {
		t.Fatalf("expected mixed report, got:\n%s", out)
	}
	if !strings.Contains(out, "ollama") || !strings.Contains(out, "HTTP 200") {
		t.Fatalf("ollama check missing:\n%s", out)
	}
	if !strings.Contains(out, "pi agent dir") {
		t.Fatalf("pi check missing:\n%s", out)
	}
}

func TestDoctorFixInstalls(t *testing.T) {
	home := t.TempDir()
	cwd := t.TempDir()
	v := filepath.Join(t.TempDir(), "v")
	if err := vault.Init(v); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	// chdir so cursor/codex install under cwd
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	err = Run(Options{
		VaultDir: v,
		Home:     home,
		CWD:      cwd,
		Fix:      true,
		Stdout:   &buf,
		HTTPGet: func(url string) (int, error) {
			return 0, os.ErrNotExist
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".pi", "agent", "skills", "worlds-brain", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude", "skills", "worlds-brain", "SKILL.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(cwd, ".cursor", "rules", "worlds-brain.mdc")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(cwd, "AGENTS.worlds.md")); err != nil {
		t.Fatal(err)
	}
}
