package runenv

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tabssh/web/src/paths"
)

func TestFilterDaemonFlag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{"removes long flag", []string{"--daemon", "--port", "8080"}, []string{"--port", "8080"}},
		{"removes short flag", []string{"-d", "--mode", "production"}, []string{"--mode", "production"}},
		{"no daemon flag", []string{"--port", "8080"}, []string{"--port", "8080"}},
		{"empty args", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterDaemonFlag(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterDaemonFlag(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestShouldDaemonizeManual(t *testing.T) {
	tests := []struct {
		name       string
		daemonFlag bool
		configVal  bool
		want       bool
	}{
		{"default foreground", false, false, false},
		{"flag wins", true, false, true},
		{"config daemonize", false, true, true},
		{"both set", true, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldDaemonize(false, tt.daemonFlag, tt.configVal); got != tt.want {
				t.Errorf("ShouldDaemonize(false, %v, %v) = %v, want %v", tt.daemonFlag, tt.configVal, got, tt.want)
			}
		})
	}
}

func TestDetectServiceManagerReturnsKnownValue(t *testing.T) {
	known := map[string]bool{
		"container": true, "systemd": true, "launchd": true, "runit": true,
		"s6": true, "sysv": true, "rcd": true, "manual": true,
	}
	got := DetectServiceManager()
	if !known[got] {
		t.Errorf("DetectServiceManager() = %q, not a known manager", got)
	}
}

func TestResolveDirPriority(t *testing.T) {
	defaults := paths.Dirs{ConfigDir: "/default/config"}
	tests := []struct {
		name    string
		flagVal string
		envVal  string
		want    string
	}{
		{"flag wins", "/flag/config", "/env/config", "/flag/config"},
		{"env wins over default", "", "/env/config", "/env/config"},
		{"default fallback", "", "", "/default/config"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CONFIG_DIR", tt.envVal)
			if got := GetConfigDir(tt.flagVal, defaults); got != tt.want {
				t.Errorf("GetConfigDir(%q) = %q, want %q", tt.flagVal, got, tt.want)
			}
		})
	}
}

func TestGetBackupDirPriority(t *testing.T) {
	defaults := paths.Dirs{BackupDir: "/default/backup"}
	tests := []struct {
		name    string
		flagVal string
		envVal  string
		want    string
	}{
		{"flag wins", "/flag/backup", "/env/backup", "/flag/backup"},
		{"env wins over default", "", "/env/backup", "/env/backup"},
		{"platform default fallback", "", "", "/default/backup"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("BACKUP_DIR", tt.envVal)
			if got := GetBackupDir(tt.flagVal, defaults); got != tt.want {
				t.Errorf("GetBackupDir(%q) = %q, want %q", tt.flagVal, got, tt.want)
			}
		})
	}
}

func TestGetDatabaseDir(t *testing.T) {
	defaults := paths.Dirs{DBDir: "/default/db"}
	t.Setenv("DATABASE_DIR", "/env/db")
	if got := GetDatabaseDir(defaults); got != "/env/db" {
		t.Errorf("GetDatabaseDir with env = %q, want /env/db", got)
	}
	t.Setenv("DATABASE_DIR", "")
	if got := GetDatabaseDir(defaults); got != "/default/db" {
		t.Errorf("GetDatabaseDir default = %q, want /default/db", got)
	}
}

func TestEnsureDir(t *testing.T) {
	base := t.TempDir()
	tests := []struct {
		name     string
		elevated bool
		perm     os.FileMode
	}{
		{"user mode 0700", false, 0o700},
		{"elevated mode 0755", true, 0o755},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := filepath.Join(base, tt.name, "nested")
			if err := EnsureDir(dir, tt.elevated); err != nil {
				t.Fatalf("EnsureDir() error: %v", err)
			}
			info, err := os.Stat(dir)
			if err != nil {
				t.Fatalf("stat: %v", err)
			}
			if perm := info.Mode().Perm(); perm != tt.perm {
				t.Errorf("dir perm = %o, want %o", perm, tt.perm)
			}
			if _, err := os.Stat(filepath.Join(dir, ".write-test")); !os.IsNotExist(err) {
				t.Error("write-test probe file was not removed")
			}
		})
	}
}

func TestEnsurePIDFileDir(t *testing.T) {
	base := t.TempDir()
	pidPath := filepath.Join(base, "run", "tabssh.pid")
	if err := EnsurePIDFileDir(pidPath, false); err != nil {
		t.Fatalf("EnsurePIDFileDir() error: %v", err)
	}
	if info, err := os.Stat(filepath.Dir(pidPath)); err != nil || !info.IsDir() {
		t.Errorf("PID parent dir missing: %v", err)
	}
}

func TestStartedElevatedStable(t *testing.T) {
	first := StartedElevated()
	second := StartedElevated()
	if first != second {
		t.Error("StartedElevated() changed between calls")
	}
	if first != paths.IsPrivileged() && os.Geteuid() >= 0 {
		t.Logf("StartedElevated()=%v (locked value)", first)
	}
}
