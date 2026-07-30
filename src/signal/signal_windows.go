//go:build windows

// Windows signal handling per AI.md PART 8: only os.Interrupt (Ctrl+C,
// Ctrl+Break) is catchable; SIGHUP/SIGUSR1/SIGUSR2/SIGQUIT do not exist.
// Windows Service Control (SERVICE_CONTROL_STOP) is wired in via
// golang.org/x/sys/windows/svc in the service layer (PART 24/25).
package signal

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// setupSignalHandler configures graceful shutdown (Windows).
func setupSignalHandler(server *http.Server, pidFile string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for sig := range sigChan {
			log.Printf("Received %v, starting graceful shutdown...", sig)
			gracefulShutdown(server, pidFile)
		}
	}()
}

// killProcess terminates a process (Windows). Windows has no graceful
// signals - Kill() calls TerminateProcess regardless of the flag.
func killProcess(pid int, graceful bool) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

// stopChildProcesses terminates children (Windows). Windows cannot send
// graceful signals - immediate termination only.
func stopChildProcesses(timeout time.Duration) {
	for _, pid := range getChildPIDs() {
		process, err := os.FindProcess(pid)
		if err != nil {
			continue
		}
		// Windows: Kill() is immediate termination (TerminateProcess).
		process.Kill()
	}
}
