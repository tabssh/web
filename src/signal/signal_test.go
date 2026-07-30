package signal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestShuttingDownFlag(t *testing.T) {
	t.Cleanup(func() { setShuttingDown(false) })
	tests := []struct {
		name string
		set  bool
		want bool
	}{
		{"initially false", false, false},
		{"set true", true, true},
		{"reset false", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setShuttingDown(tt.set)
			if got := IsShuttingDown(); got != tt.want {
				t.Errorf("IsShuttingDown() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChildPIDRegistry(t *testing.T) {
	tests := []struct {
		name       string
		register   []int
		unregister []int
		want       int
	}{
		{"register two", []int{101, 102}, nil, 2},
		{"unregister one", []int{201, 202}, []int{201}, 1},
		{"unregister unknown", []int{301}, []int{999}, 1},
		{"unregister all", []int{401, 402}, []int{401, 402}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.register {
				RegisterChildPID(p)
			}
			t.Cleanup(func() {
				for _, p := range tt.register {
					UnregisterChildPID(p)
				}
			})
			for _, p := range tt.unregister {
				UnregisterChildPID(p)
			}
			if got := len(getChildPIDs()); got != tt.want {
				t.Errorf("len(getChildPIDs()) = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestHooksInvoked(t *testing.T) {
	tests := []struct {
		name string
		set  func(func())
		run  func()
	}{
		{"reopen logs", SetReopenLogsHook, reopenLogs},
		{"dump status", SetDumpStatusHook, dumpStatus},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			tt.set(func() { called = true })
			t.Cleanup(func() { tt.set(nil) })
			tt.run()
			if !called {
				t.Error("hook was not invoked")
			}
		})
	}
}

func TestRunWithTimeoutBounded(t *testing.T) {
	tests := []struct {
		name    string
		fn      func()
		timeout time.Duration
		maxWait time.Duration
	}{
		{"nil hook returns immediately", nil, time.Second, 100 * time.Millisecond},
		{"fast hook completes", func() {}, time.Second, 500 * time.Millisecond},
		{"hanging hook bounded", func() { time.Sleep(5 * time.Second) }, 100 * time.Millisecond, time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			runWithTimeout("test phase", tt.fn, tt.timeout)
			if elapsed := time.Since(start); elapsed > tt.maxWait {
				t.Errorf("runWithTimeout took %v, want <= %v", elapsed, tt.maxWait)
			}
		})
	}
}

func TestGracefulShutdownSequence(t *testing.T) {
	dbClosed := false
	logsFlushed := false
	SetCloseDatabaseHook(func() { dbClosed = true })
	SetFlushLogsHook(func() { logsFlushed = true })
	t.Cleanup(func() {
		SetCloseDatabaseHook(nil)
		SetFlushLogsHook(nil)
		setShuttingDown(false)
		exitFunc = os.Exit
	})

	exitCode := -1
	exitFunc = func(code int) { exitCode = code }

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server := ts.Config
	defer ts.Close()

	pidFile := filepath.Join(t.TempDir(), "tabssh.pid")
	if err := os.WriteFile(pidFile, []byte("12345"), 0o644); err != nil {
		t.Fatalf("writing pid file: %v", err)
	}

	gracefulShutdown(server, pidFile)

	if !IsShuttingDown() {
		t.Error("shutdown flag not set")
	}
	if !dbClosed {
		t.Error("database close hook not invoked")
	}
	if !logsFlushed {
		t.Error("log flush hook not invoked")
	}
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Error("PID file not removed")
	}
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
}

func TestKillProcess(t *testing.T) {
	tests := []struct {
		name     string
		graceful bool
	}{
		{"graceful sigterm", true},
		{"force kill", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("sleep", "30")
			if err := cmd.Start(); err != nil {
				t.Skipf("cannot start helper process: %v", err)
			}
			defer cmd.Process.Kill()
			defer cmd.Wait()
			if err := Kill(cmd.Process.Pid, tt.graceful); err != nil {
				t.Errorf("Kill(%d, %v) error: %v", cmd.Process.Pid, tt.graceful, err)
			}
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Error("process did not exit after Kill")
			}
		})
	}
}

func TestStopChildProcessesTermination(t *testing.T) {
	cmd := exec.Command("sleep", "30")
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot start helper process: %v", err)
	}
	defer cmd.Process.Kill()
	RegisterChildPID(cmd.Process.Pid)
	t.Cleanup(func() { UnregisterChildPID(cmd.Process.Pid) })

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	stopChildProcesses(2 * time.Second)

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Error("child process still running after stopChildProcesses")
	}
}
