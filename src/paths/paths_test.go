package paths

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigFileName(t *testing.T) {
	if ConfigFileName != "server.yml" {
		t.Errorf("ConfigFileName = %q, want server.yml (never server.yaml)", ConfigFileName)
	}
}

// checkDirs verifies structural invariants every platform's Dirs must hold.
func checkDirs(t *testing.T, name string, d Dirs) {
	t.Helper()
	if d.ConfigFile != filepath.Join(d.ConfigDir, ConfigFileName) {
		t.Errorf("%s: ConfigFile %q not ConfigDir + server.yml", name, d.ConfigFile)
	}
	if !strings.HasPrefix(d.LogFile, d.LogDir) {
		t.Errorf("%s: LogFile %q outside LogDir %q", name, d.LogFile, d.LogDir)
	}
	fields := map[string]string{
		"Binary": d.Binary, "ConfigDir": d.ConfigDir, "DataDir": d.DataDir,
		"CacheDir": d.CacheDir, "LogDir": d.LogDir, "BackupDir": d.BackupDir,
		"SSLDir": d.SSLDir, "SecDir": d.SecDir, "DBDir": d.DBDir,
	}
	for field, value := range fields {
		if value == "" {
			t.Errorf("%s: %s is empty", name, field)
		}
	}
}

func TestLinuxDirs(t *testing.T) {
	priv := linuxDirs(true)
	checkDirs(t, "linux privileged", priv)
	if priv.ConfigFile != "/etc/tabssh/tabssh/server.yml" {
		t.Errorf("linux privileged ConfigFile = %q", priv.ConfigFile)
	}
	if priv.DataDir != "/var/lib/tabssh/tabssh" {
		t.Errorf("linux privileged DataDir = %q", priv.DataDir)
	}
	if priv.LogDir != "/var/log/tabssh/tabssh" {
		t.Errorf("linux privileged LogDir = %q", priv.LogDir)
	}

	user := linuxDirs(false)
	checkDirs(t, "linux user", user)
	if !strings.HasSuffix(user.ConfigFile, ".config/tabssh/tabssh/server.yml") {
		t.Errorf("linux user ConfigFile = %q", user.ConfigFile)
	}
	if !strings.HasSuffix(user.DataDir, ".local/share/tabssh/tabssh") {
		t.Errorf("linux user DataDir = %q", user.DataDir)
	}
}

func TestDarwinDirs(t *testing.T) {
	priv := darwinDirs(true)
	checkDirs(t, "darwin privileged", priv)
	if priv.ConfigFile != "/Library/Application Support/tabssh/tabssh/server.yml" {
		t.Errorf("darwin privileged ConfigFile = %q", priv.ConfigFile)
	}

	user := darwinDirs(false)
	checkDirs(t, "darwin user", user)
	if !strings.HasSuffix(user.ConfigFile, "Library/Application Support/tabssh/tabssh/server.yml") {
		t.Errorf("darwin user ConfigFile = %q", user.ConfigFile)
	}
}

func TestBSDDirs(t *testing.T) {
	priv := bsdDirs(true)
	checkDirs(t, "bsd privileged", priv)
	if priv.ConfigFile != "/usr/local/etc/tabssh/tabssh/server.yml" {
		t.Errorf("bsd privileged ConfigFile = %q", priv.ConfigFile)
	}

	user := bsdDirs(false)
	checkDirs(t, "bsd user", user)
	if !strings.HasSuffix(user.ConfigFile, ".config/tabssh/tabssh/server.yml") {
		t.Errorf("bsd user ConfigFile = %q", user.ConfigFile)
	}
}

func TestWindowsDirs(t *testing.T) {
	t.Setenv("ProgramData", `C:\ProgramData`)
	t.Setenv("AppData", `C:\Users\u\AppData\Roaming`)
	t.Setenv("LocalAppData", `C:\Users\u\AppData\Local`)

	priv := windowsDirs(true)
	if priv.ConfigFile != filepath.Join(`C:\ProgramData`, "tabssh", "tabssh", "server.yml") {
		t.Errorf("windows privileged ConfigFile = %q", priv.ConfigFile)
	}

	user := windowsDirs(false)
	if user.ConfigFile != filepath.Join(`C:\Users\u\AppData\Roaming`, "tabssh", "tabssh", "server.yml") {
		t.Errorf("windows user ConfigFile = %q", user.ConfigFile)
	}
	if user.DataDir != filepath.Join(`C:\Users\u\AppData\Local`, "tabssh", "tabssh") {
		t.Errorf("windows user DataDir = %q", user.DataDir)
	}
}

func TestDockerDirs(t *testing.T) {
	d := dockerDirs()
	checkDirs(t, "docker", d)
	if d.ConfigFile != "/config/tabssh/server.yml" {
		t.Errorf("docker ConfigFile = %q", d.ConfigFile)
	}
	if d.DataDir != "/data/tabssh" {
		t.Errorf("docker DataDir = %q", d.DataDir)
	}
	if d.DBDir != "/data/db/sqlite" {
		t.Errorf("docker DBDir = %q", d.DBDir)
	}
}

func TestResolve(t *testing.T) {
	// Resolve depends on the runtime OS/container context, so only invariants
	// are asserted; the per-platform helpers are tested deterministically.
	d := Resolve()
	checkDirs(t, "resolve", d)
	if filepath.Base(d.ConfigFile) != ConfigFileName {
		t.Errorf("Resolve ConfigFile = %q, want basename server.yml", d.ConfigFile)
	}
}
