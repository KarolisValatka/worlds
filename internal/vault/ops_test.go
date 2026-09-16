package vault

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListWorldsAndWikiOps(t *testing.T) {
	dir := t.TempDir()
	if err := Init(filepath.Join(dir, "v")); err != nil {
		t.Fatal(err)
	}
	v := filepath.Join(dir, "v")
	worlds, err := ListWorlds(v)
	if err != nil {
		t.Fatal(err)
	}
	if len(worlds) != 1 || worlds[0] != "learn" {
		t.Fatalf("worlds=%v", worlds)
	}
	brain, err := ReadBrain(v)
	if err != nil || !strings.Contains(brain, "Brain index") {
		t.Fatalf("brain=%q err=%v", brain, err)
	}
	wb, err := ReadWorldBrain(v, "learn")
	if err != nil || !strings.Contains(wb, "learn") {
		t.Fatalf("world brain=%q err=%v", wb, err)
	}
	pages, err := ListWikiPages(v, "learn")
	if err != nil || len(pages) != 1 || pages[0] != "hello" {
		t.Fatalf("pages=%v err=%v", pages, err)
	}
	if err := WriteWiki(v, "learn", "hello", WriteWikiOptions{
		Content: "# hello\n\nUpdated.\n",
		Tags:    []string{"learn", "test"},
		Source:  "test",
	}); err != nil {
		t.Fatal(err)
	}
	body, err := ReadWiki(v, "learn", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(body, "updated:") || !strings.Contains(body, "created:") {
		t.Fatalf("missing dates: %s", body)
	}
	if !strings.Contains(body, "Updated.") {
		t.Fatalf("body not updated: %s", body)
	}
	// created should be preserved from init
	if !strings.Contains(body, "created: 2026-09-15") && !strings.Contains(body, "created:") {
		t.Fatalf("created missing: %s", body)
	}
}

func TestSlugify(t *testing.T) {
	if got := Slugify("SKILL.md"); got != "skill" {
		t.Fatalf("got %q", got)
	}
	if got := Slugify("Worlds Brain.mdc"); got != "worlds-brain" {
		t.Fatalf("got %q", got)
	}
}

func TestWriteWikiCreatesWorld(t *testing.T) {
	dir := t.TempDir()
	if err := EnsureRootBrain(dir, "# Brain\n"); err != nil {
		t.Fatal(err)
	}
	if err := WriteWiki(dir, "imported", "from-pi", WriteWikiOptions{
		Title:   "from-pi",
		Content: "hello",
		Source:  "import",
		Tags:    []string{"import"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "imported", "wiki", "from-pi.md")); err != nil {
		t.Fatal(err)
	}
}
