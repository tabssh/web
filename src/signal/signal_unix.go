//go:build !windows

// Unix signal handling per AI.md PART 8: SIGTERM/SIGINT/SIGQUIT and
// SIGRTMIN+3 (Docker STOPSIGNAL, 37) trigger graceful shutdown, SIGUSR1
// reopens logs, SIGUSR2 dumps status, SIGHUP is ignored (config
// auto-reloads via file watcher).
package signal

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// setupSignalHandler configures graceful shutdown (Unix).
func setupSignalHandler(server *http.Server, pidFile string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	// Handle SIGRTMIN+3 (Docker STOPSIGNAL) - signal 37.
	signal.Notify(sigChan, syscall.Signal(37))

	// Ignore SIGHUP - config reloads automatically via file watcher.
	signal.Ignore(syscall.SIGHUP)

	go func() {
		for sig := range sigChan {
			switch sig {
			case syscall.SIGUSR1:
				log.Println("Received SIGUSR1, reopening logs...")
				reopenLogs()

			case syscall.SIGUSR2:
				log.Println("Received SIGUSR2, dumping status...")
				dumpStatus()

			default:
				// Graceful shutdown (SIGTERM, SIGINT, SIGQUIT, SIGRTMIN+3).
				log.Printf("Received %v, starting graceful shutdown...", sig)
				gracefulShutdown(server, pidFile)
			}
		}
	}()
}

// killProcess sends a signal to a process (Unix): SIGTERM when graceful,
// SIGKILL otherwise.
func killProcess(pid int, graceful bool) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if graceful {
		return process.Signal(syscall.SIGTERM)
	}
	return process.Signal(syscall.SIGKILL)
}

// stopChildProcesses sends SIGTERM to children, then SIGKILL after the
// timeout (Unix).
func stopChildProcesses(timeout time.Duration) {
	for _, pid := range getChildPIDs() {
		process, err := os.FindProcess(pid)
		if err != nil {
			continue
		}

		// Send SIGTERM (graceful).
		process.Signal(syscall.SIGTERM)
	}

	// Wait with timeout, then SIGKILL.
	deadline := time.Now().Add(timeout)
	for _, pid := range getChildPIDs() {
		process, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		exited := false
		for time.Now().Before(deadline) {
			if err := process.Signal(syscall.Signal(0)); err != nil {
				// Process exited.
				exited = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		// Force kill if still running.
		if !exited {
			process.Signal(syscall.SIGKILL)
		}
	}
}
