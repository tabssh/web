package pid

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestCheckPIDFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		noFile      bool
		wantRunning bool
		wantRemoved bool
	}{
		{"missing file", "", true, false, false},
		{"corrupt content", "not-a-pid", false, false, true},
		{"dead process", "999999", false, false, true},
		// The current test process is running but is not the tabssh binary,
		// so PID-reuse detection must treat the file as stale.
		{"pid reuse different binary", strconv.Itoa(os.Getpid()), false, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pidPath := filepath.Join(t.TempDir(), "test.pid")
			if !tt.noFile {
				if err := os.WriteFile(pidPath, []byte(tt.content), 0o644); err != nil {
					t.Fatalf("writing pid file: %v", err)
				}
			}
			running, _, err := CheckPIDFile(pidPath)
			if err != nil {
				t.Fatalf("CheckPIDFile() error: %v", err)
			}
			if running != tt.wantRunning {
				t.Errorf("CheckPIDFile() running = %v, want %v", running, tt.wantRunning)
			}
			if tt.wantRemoved {
				if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
					t.Error("stale/corrupt pid file was not removed")
				}
			}
		})
	}
}

func TestWriteAndRemovePIDFile(t *testing.T) {
	pidPath := filepath.Join(t.TempDir(), "tabssh.pid")
	if err := WritePIDFile(pidPath); err != nil {
		t.Fatalf("WritePIDFile() error: %v", err)
	}
	// Inside a container WritePIDFile is a documented no-op.
	data, err := os.ReadFile(pidPath)
	if os.IsNotExist(err) {
		t.Skip("running in container: PID file skipped by design")
	}
	if err != nil {
		t.Fatalf("reading pid file: %v", err)
	}
	got, err := strconv.Atoi(string(data))
	if err != nil || got != os.Getpid() {
		t.Errorf("pid file content = %q, want %d", data, os.Getpid())
	}
	if err := RemovePIDFile(pidPath); err != nil {
		t.Fatalf("RemovePIDFile() error: %v", err)
	}
	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Error("pid file still exists after RemovePIDFile")
	}
}

func TestIsProcessRunning(t *testing.T) {
	tests := []struct {
		name string
		pid  int
		want bool
	}{
		{"current process", os.Getpid(), true},
		{"nonexistent pid", 999999, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isProcessRunning(tt.pid); got != tt.want {
				t.Errorf("isProcessRunning(%d) = %v, want %v", tt.pid, got, tt.want)
			}
		})
	}
}

func TestIsOurProcessRejectsTestBinary(t *testing.T) {
	// The test binary is never named exactly "tabssh".
	if isOurProcess(os.Getpid()) {
		t.Error("isOurProcess() = true for test binary, want false (exact-name match)")
	}
}
