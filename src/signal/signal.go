// Package signal implements graceful shutdown and platform signal
// handling per AI.md PART 8 (Signal Handling). Unix and Windows differ
// via build tags; this file holds the shared, platform-independent core.
package signal

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// shuttingDown is the shutdown flag surfaced to health checks (503).
var shuttingDown atomic.Bool

// exitFunc allows tests to intercept process exit.
var exitFunc = os.Exit

// childPIDs tracks child processes (Tor, etc.) to stop on shutdown.
var (
	childMu   sync.Mutex
	childPIDs []int
)

// Hook functions are the seams later PARTs plug into: log reopen/flush
// (PART 11), status dump (PART 8), database close (PART 10). Each defaults
// to a safe no-op so PART 8 is complete on its own.
var (
	hookMu       sync.Mutex
	onReopenLogs func()
	onDumpStatus func()
	onCloseDB    func()
	onFlushLogs  func()
)

// Setup installs the platform signal handler for the given HTTP server
// and PID file path.
func Setup(server *http.Server, pidFile string) {
	setupSignalHandler(server, pidFile)
}

// Kill sends a termination signal to a process: SIGTERM when graceful on
// Unix, SIGKILL otherwise; TerminateProcess on Windows.
func Kill(pid int, graceful bool) error {
	return killProcess(pid, graceful)
}

// IsShuttingDown reports whether graceful shutdown has begun; health
// checks use it to return 503.
func IsShuttingDown() bool {
	return shuttingDown.Load()
}

// setShuttingDown sets the shutdown flag for health checks.
func setShuttingDown(v bool) {
	shuttingDown.Store(v)
}

// RegisterChildPID records a child process to stop during shutdown.
func RegisterChildPID(pid int) {
	childMu.Lock()
	defer childMu.Unlock()
	childPIDs = append(childPIDs, pid)
}

// UnregisterChildPID removes a child process from the shutdown list.
func UnregisterChildPID(pid int) {
	childMu.Lock()
	defer childMu.Unlock()
	for i, p := range childPIDs {
		if p == pid {
			childPIDs = append(childPIDs[:i], childPIDs[i+1:]...)
			return
		}
	}
}

// getChildPIDs returns a snapshot of registered child PIDs.
func getChildPIDs() []int {
	childMu.Lock()
	defer childMu.Unlock()
	out := make([]int, len(childPIDs))
	copy(out, childPIDs)
	return out
}

// SetReopenLogsHook registers the log-rotation handler (SIGUSR1).
func SetReopenLogsHook(fn func()) {
	hookMu.Lock()
	defer hookMu.Unlock()
	onReopenLogs = fn
}

// SetDumpStatusHook registers the status-dump handler (SIGUSR2).
func SetDumpStatusHook(fn func()) {
	hookMu.Lock()
	defer hookMu.Unlock()
	onDumpStatus = fn
}

// SetCloseDatabaseHook registers the database close handler.
func SetCloseDatabaseHook(fn func()) {
	hookMu.Lock()
	defer hookMu.Unlock()
	onCloseDB = fn
}

// SetFlushLogsHook registers the log flush handler.
func SetFlushLogsHook(fn func()) {
	hookMu.Lock()
	defer hookMu.Unlock()
	onFlushLogs = fn
}

// getHook returns a registered hook under lock.
func getHook(which *func()) func() {
	hookMu.Lock()
	defer hookMu.Unlock()
	return *which
}

// reopenLogs runs the registered log-rotation hook (SIGUSR1).
func reopenLogs() {
	if fn := getHook(&onReopenLogs); fn != nil {
		fn()
	}
}

// dumpStatus runs the registered status-dump hook (SIGUSR2).
func dumpStatus() {
	if fn := getHook(&onDumpStatus); fn != nil {
		fn()
	}
}

// runWithTimeout runs fn but returns after at most timeout, per the
// shutdown timeout table (exceeding phases are logged and skipped).
func runWithTimeout(name string, fn func(), timeout time.Duration) {
	if fn == nil {
		return
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		log.Printf("%s did not finish within %v, continuing shutdown", name, timeout)
	}
}

// closeDatabase closes database connections with a timeout.
func closeDatabase(timeout time.Duration) {
	runWithTimeout("database close", getHook(&onCloseDB), timeout)
}

// flushLogs flushes buffered logs with a timeout.
func flushLogs(timeout time.Duration) {
	runWithTimeout("log flush", getHook(&onFlushLogs), timeout)
}

// gracefulShutdown performs orderly shutdown (cross-platform): stop
// accepting connections, drain in-flight requests (30s), stop children
// (10s), close database (5s), flush logs (2s), remove PID file, exit 0.
func gracefulShutdown(server *http.Server, pidFile string) {
	// Set shutdown flag for health checks
	setShuttingDown(true)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop accepting new connections, wait for in-flight requests.
	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

	// Stop child processes (Tor, etc.) - platform-specific.
	stopChildProcesses(10 * time.Second)

	closeDatabase(5 * time.Second)

	flushLogs(2 * time.Second)

	if pidFile != "" {
		os.Remove(pidFile)
	}

	log.Println("Graceful shutdown complete")
	exitFunc(0)
}
