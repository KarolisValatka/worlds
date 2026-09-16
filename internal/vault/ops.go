package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var worldNameRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
var pageNameRe = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ListWorlds returns sorted world directory names under the vault.
func ListWorlds(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var worlds []string
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		name := e.Name()
		if name == "out" || name == ".git" {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, name, "BRAIN.md")); err != nil {
			continue
		}
		worlds = append(worlds, name)
	}
	sort.Strings(worlds)
	return worlds, nil
}

// ReadBrain returns root BRAIN.md contents.
func ReadBrain(dir string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, "BRAIN.md"))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadWorldBrain returns <world>/BRAIN.md contents.
func ReadWorldBrain(dir, world string) (string, error) {
	if err := checkWorldName(world); err != nil {
		return "", err
	}
	b, err := os.ReadFile(filepath.Join(dir, world, "BRAIN.md"))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ListWikiPages returns wiki page basenames without .md, sorted.
func ListWikiPages(dir, world string) ([]string, error) {
	if err := checkWorldName(world); err != nil {
		return nil, err
	}
	wiki := filepath.Join(dir, world, "wiki")
	entries, err := os.ReadDir(wiki)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var pages []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		pages = append(pages, strings.TrimSuffix(e.Name(), ".md"))
	}
	sort.Strings(pages)
	return pages, nil
}

// ReadWiki returns a wiki page body (including frontmatter).
func ReadWiki(dir, world, page string) (string, error) {
	if err := checkWorldName(world); err != nil {
		return "", err
	}
	if err := checkPageName(page); err != nil {
		return "", err
	}
	b, err := os.ReadFile(wikiPath(dir, world, page))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// WriteWikiOptions controls WriteWiki.
type WriteWikiOptions struct {
	Title   string
	Content string // markdown body after frontmatter; if empty and creating, uses "# title"
	Tags    []string
	Source  string
}

// WriteWiki creates or updates a wiki page, maintaining created/updated dates.
func WriteWiki(dir, world, page string, opts WriteWikiOptions) error {
	if err := checkWorldName(world); err != nil {
		return err
	}
	if err := checkPageName(page); err != nil {
		return err
	}
	for _, sub := range []string{"wiki", "labs", "raw"} {
		if err := os.MkdirAll(filepath.Join(dir, world, sub), 0o755); err != nil {
			return err
		}
	}
	brainPath := filepath.Join(dir, world, "BRAIN.md")
	if _, err := os.Stat(brainPath); os.IsNotExist(err) {
		if err := os.WriteFile(brainPath, []byte(fmt.Sprintf("# %s index\n\n## wiki\n", world)), 0o644); err != nil {
			return err
		}
	}
	path := wikiPath(dir, world, page)
	today := time.Now().Format("2006-01-02")
	title := opts.Title
	if title == "" {
		title = page
	}
	source := opts.Source
	if source == "" {
		source = "worlds"
	}
	tags := opts.Tags
	if tags == nil {
		tags = []string{}
	}
	created := today
	body := opts.Content
	if existing, err := os.ReadFile(path); err == nil {
		fm, rest, ok := splitFrontmatter(string(existing))
		if ok {
			if c := fmValue(fm, "created"); c != "" {
				created = c
			}
			if title == page {
				if t := fmValue(fm, "title"); t != "" && opts.Title == "" {
					title = t
				}
			}
			if opts.Source == "" {
				if s := fmValue(fm, "source"); s != "" {
					source = s
				}
			}
			if opts.Tags == nil {
				if t := fmValue(fm, "tags"); t != "" {
					tags = parseTags(t)
				}
			}
			if opts.Content == "" {
				body = strings.TrimPrefix(rest, "\n")
			}
		}
	}
	if body == "" {
		body = "# " + title + "\n"
	}
	content := formatWikiPage(title, created, today, source, tags, body)
	return os.WriteFile(path, []byte(content), 0o644)
}

func wikiPath(dir, world, page string) string {
	return filepath.Join(dir, world, "wiki", page+".md")
}

func checkWorldName(world string) error {
	if !worldNameRe.MatchString(world) {
		return fmt.Errorf("invalid world name %q (use lowercase kebab-case)", world)
	}
	return nil
}

func checkPageName(page string) error {
	page = strings.TrimSuffix(page, ".md")
	if !pageNameRe.MatchString(page) {
		return fmt.Errorf("invalid page name %q (use lowercase kebab-case)", page)
	}
	return nil
}

func formatWikiPage(title, created, updated, source string, tags []string, body string) string {
	tagStr := "[" + strings.Join(tags, ", ") + "]"
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("title: " + title + "\n")
	b.WriteString("created: " + created + "\n")
	b.WriteString("updated: " + updated + "\n")
	b.WriteString("source: " + source + "\n")
	b.WriteString("tags: " + tagStr + "\n")
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimLeft(body, "\n"))
	if !strings.HasSuffix(body, "\n") {
		b.WriteByte('\n')
	}
	return b.String()
}

func splitFrontmatter(body string) (fm, rest string, ok bool) {
	if !strings.HasPrefix(body, "---\n") {
		return "", body, false
	}
	end := strings.Index(body[4:], "\n---")
	if end < 0 {
		return "", body, false
	}
	fm = body[4 : 4+end]
	rest = body[4+end+4:] // after \n---
	if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	}
	return fm, rest, true
}

func fmValue(fm, key string) string {
	prefix := key + ":"
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func parseTags(raw string) []string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.Trim(p, `"'`)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// EnsureWorld creates world dirs and BRAIN.md if missing.
func EnsureWorld(dir, world string) error {
	if err := checkWorldName(world); err != nil {
		return err
	}
	for _, sub := range []string{"wiki", "labs", "raw"} {
		if err := os.MkdirAll(filepath.Join(dir, world, sub), 0o755); err != nil {
			return err
		}
	}
	brain := filepath.Join(dir, world, "BRAIN.md")
	if _, err := os.Stat(brain); os.IsNotExist(err) {
		return os.WriteFile(brain, []byte(fmt.Sprintf("# %s index\n\n## wiki\n", world)), 0o644)
	}
	return nil
}

// EnsureRootBrain creates root BRAIN.md if missing.
func EnsureRootBrain(dir, content string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "BRAIN.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if content == "" {
		content = "# Brain index\n\nImported vault.\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// Slugify turns a filename/title into a kebab-case page name.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSuffix(s, ".md")
	s = strings.TrimSuffix(s, ".mdc")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-'
		if !ok {
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
			continue
		}
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteRune(r)
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "page"
	}
	return out
}
