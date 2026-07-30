//go:build !windows

// Unix process liveness and identity checks, per AI.md PART 8 (PID File).
package pid

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// isProcessRunning checks whether a process with the given PID exists. On
// Unix, FindProcess always succeeds, so signal 0 probes the process; EPERM
// means it exists but belongs to another user.
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	return errors.Is(err, syscall.EPERM)
}

// isOurProcess verifies the process is actually the tabssh server binary.
// Exact-name matching is required: substring matching would also match
// tabssh-cli.
func isOurProcess(pid int) bool {
	exePath, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		// No procfs (macOS/BSD) - fall back to ps.
		return isOurProcessDarwin(pid)
	}
	return filepath.Base(exePath) == "tabssh"
}

// isOurProcessDarwin checks the process name on macOS/BSD via ps.
func isOurProcessDarwin(pid int) bool {
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "tabssh"
}
