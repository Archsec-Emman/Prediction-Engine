package engine

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type Dep struct {
	Bin  string
	Args []string
	Desc string
}

type Step struct {
	Cmd []string
	Dir string
}

type RunMode int

const (
	RunCLI RunMode = iota
	RunServer
)

type Engine struct {
	Name        string
	Description string
	Upstream    string
	Credit      string
	DirName     string
	Requires    []Dep
	Install     []Step
	DefaultRun  []string
	Mode        RunMode
}

func (e Engine) InstallSummary() string {
	parts := make([]string, 0, len(e.Install))
	for _, s := range e.Install {
		parts = append(parts, strings.Join(s.Cmd, " "))
	}
	return strings.Join(parts, " && ")
}

func ResolveTool(bin string) (string, error) {
	if runtime.GOOS == "windows" {
		for _, ext := range []string{"", ".cmd", ".bat", ".exe"} {
			if p, err := exec.LookPath(bin + ext); err == nil {
				return p, nil
			}
		}
		return "", fmt.Errorf("%s not found in PATH", bin)
	}
	return exec.LookPath(bin)
}

func CheckDep(dep Dep) (string, error) {
	path, err := ResolveTool(dep.Bin)
	if err != nil {
		return "", err
	}
	if len(dep.Args) == 0 {
		return path, nil
	}
	out, err := exec.Command(path, dep.Args...).Output()
	if err != nil {
		return path, nil
	}
	first := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	return first, nil
}
