// Directory resolution and creation, per AI.md PART 8 (Directory Flags and
// Environment Variable Fallbacks). Every value resolves as CLI flag >
// environment variable > compiled default; missing directories are created
// automatically with mode-appropriate permissions.
package runenv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tabssh/web/src/paths"
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

// GetDatabaseDir resolves the database directory: DATABASE_DIR env overrides
// the platform default. defaults.DBDir already encodes the PART 4 location
// for every OS/privilege level and the Docker /data/db/sqlite layout, so it
// is the authoritative default rather than a re-derivation from the data dir.
func GetDatabaseDir(defaults paths.Dirs) string {
	if v := os.Getenv("DATABASE_DIR"); v != "" {
		return v
	}
	return defaults.DBDir
}

// GetBackupDir resolves the backup directory: --backup flag > BACKUP_DIR env
// > the PART 4 platform default (defaults.BackupDir). The platform default is
// authoritative and is never probed for writability or re-derived from the
// data dir — that produced off-spec locations on macOS/Windows/Docker.
func GetBackupDir(flagVal string, defaults paths.Dirs) string {
	if flagVal != "" {
		return flagVal
	}
	if v := os.Getenv("BACKUP_DIR"); v != "" {
		return v
	}
	return defaults.BackupDir
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
