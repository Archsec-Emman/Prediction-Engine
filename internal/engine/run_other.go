//go:build !windows

package engine

import "os/exec"

func commandFor(tool string, args []string) *exec.Cmd {
	return exec.Command(tool, args...)
}
