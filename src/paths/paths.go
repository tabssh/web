// Package paths resolves OS-specific and Docker filesystem locations for
// TabSSH Web. It implements AI.md PART 4 (OS-Specific Paths): distinct
// privileged and non-privileged path sets per platform, plus the simplified
// Docker-only /config and /data layout. internal_org and internal_name are
// frozen values from IDEA.md and must never change once set.
package paths

import (
	"os"
	"path/filepath"
	"runtime"
)

// internalOrg and internalName are frozen project identity values from
// IDEA.md ## Project variables. They must never be edited after first-time
// setup, even if project_name/project_org (the display/repo names) change.
const (
	internalOrg  = "tabssh"
	internalName = "tabssh"
	projectName  = "tabssh"
)

// ConfigFileName is always server.yml, never server.yaml, on every platform.
const ConfigFileName = "server.yml"

// IsRunningInDocker reports whether the process is running inside a Docker
// container, using the presence of /.dockerenv as the signal.
func IsRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

// IsPrivileged reports whether the process is running with elevated rights:
// UID 0 on Unix-like systems. Windows privilege detection is handled by the
// caller via an Administrator-elevation check (not covered by this pure
// path-resolution package).
func IsPrivileged() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	return os.Geteuid() == 0
}

// Dirs holds every resolved runtime directory/file location for the current
// platform and privilege level.
type Dirs struct {
	Binary     string
	ConfigDir  string
	ConfigFile string
	DataDir    string
	CacheDir   string
	LogDir     string
	LogFile    string
	BackupDir  string
	PIDFile    string
	SSLDir     string
	SecDir     string
	DBDir      string
}

// Resolve returns the correct Dirs for the current OS/container context.
func Resolve() Dirs {
	switch {
	case IsRunningInDocker():
		return dockerDirs()
	case runtime.GOOS == "windows":
		return windowsDirs(IsPrivileged())
	case runtime.GOOS == "darwin":
		return darwinDirs(IsPrivileged())
	case isBSD():
		return bsdDirs(IsPrivileged())
	default:
		return linuxDirs(IsPrivileged())
	}
}

func isBSD() bool {
	switch runtime.GOOS {
	case "freebsd", "openbsd", "netbsd":
		return true
	default:
		return false
	}
}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}

func linuxDirs(privileged bool) Dirs {
	if privileged {
		return Dirs{
			Binary:     filepath.Join("/usr/local/bin", projectName),
			ConfigDir:  filepath.Join("/etc", internalOrg, internalName),
			ConfigFile: filepath.Join("/etc", internalOrg, internalName, ConfigFileName),
			DataDir:    filepath.Join("/var/lib", internalOrg, internalName),
			CacheDir:   filepath.Join("/var/cache", internalOrg, internalName),
			LogDir:     filepath.Join("/var/log", internalOrg, internalName),
			LogFile:    filepath.Join("/var/log", internalOrg, internalName, "server.log"),
			BackupDir:  filepath.Join("/mnt/Backups", internalOrg, internalName),
			PIDFile:    filepath.Join("/var/run", internalOrg, internalName+".pid"),
			SSLDir:     filepath.Join("/etc", internalOrg, internalName, "ssl"),
			SecDir:     filepath.Join("/var/lib", internalOrg, internalName, "security"),
			DBDir:      filepath.Join("/var/lib", internalOrg, internalName, "db"),
		}
	}
	home := homeDir()
	return Dirs{
		Binary:     filepath.Join(home, ".local/bin", projectName),
		ConfigDir:  filepath.Join(home, ".config", internalOrg, internalName),
		ConfigFile: filepath.Join(home, ".config", internalOrg, internalName, ConfigFileName),
		DataDir:    filepath.Join(home, ".local/share", internalOrg, internalName),
		CacheDir:   filepath.Join(home, ".cache", internalOrg, internalName),
		LogDir:     filepath.Join(home, ".local/log", internalOrg, internalName),
		LogFile:    filepath.Join(home, ".local/log", internalOrg, internalName, "server.log"),
		BackupDir:  filepath.Join(home, ".local/share/Backups", internalOrg, internalName),
		PIDFile:    filepath.Join(home, ".local/share", internalOrg, internalName, internalName+".pid"),
		SSLDir:     filepath.Join(home, ".config", internalOrg, internalName, "ssl"),
		SecDir:     filepath.Join(home, ".local/share", internalOrg, internalName, "security"),
		DBDir:      filepath.Join(home, ".local/share", internalOrg, internalName, "db"),
	}
}

func darwinDirs(privileged bool) Dirs {
	if privileged {
		base := filepath.Join("/Library/Application Support", internalOrg, internalName)
		return Dirs{
			Binary:     filepath.Join("/usr/local/bin", projectName),
			ConfigDir:  base,
			ConfigFile: filepath.Join(base, ConfigFileName),
			DataDir:    filepath.Join(base, "data"),
			CacheDir:   filepath.Join("/Library/Caches", internalOrg, internalName),
			LogDir:     filepath.Join("/Library/Logs", internalOrg, internalName),
			LogFile:    filepath.Join("/Library/Logs", internalOrg, internalName, "server.log"),
			BackupDir:  filepath.Join("/Library/Backups", internalOrg, internalName),
			PIDFile:    filepath.Join("/var/run", internalOrg, internalName+".pid"),
			SSLDir:     filepath.Join(base, "ssl"),
			SecDir:     filepath.Join(base, "data/security"),
			DBDir:      filepath.Join(base, "db"),
		}
	}
	home := homeDir()
	base := filepath.Join(home, "Library/Application Support", internalOrg, internalName)
	return Dirs{
		Binary:     filepath.Join(home, "bin", projectName),
		ConfigDir:  base,
		ConfigFile: filepath.Join(base, ConfigFileName),
		DataDir:    base,
		CacheDir:   filepath.Join(home, "Library/Caches", internalOrg, internalName),
		LogDir:     filepath.Join(home, "Library/Logs", internalOrg, internalName),
		LogFile:    filepath.Join(home, "Library/Logs", internalOrg, internalName, "server.log"),
		BackupDir:  filepath.Join(home, "Library/Backups", internalOrg, internalName),
		PIDFile:    filepath.Join(base, internalName+".pid"),
		SSLDir:     filepath.Join(base, "ssl"),
		SecDir:     filepath.Join(base, "data/security"),
		DBDir:      filepath.Join(base, "db"),
	}
}

