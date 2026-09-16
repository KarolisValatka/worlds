package doctor

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KarolisValatka/worlds/internal/sync"
	"github.com/KarolisValatka/worlds/internal/vault"
)

// Check describes one doctor probe.
type Check struct {
	Name   string
	Path   string
	OK     bool
	Detail string
	Target string // sync target if --fix should install this
}

// Options configures Doctor.
type Options struct {
	VaultDir string
	CWD      string
	Home     string
	Fix      bool
	Stdout   io.Writer
	HTTPGet  func(url string) (status int, err error)
}

// Run prints an OK/missing report and optionally installs missing harness files.
func Run(opts Options) error {
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}
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
	vaultDir := opts.VaultDir
	if vaultDir == "" {
		vaultDir = "."
	}
	abs, err := vault.Abs(vaultDir)
	if err != nil {
		return err
	}

	checks := collectChecks(home, cwd)
	checks = append(checks, checkOllama(opts.HTTPGet))

	var missingTargets []string
	seen := map[string]bool{}
	for _, c := range checks {
		status := "OK"
		if !c.OK {
			status = "missing"
			if c.Target != "" && !seen[c.Target] {
				missingTargets = append(missingTargets, c.Target)
				seen[c.Target] = true
			}
		}
		line := fmt.Sprintf("%-8s %s", status, c.Name)
		if c.Path != "" {
			line += "  (" + c.Path + ")"
		}
		if c.Detail != "" {
			line += "  " + c.Detail
		}
		fmt.Fprintln(out, line)
	}

	if !opts.Fix {
		return nil
	}
	if len(missingTargets) == 0 {
		fmt.Fprintln(out, "fix: nothing relevant to install")
		return nil
	}
	if err := vault.Validate(abs); err != nil {
		return fmt.Errorf("fix requires a valid --vault: %w", err)
	}
	if opts.Home != "" {
		oldHome := os.Getenv("HOME")
		if err := os.Setenv("HOME", opts.Home); err != nil {
			return err
		}
		defer os.Setenv("HOME", oldHome)
	}
	for _, t := range missingTargets {
		fmt.Fprintf(out, "fix: sync --install --target %s --vault %s\n", t, abs)
		if err := sync.Sync(abs, t, filepath.Join("out", t), true); err != nil {
			return err
		}
	}
	return nil
}

func collectChecks(home, cwd string) []Check {
	piAgent := filepath.Join(home, ".pi", "agent")
	piSkill := filepath.Join(piAgent, "skills", "worlds-brain")
	claudeSkills := filepath.Join(home, ".claude", "skills")
	claudeBrain := filepath.Join(claudeSkills, "worlds-brain")
	cursorRules := filepath.Join(cwd, ".cursor", "rules")
	cursorBrain := filepath.Join(cursorRules, "worlds-brain.mdc")
	agentsWorlds := filepath.Join(cwd, "AGENTS.worlds.md")
	agents := filepath.Join(cwd, "AGENTS.md")

	return []Check{
		statCheck("pi agent dir", piAgent, "pi", true),
		statCheck("pi worlds-brain skill", piSkill, "pi", false),
		statCheck("claude skills dir", claudeSkills, "claude", true),
		statCheck("claude worlds-brain skill", claudeBrain, "claude", false),
		statCheck("cursor rules dir", cursorRules, "cursor", true),
		statCheck("cursor worlds-brain rule", cursorBrain, "cursor", false),
		eitherCheck("AGENTS.worlds.md or AGENTS.md", "codex", agentsWorlds, agents),
	}
}

func statCheck(name, path, target string, dirOK bool) Check {
	info, err := os.Stat(path)
	ok := err == nil
	if ok && dirOK && !info.IsDir() {
		ok = false
	}
	return Check{Name: name, Path: path, OK: ok, Target: target}
}

func eitherCheck(name, target string, paths ...string) Check {
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return Check{Name: name, Path: p, OK: true, Target: target}
		}
	}
	return Check{Name: name, Path: strings.Join(paths, " | "), OK: false, Target: target}
}

func checkOllama(get func(url string) (int, error)) Check {
	c := Check{Name: "ollama (optional)", Path: "http://127.0.0.1:11434", Target: ""}
	if get == nil {
		get = defaultHTTPGet
	}
	status, err := get("http://127.0.0.1:11434")
	if err != nil {
		c.OK = false
		c.Detail = err.Error()
		return c
	}
	c.OK = status > 0
	c.Detail = fmt.Sprintf("HTTP %d", status)
	return c
}

func defaultHTTPGet(url string) (int, error) {
	client := &http.Client{Timeout: 800 * time.Millisecond}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}
