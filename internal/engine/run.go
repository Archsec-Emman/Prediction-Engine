package engine

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func InstallRoot() string {
	return filepath.Join(DataDir(), "prediction-engine")
}

func EngineDir(e Engine) string {
	return filepath.Join(InstallRoot(), e.DirName)
}

func Installed(e Engine) bool {
	return dirExists(EngineDir(e))
}

func Run(name string, dir string, args []string, env []string) error {
	tool, err := ResolveTool(name)
	if err != nil {
		return err
	}
	cmd := commandFor(tool, args)
	return startStreaming(cmd, dir, env)
}

func startStreaming(cmd *exec.Cmd, dir string, env []string) error {
	cmd.Dir = dir
	cmd.WaitDelay = 30 * time.Second
	if len(env) > 0 {
		cmd.Env = env
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scan(stderr)
	}()
	go func() {
		defer wg.Done()
		scan(stdout)
	}()
	wg.Wait()
	return cmd.Wait()
}

func scan(r io.Reader) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for s.Scan() {
		fmt.Println(s.Text())
	}
}

func dirExists(p string) bool {
	info, err := Stat(p)
	return err == nil && info.IsDir()
}

func FormatSummary(results map[string]error, order []string) string {
	var b strings.Builder
	b.WriteString("\n=== Summary ===\n")
	for _, name := range order {
		if err, ok := results[name]; ok {
			if err != nil {
				fmt.Fprintf(&b, "  %-12s FAILED: %v\n", name, err)
			} else {
				fmt.Fprintf(&b, "  %-12s OK\n", name)
			}
		}
	}
	return b.String()
}
