package imprt

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/KarolisValatka/worlds/internal/vault"
)

// Options configures Import.
type Options struct {
	From string // pi|claude|cursor
	Out  string
	Home string
	CWD  string
}

// Run reads harness skill/rule files and writes an imported vault. Never modifies sources.
func Run(opts Options) error {
	home := opts.Home
	if home == "" {
		var err error
		home, err = os.UserHomeDir()
		if err != nil {
			return err
		}
	}
	cwd := opts.CWD
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	out := opts.Out
	if out == "" {
		out = "./imported-vault"
	}
	sources, err := collectSources(opts.From, home, cwd)
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return fmt.Errorf("no skill/rule files found for --from %s", opts.From)
	}
	if err := vault.EnsureRootBrain(out, "# Brain index\n\n## imported\nPages imported from harness skills/rules.\nIndex: [[imported/BRAIN]]\n"); err != nil {
		return err
	}
	if err := vault.EnsureWorld(out, "imported"); err != nil {
		return err
	}
	var titles []string
	for _, src := range sources {
		body, err := os.ReadFile(src.Path)
		if err != nil {
			return err
		}
		page := vault.Slugify(src.Name)
		title := src.Name
		if err := vault.WriteWiki(out, "imported", page, vault.WriteWikiOptions{
			Title:   title,
			Content: string(body),
			Tags:    []string{"import", opts.From},
			Source:  "import",
		}); err != nil {
			return err
		}
		titles = append(titles, page)
		fmt.Printf("imported %s -> imported/wiki/%s.md\n", src.Path, page)
	}
	var index strings.Builder
	index.WriteString("# imported index\n\n## wiki\n")
	for _, t := range titles {
		index.WriteString(fmt.Sprintf("- [[%s]]\n", t))
	}
	if err := os.WriteFile(filepath.Join(out, "imported", "BRAIN.md"), []byte(index.String()), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote vault %s (%d pages)\n", out, len(titles))
	return nil
}

type sourceFile struct {
	Path string
	Name string
}

func collectSources(from, home, cwd string) ([]sourceFile, error) {
	switch from {
	case "pi":
		root := filepath.Join(home, ".pi", "agent", "skills")
		return listMarkdownTree(root, []string{".md"})
	case "claude":
		root := filepath.Join(home, ".claude", "skills")
		return listMarkdownTree(root, []string{".md"})
	case "cursor":
		root := filepath.Join(cwd, ".cursor", "rules")
		return listMarkdownTree(root, []string{".md", ".mdc"})
	default:
		return nil, fmt.Errorf("unknown --from %q (pi|claude|cursor)", from)
	}
}

func listMarkdownTree(root string, exts []string) ([]sourceFile, error) {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, nil
	}
	var out []sourceFile
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		lower := strings.ToLower(name)
		ok := false
		for _, ext := range exts {
			if strings.HasSuffix(lower, ext) {
				ok = true
				break
			}
		}
		if !ok {
			return nil
		}
		base := name
		for _, ext := range exts {
			if strings.HasSuffix(strings.ToLower(base), ext) {
				base = base[:len(base)-len(ext)]
				break
			}
		}
		// Prefer parent skill folder name when file is SKILL.md
		if strings.EqualFold(name, "SKILL.md") {
			base = filepath.Base(filepath.Dir(path))
		}
		out = append(out, sourceFile{Path: path, Name: base})
		return nil
	})
	return out, err
}
