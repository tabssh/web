//go:build !windows

// Unix daemonization by self re-exec, per AI.md PART 8 (Daemonization
// Process). The parent re-executes the binary in a new session with a
// marker env var, prints the child PID, and exits.
package runenv

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Daemonize forks the process and detaches it from the terminal. It returns
// nil immediately when already daemonized (parent is init) or when running
// as the re-exec'd child.
func Daemonize() error {
	if os.Getppid() == 1 {
		return nil
	}

	// The re-exec'd child carries a marker env var and just continues.
	if os.Getenv("_DAEMON_CHILD") != "" {
		return nil
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	// Re-exec with the same args minus --daemon to prevent a fork loop.
	args := FilterDaemonFlag(os.Args[1:])

	cmd := exec.Command(execPath, args...)
	cmd.Env = append(os.Environ(), "_DAEMON_CHILD=1")

	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Create a new session, detaching from the controlling terminal.
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting daemon: %w", err)
	}

	fmt.Printf("Daemon started with PID %d\n", cmd.Process.Pid)
	os.Exit(0)
	return nil
}
