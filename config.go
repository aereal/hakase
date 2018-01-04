package main

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// from https://github.com/motemen/ghq/blob/master/git.go
func getConfig(key string) (string, error) {
	args := []string{"config", "--path", "--null", "--get", key}
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr

	buf, err := cmd.Output()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				// The key was not found, do not treat as an error
				return "", nil
			}
		}

		return "", err
	}

	return strings.TrimRight(string(buf), "\000"), nil
}
