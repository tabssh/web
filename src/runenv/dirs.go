// Directory resolution and creation, per AI.md PART 8 (Directory Flags and
// Environment Variable Fallbacks). Every value resolves as CLI flag >
// environment variable > compiled default; missing directories are created
// automatically with mode-appropriate permissions.
package runenv

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/tabssh/web/src/paths"
)

const (
	internalOrg  = "tabssh"
	internalName = "tabssh"
)

// resolveDir applies the flag > env > default priority for one directory.
func resolveDir(flagVal, envName, defaultVal string) string {
	if flagVal != "" {
		return flagVal
	}
	if v := os.Getenv(envName); v != "" {
		return v
	}
	return defaultVal
}

// GetConfigDir resolves the config directory: --config > CONFIG_DIR > default.
func GetConfigDir(flagVal string, defaults paths.Dirs) string {
	return resolveDir(flagVal, "CONFIG_DIR", defaults.ConfigDir)
}

// GetDataDir resolves the data directory: --data > DATA_DIR > default.
func GetDataDir(flagVal string, defaults paths.Dirs) string {
	return resolveDir(flagVal, "DATA_DIR", defaults.DataDir)
}

// GetCacheDir resolves the cache directory: --cache > CACHE_DIR > default.
func GetCacheDir(flagVal string, defaults paths.Dirs) string {
	return resolveDir(flagVal, "CACHE_DIR", defaults.CacheDir)
}

// GetLogDir resolves the log directory: --log > LOG_DIR > default.
func GetLogDir(flagVal string, defaults paths.Dirs) string {
	return resolveDir(flagVal, "LOG_DIR", defaults.LogDir)
}

// GetPIDFile resolves the PID file path: --pid > PID_FILE > default.
func GetPIDFile(flagVal string, defaults paths.Dirs) string {
	return resolveDir(flagVal, "PID_FILE", defaults.PIDFile)
}

// GetDatabaseDir resolves the database directory: DATABASE_DIR > container
// default /data/db/sqlite > {data_dir}/db.
func GetDatabaseDir(dataDir string) string {
	if v := os.Getenv("DATABASE_DIR"); v != "" {
		return v
	}
	if paths.IsRunningInDocker() {
		return "/data/db/sqlite"
	}
	return filepath.Join(dataDir, "db")
}

// GetBackupDir resolves the backup directory: --backup > BACKUP_DIR > the
// OS system backup location when writable > {data_dir}/backup (elevated) or
// the per-user backup location.
func GetBackupDir(flagVal, dataDir string) string {
	if flagVal != "" {
		return flagVal
	}
	if v := os.Getenv("BACKUP_DIR"); v != "" {
		return v
	}
	if sys := systemBackupDir(); sys != "" && isWritable(sys) {
		return sys
	}
	if StartedElevated() {
		return filepath.Join(dataDir, "backup")
	}
	return userBackupDir()
}

// systemBackupDir returns the OS-level shared backup location.
func systemBackupDir() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join("/Library/Backups", internalOrg, internalName)
	case "windows":
		pd := os.Getenv("ProgramData")
		if pd == "" {
			return ""
		}
		return filepath.Join(pd, "Backups", internalOrg, internalName)
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		return filepath.Join("/var/backups", internalOrg, internalName)
	default:
		return filepath.Join("/mnt/Backups", internalOrg, internalName)
	}
}

// userBackupDir returns the per-user backup location.
func userBackupDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library/Backups", internalOrg, internalName)
	case "windows":
		lad := os.Getenv("LOCALAPPDATA")
		if lad == "" {
			lad = home
		}
		return filepath.Join(lad, "Backups", internalOrg, internalName)
	default:
		return filepath.Join(home, ".local/share/Backups", internalOrg, internalName)
	}
}

// isWritable reports whether a test file can be created under the parent of
// path, without requiring path itself to exist yet.
func isWritable(path string) bool {
	parent := filepath.Dir(path)
	info, err := os.Stat(parent)
	if err != nil || !info.IsDir() {
		return false
	}
	probe := filepath.Join(parent, ".write_test_"+strconv.FormatInt(time.Now().UnixNano(), 36))
	f, err := os.Create(probe)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(probe)
	return true
}

// EnsureDir creates the directory (with parents) if missing and verifies it
// is writable. Elevated mode uses 0755; user mode uses 0700.
func EnsureDir(path string, elevated bool) error {
	perm := os.FileMode(0o700)
	if elevated {
		perm = 0o755
	}
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("creating directory %s: %w", path, err)
	}
	probe := filepath.Join(path, ".write-test")
	f, err := os.Create(probe)
	if err != nil {
		return fmt.Errorf("directory %s is not writable: %w", path, err)
	}
	f.Close()
	os.Remove(probe)
	return nil
}

// EnsurePIDFileDir creates the parent directory of the PID file path.
func EnsurePIDFileDir(pidPath string, elevated bool) error {
	return EnsureDir(filepath.Dir(pidPath), elevated)
}
