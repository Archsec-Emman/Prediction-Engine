//go:build windows

package engine

import (
	"os/exec"
	"strings"
	"syscall"
)

func commandFor(tool string, args []string) *exec.Cmd {
	lower := strings.ToLower(tool)
	if strings.HasSuffix(lower, ".cmd") || strings.HasSuffix(lower, ".bat") {
		quoted := `"` + tool + `"`
		if len(args) > 0 {
			parts := make([]string, len(args))
			for i, a := range args {
				parts[i] = batchQuote(a)
			}
			quoted += " " + strings.Join(parts, " ")
		}
		cmd := exec.Command("cmd")
		cmd.Args = []string{"cmd"}
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `/d /s /c "` + quoted + `"`}
		return cmd
	}
	return exec.Command(tool, args...)
}

func batchQuote(a string) string {
	if !strings.ContainsAny(a, " \t\"") {
		return a
	}
	return `"` + strings.ReplaceAll(a, `"`, `""`) + `"`
}
