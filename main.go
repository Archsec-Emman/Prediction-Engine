package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Archsec-Emman/Prediction-Engine/internal/engine"
)

const version = "1.0.0"

const banner = `
 Prediction Engine %s
 A single launcher for three open-source prediction-market tools.

 This binary does not implement any trading logic itself. It installs,
 checks and runs three upstream projects in parallel and reports their
 real exit status. All credit belongs to the upstream authors - see
 'prediction-engine status' for the full attribution list.
`

var engines = engine.Registry()

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "status":
		cmdStatus()
	case "install":
		cmdInstall(os.Args[2:])
	case "run":
		cmdRun(os.Args[2:])
	case "version", "-v", "--version":
		fmt.Println("prediction-engine", version)
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(banner+"\n", version)
	fmt.Println(`Usage:
  prediction-engine status              show engines, dependencies and install state
  prediction-engine install <name|all>  clone upstream repo and run install steps
  prediction-engine run <name|all> [--] [args...]
                                        run an engine (extra args passed through)
  prediction-engine version`)
}

func cmdStatus() {
	fmt.Printf("%-12s %-9s %-14s %s\n", "ENGINE", "INSTALLED", "DEPS", "UPSTREAM")
	for _, e := range engines {
		inst := "no"
		if engine.Installed(e) {
			inst = "yes"
		}
		var deps []string
		for _, d := range e.Requires {
			_, err := engine.CheckDep(d)
			mark := "ok"
			if err != nil {
				mark = "MISSING"
			}
			deps = append(deps, d.Bin+":"+mark)
		}
		fmt.Printf("%-12s %-9s %-14s %s\n", e.Name, inst, strings.Join(deps, " "), e.Upstream)
	}
	fmt.Println("\nCredits:")
	for _, e := range engines {
		fmt.Printf("  %-12s %s (%s)\n", e.Name, e.Credit, e.Upstream)
	}
}

func cmdInstall(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: prediction-engine install <name|all>")
		os.Exit(1)
	}
	targets, err := selectEngines(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, e := range targets {
		fmt.Printf("==> %s: %s\n", e.Name, e.Upstream)
		if err := installEngine(e); err != nil {
			fmt.Fprintf(os.Stderr, "install %s failed: %v\n", e.Name, err)
			os.Exit(1)
		}
	}
	fmt.Println("install complete")
}

func installEngine(e engine.Engine) error {
	if err := os.MkdirAll(engine.InstallRoot(), 0o755); err != nil {
		return err
	}
	if !engine.Installed(e) {
		if err := engine.Run("git", filepath.Dir(engine.InstallRoot()),
			[]string{"clone", "--depth", "1", e.Upstream, engine.EngineDir(e)}, nil); err != nil {
			return fmt.Errorf("clone: %w", err)
		}
	} else {
		fmt.Printf("==> %s already cloned, updating\n", e.Name)
		if err := engine.Run("git", engine.EngineDir(e), []string{"pull", "--ff-only"}, nil); err != nil {
			fmt.Printf("==> update skipped (%v)\n", err)
		}
	}
	for _, step := range e.Install {
		fmt.Printf("==> %s: %s\n", e.Name, strings.Join(step.Cmd, " "))
		if err := engine.Run(step.Cmd[0], engine.EngineDir(e), step.Cmd[1:], nil); err != nil {
			return fmt.Errorf("%s: %w", strings.Join(step.Cmd, " "), err)
		}
	}
	return nil
}

func cmdRun(args []string) {
	if len(args) < 1 {
		fmt.Println("usage: prediction-engine run <name|all> [--] [args...]")
		os.Exit(1)
	}
	passThrough := args[1:]
	if len(passThrough) > 0 && passThrough[0] == "--" {
		passThrough = passThrough[1:]
	}
	targets, err := selectEngines(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	results := make(map[string]error)
	var mu sync.Mutex
	order := make([]string, 0, len(targets))
	for _, e := range targets {
		if !engine.Installed(e) {
			fmt.Printf("[%-12s] not installed - run 'prediction-engine install %s' first\n", e.Name, e.Name)
			continue
		}
		order = append(order, e.Name)
		wg.Add(1)
		go func(e engine.Engine) {
			defer wg.Done()
			cmdArgs := append(append([]string{}, e.DefaultRun...), passThrough...)
			fmt.Printf("[%-12s] starting: %s\n", e.Name, strings.Join(cmdArgs, " "))
			err := engine.Run(cmdArgs[0], engine.EngineDir(e), cmdArgs[1:], nil)
			mu.Lock()
			results[e.Name] = err
			mu.Unlock()
		}(e)
	}
	wg.Wait()
	fmt.Print(engine.FormatSummary(results, order))
	failed := false
	for _, err := range results {
		if err != nil {
			failed = true
			break
		}
	}
	if failed {
		os.Exit(1)
	}
}

func selectEngines(name string) ([]engine.Engine, error) {
	if name == "all" {
		return engines, nil
	}
	e, err := engine.ByName(name)
	if err != nil {
		return nil, err
	}
	return []engine.Engine{*e}, nil
}
