package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func Init(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	rootBrain := filepath.Join(dir, "BRAIN.md")
	if _, err := os.Stat(rootBrain); err == nil {
		return fmt.Errorf("%s already exists", rootBrain)
	}
	if err := os.WriteFile(rootBrain, []byte("# Brain index\n\n## learn\nIndex: [[learn/BRAIN]]\n"), 0o644); err != nil {
		return err
	}
	world := filepath.Join(dir, "learn")
	for _, sub := range []string{"wiki", "labs", "raw"} {
		if err := os.MkdirAll(filepath.Join(world, sub), 0o755); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(world, "BRAIN.md"), []byte("# learn index\n\n## wiki\n- [[hello]] - sample page\n"), 0o644); err != nil {
		return err
	}
	page := "---\ntitle: hello\ncreated: 2026-09-15\nupdated: 2026-09-15\nsource: worlds init\ntags: [learn]\n---\n\n# hello\n\nSample wiki page.\n"
	return os.WriteFile(filepath.Join(world, "wiki", "hello.md"), []byte(page), 0o644)
}

func Validate(dir string) error {
	root := filepath.Join(dir, "BRAIN.md")
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("missing BRAIN.md: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	worlds := 0
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		name := e.Name()
		if name == "out" || name == ".git" {
			continue
		}
		wdir := filepath.Join(dir, name)
		if _, err := os.Stat(filepath.Join(wdir, "BRAIN.md")); err != nil {
			return fmt.Errorf("world %s missing BRAIN.md", name)
		}
		for _, sub := range []string{"wiki", "labs", "raw"} {
			if err := os.MkdirAll(filepath.Join(wdir, sub), 0o755); err != nil {
				return err
			}
		}
		wiki := filepath.Join(wdir, "wiki")
		pages, err := os.ReadDir(wiki)
		if err != nil {
			return err
		}
		for _, p := range pages {
			if p.IsDir() || !strings.HasSuffix(p.Name(), ".md") {
				continue
			}
			body, err := os.ReadFile(filepath.Join(wiki, p.Name()))
			if err != nil {
				return err
			}
			if err := checkFrontmatter(string(body), filepath.Join(name, "wiki", p.Name())); err != nil {
				return err
			}
		}
		worlds++
	}
	if worlds == 0 {
		return fmt.Errorf("no world folders found under %s", dir)
	}
	fmt.Printf("ok: %s (%d worlds)\n", dir, worlds)
	return nil
}

func checkFrontmatter(body, path string) error {
	if !strings.HasPrefix(body, "---\n") {
		return fmt.Errorf("%s: missing YAML frontmatter", path)
	}
	end := strings.Index(body[4:], "\n---")
	if end < 0 {
		return fmt.Errorf("%s: unterminated frontmatter", path)
	}
	fm := body[4 : 4+end]
	need := []string{"title:", "created:", "updated:", "source:", "tags:"}
	for _, n := range need {
		if !strings.Contains(fm, n) {
			return fmt.Errorf("%s: frontmatter missing %s", path, strings.TrimSuffix(n, ":"))
		}
	}
	for _, key := range []string{"created", "updated"} {
		for _, line := range strings.Split(fm, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, key+":") {
				val := strings.TrimSpace(strings.TrimPrefix(line, key+":"))
				if !dateRe.MatchString(val) {
					return fmt.Errorf("%s: %s must be YYYY-MM-DD", path, key)
				}
			}
		}
	}
	return nil
}

// Abs resolves vault path.
func Abs(dir string) (string, error) {
	return filepath.Abs(dir)
}
