// Package pid implements PID file handling with stale-file and PID-reuse
// detection, per AI.md PART 8 (PID File). Containers skip PID files
// entirely: the runtime supervises the process and namespace-local PIDs on
// shared volumes would point at the wrong process.
package pid

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tabssh/web/src/runenv"
)

// CheckPIDFile checks whether the PID file exists and whether the recorded
// process is still running and is actually our binary. Corrupt and stale
// files are removed. Returns (isRunning, pid, err).
func CheckPIDFile(pidPath string) (bool, int, error) {
	data, err := os.ReadFile(pidPath)
	if os.IsNotExist(err) {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, fmt.Errorf("reading pid file: %w", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		// Corrupt PID file - remove it.
		os.Remove(pidPath)
		return false, 0, nil
	}

	if !isProcessRunning(pid) {
		// Stale PID file - remove it.
		os.Remove(pidPath)
		return false, 0, nil
	}

	// The PID exists but may have been reused by another process.
	if !isOurProcess(pid) {
		os.Remove(pidPath)
		return false, 0, nil
	}

	return true, pid, nil
}

// WritePIDFile writes the current process PID to the file after verifying
// no other instance is running. Inside a container it is a no-op.
func WritePIDFile(pidPath string) error {
	if runenv.IsContainer() {
		return nil
	}

	running, existingPID, err := CheckPIDFile(pidPath)
	if err != nil {
		return err
	}
	if running {
		return fmt.Errorf("already running (pid %d)", existingPID)
	}

	pid := os.Getpid()
	return os.WriteFile(pidPath, []byte(strconv.Itoa(pid)), 0o644)
}

// RemovePIDFile removes the PID file on shutdown.
func RemovePIDFile(pidPath string) error {
	return os.Remove(pidPath)
}
