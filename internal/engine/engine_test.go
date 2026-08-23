package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func writeScript(t *testing.T, dir, name, body string) string {
	t.Helper()
	full := filepath.Join(dir, name)
	if runtime.GOOS == "windows" {
		full += ".cmd"
		body = "@echo off\r\n" + body
	} else {
		body = "#!/bin/sh\n" + body
	}
	if err := os.WriteFile(full, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	return full
}

func TestResolveToolFindsScriptOnPath(t *testing.T) {
	dir := t.TempDir()
	writeScript(t, dir, "fake-tool", "echo ok")
	t.Setenv("PATH", dir)
	if _, err := ResolveTool("fake-tool"); err != nil {
		t.Fatalf("expected fake-tool to resolve: %v", err)
	}
	if _, err := ResolveTool("definitely-not-a-tool-xyz"); err == nil {
		t.Fatal("expected missing tool to error")
	}
}

func TestCheckDepReportsVersion(t *testing.T) {
	dir := t.TempDir()
	writeScript(t, dir, "fake-ver", "echo v9.9.9")
	t.Setenv("PATH", dir)
	got, err := CheckDep(Dep{Bin: "fake-ver", Args: []string{"--version"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "v9.9.9") {
		t.Fatalf("expected version output, got %q", got)
	}
}

const helperEnv = "GO_ENGINE_HELPER_MODE"

func TestHelperProcess(t *testing.T) {
	switch os.Getenv(helperEnv) {
	case "stdout":
		fmt.Println("hello-from-ok")
		os.Exit(0)
	case "exit3":
		os.Exit(3)
	}
}

func helperCommand(mode string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "-test.v=false")
	cmd.Env = append(os.Environ(), helperEnv+"="+mode)
	return cmd
}

func TestStartStreamingOutputAndExitCode(t *testing.T) {
	if err := startStreaming(helperCommand("stdout"), t.TempDir(), nil); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	err := startStreaming(helperCommand("exit3"), t.TempDir(), nil)
	if err == nil {
		t.Fatal("expected failure exit code")
	}
	if !strings.Contains(err.Error(), "3") {
		t.Fatalf("expected exit code 3 in error, got %v", err)
	}
}

func TestFormatSummary(t *testing.T) {
	results := map[string]error{"a": nil, "b": fmt.Errorf("boom")}
	out := FormatSummary(results, []string{"a", "b"})
	okRe := regexp.MustCompile(`(?m)^\s+a\s+OK\s*$`)
	failRe := regexp.MustCompile(`(?m)^\s+b\s+FAILED: boom\s*$`)
	if !okRe.MatchString(out) || !failRe.MatchString(out) {
		t.Fatalf("unexpected summary:\n%s", out)
	}
}

func TestRegistryIsComplete(t *testing.T) {
	reg := Registry()
	if len(reg) != 3 {
		t.Fatalf("expected 3 engines, got %d", len(reg))
	}
	for _, e := range reg {
		if !strings.HasPrefix(e.Upstream, "https://github.com/") {
			t.Errorf("%s: upstream must be a github URL", e.Name)
		}
		if e.Credit == "" || e.Description == "" {
			t.Errorf("%s: missing attribution or description", e.Name)
		}
		if len(e.Requires) == 0 {
			t.Errorf("%s: no dependency list", e.Name)
		}
	}
	names := map[string]bool{}
	for _, e := range reg {
		names[e.Name] = true
	}
	for _, want := range []string{"polyseer", "papertrader", "backtest"} {
		if !names[want] {
			t.Errorf("engine %s missing from registry", want)
		}
	}
}

func TestByName(t *testing.T) {
	if _, err := ByName("papertrader"); err != nil {
		t.Fatalf("papertrader should resolve: %v", err)
	}
	if _, err := ByName("nope"); err == nil {
		t.Fatal("unknown engine should error")
	}
}