func bsdDirs(privileged bool) Dirs {
	if privileged {
		return Dirs{
			Binary:     filepath.Join("/usr/local/bin", projectName),
			ConfigDir:  filepath.Join("/usr/local/etc", internalOrg, internalName),
			ConfigFile: filepath.Join("/usr/local/etc", internalOrg, internalName, ConfigFileName),
			DataDir:    filepath.Join("/var/db", internalOrg, internalName),
			CacheDir:   filepath.Join("/var/cache", internalOrg, internalName),
			LogDir:     filepath.Join("/var/log", internalOrg, internalName),
			LogFile:    filepath.Join("/var/log", internalOrg, internalName, "server.log"),
			BackupDir:  filepath.Join("/var/backups", internalOrg, internalName),
			PIDFile:    filepath.Join("/var/run", internalOrg, internalName+".pid"),
			SSLDir:     filepath.Join("/usr/local/etc", internalOrg, internalName, "ssl"),
			SecDir:     filepath.Join("/var/db", internalOrg, internalName, "security"),
			DBDir:      filepath.Join("/var/db", internalOrg, internalName, "db"),
		}
	}
	home := homeDir()
	return Dirs{
		Binary:     filepath.Join(home, ".local/bin", projectName),
		ConfigDir:  filepath.Join(home, ".config", internalOrg, internalName),
		ConfigFile: filepath.Join(home, ".config", internalOrg, internalName, ConfigFileName),
		DataDir:    filepath.Join(home, ".local/share", internalOrg, internalName),
		CacheDir:   filepath.Join(home, ".cache", internalOrg, internalName),
		LogDir:     filepath.Join(home, ".local/log", internalOrg, internalName),
		LogFile:    filepath.Join(home, ".local/log", internalOrg, internalName, "server.log"),
		BackupDir:  filepath.Join(home, ".local/share/Backups", internalOrg, internalName),
		PIDFile:    filepath.Join(home, ".local/share", internalOrg, internalName, internalName+".pid"),
		SSLDir:     filepath.Join(home, ".config", internalOrg, internalName, "ssl"),
		SecDir:     filepath.Join(home, ".local/share", internalOrg, internalName, "security"),
		DBDir:      filepath.Join(home, ".local/share", internalOrg, internalName, "db"),
	}
}

func windowsDirs(privileged bool) Dirs {
	if privileged {
		programData := os.Getenv("ProgramData")
		base := filepath.Join(programData, internalOrg, internalName)
		return Dirs{
			Binary:     filepath.Join(`C:\Program Files`, internalOrg, internalName, projectName+".exe"),
			ConfigDir:  base,
			ConfigFile: filepath.Join(base, ConfigFileName),
			DataDir:    filepath.Join(base, "data"),
			CacheDir:   filepath.Join(base, "cache"),
			LogDir:     filepath.Join(base, "logs"),
			LogFile:    filepath.Join(base, "logs", "server.log"),
			BackupDir:  filepath.Join(programData, "Backups", internalOrg, internalName),
			SSLDir:     filepath.Join(base, "ssl"),
			SecDir:     filepath.Join(base, "data/security"),
			DBDir:      filepath.Join(base, "db"),
		}
	}
	localAppData := os.Getenv("LocalAppData")
	appData := os.Getenv("AppData")
	configBase := filepath.Join(appData, internalOrg, internalName)
	dataBase := filepath.Join(localAppData, internalOrg, internalName)
	return Dirs{
		Binary:     filepath.Join(localAppData, internalOrg, internalName, projectName+".exe"),
		ConfigDir:  configBase,
		ConfigFile: filepath.Join(configBase, ConfigFileName),
		DataDir:    dataBase,
		CacheDir:   filepath.Join(dataBase, "cache"),
		LogDir:     filepath.Join(dataBase, "logs"),
		LogFile:    filepath.Join(dataBase, "logs", "server.log"),
		BackupDir:  filepath.Join(localAppData, "Backups", internalOrg, internalName),
		SSLDir:     filepath.Join(configBase, "ssl"),
		SecDir:     filepath.Join(dataBase, "security"),
		DBDir:      filepath.Join(dataBase, "db"),
	}
}

// dockerDirs returns the simplified /config and /data layout used ONLY
// inside Docker containers, per AI.md PART 4 "Docker/Container".
func dockerDirs() Dirs {
	configDir := filepath.Join("/config", projectName)
	dataDir := filepath.Join("/data", projectName)
	return Dirs{
		Binary:     filepath.Join("/usr/local/bin", projectName),
		ConfigDir:  configDir,
		ConfigFile: filepath.Join(configDir, ConfigFileName),
		DataDir:    dataDir,
		CacheDir:   filepath.Join(dataDir, "cache"),
		LogDir:     filepath.Join("/data/log", projectName),
		LogFile:    filepath.Join("/data/log", projectName, "server.log"),
		BackupDir:  filepath.Join("/data/backups", projectName),
		SSLDir:     filepath.Join(configDir, "ssl"),
		SecDir:     filepath.Join(dataDir, "security"),
		DBDir:      "/data/db/sqlite",
	}
}
