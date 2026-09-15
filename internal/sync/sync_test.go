package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KarolisValatka/worlds/internal/vault"
)

func TestSyncPi(t *testing.T) {
	dir := t.TempDir()
	v := filepath.Join(dir, "v")
	if err := vault.Init(v); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "out")
	if err := Sync(v, "pi", out, false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(out, "SKILL.md")); err != nil {
		t.Fatal(err)
	}
}
