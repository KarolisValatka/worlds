package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KarolisValatka/worlds/internal/sync"
	"github.com/KarolisValatka/worlds/internal/vault"
)

func Run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return fmt.Errorf("missing command")
	}
	switch args[0] {
	case "help", "-h", "--help":
		printHelp()
		return nil
	case "init":
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}
		return vault.Init(dir)
	case "validate":
		dir := "."
		if len(args) > 1 {
			dir = args[1]
		}
		return vault.Validate(dir)
	case "sync":
		fs := flag.NewFlagSet("sync", flag.ContinueOnError)
		vaultDir := fs.String("vault", ".", "path to vault")
		target := fs.String("target", "", "pi|claude|cursor|codex")
		outDir := fs.String("out", "", "output directory (default out/<target>)")
		install := fs.Bool("install", false, "write into harness home paths")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *target == "" {
			return fmt.Errorf("--target is required (pi|claude|cursor|codex)")
		}
		out := *outDir
		if out == "" {
			out = filepath.Join("out", *target)
		}
		return sync.Sync(*vaultDir, *target, out, *install)
	default:
		printHelp()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printHelp() {
	fmt.Fprint(os.Stderr, `worlds - one vault, every coding agent

Usage:
  worlds init [dir]
  worlds validate [dir]
  worlds sync --vault <dir> --target <pi|claude|cursor|codex> [--out dir] [--install]

`)
}
