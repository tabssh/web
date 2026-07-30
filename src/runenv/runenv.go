// Package runenv detects TabSSH Web's runtime execution context — container,
// service manager, parent process, and elevation — per AI.md PART 8
// (Daemonize Flag, Container Detection, Service Manager Detection).
package runenv

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/tabssh/web/src/paths"
)

var (
	elevatedOnce  sync.Once
	elevatedValue bool
)

// StartedElevated reports whether the process started with elevated
// privileges. The value is locked on first call and never re-derived, so a
// later privilege drop cannot flip directory-mode decisions mid-run.
func StartedElevated() bool {
	elevatedOnce.Do(func() {
		elevatedValue = paths.IsPrivileged()
	})
	return elevatedValue
}

// IsContainer reports whether the process is running inside a container
// (Docker, Podman, LXC, Kubernetes, or a container init/entrypoint).
func IsContainer() bool {
	containerFiles := []string{
		// Docker
		"/.dockerenv",
		// Podman
		"/run/.containerenv",
		// LXC/LXD/Incus
		"/dev/lxc",
	}
	for _, f := range containerFiles {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}

	// Generic container env var (systemd-nspawn, lxc, etc.), any value.
	if os.Getenv("container") != "" {
		return true
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	// Container init systems supervise as the parent process.
	parentName := ParentProcessName()
	switch parentName {
	case "tini", "dumb-init", "s6-svscan", "runsv", "runsvdir", "catatonit":
		return true
	case "tabssh":
		// Parent is our own binary - likely a container entrypoint wrapper.
		return true
	}

	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") ||
			strings.Contains(content, "kubepods") ||
			strings.Contains(content, "lxc") {
			return true
		}
	}

	return false
}

// DetectServiceManager returns the active service manager: "container",
// "systemd", "launchd", "runit", "s6", "sysv", "rcd", or "manual".
func DetectServiceManager() string {
	if IsContainer() {
		return "container"
	}

	ppid := os.Getppid()

	if ppid == 1 {
		if _, err := os.Stat("/run/systemd/system"); err == nil {
			return "systemd"
		}
	}
	// INVOCATION_ID is set by systemd for supervised units.
	if os.Getenv("INVOCATION_ID") != "" {
		return "systemd"
	}

	if runtime.GOOS == "darwin" && ppid == 1 {
		return "launchd"
	}

	if os.Getenv("SVDIR") != "" {
		return "runit"
	}

	if os.Getenv("S6_LOGGING") != "" {
		return "s6"
	}

	if ppid == 1 {
		if _, err := os.Stat("/etc/init.d"); err == nil {
			if _, err := os.Stat("/run/systemd/system"); os.IsNotExist(err) {
				return "sysv"
			}
		}
	}

	if _, err := os.Stat("/etc/rc.subr"); err == nil {
		return "rcd"
	}

	return "manual"
}

// ParentProcessName returns the name of the parent process, or "" when it
// cannot be determined.
func ParentProcessName() string {
	ppid := os.Getppid()

	// Linux: read /proc/{ppid}/comm.
	if data, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", ppid)); err == nil {
		return strings.TrimSpace(string(data))
	}

	// macOS/BSD: fall back to the ps command.
	cmd := exec.Command("ps", "-p", strconv.Itoa(ppid), "-o", "comm=")
	if output, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(output))
	}

	return ""
}

// ShouldDaemonize determines whether to daemonize based on start context.
// A service start ignores the flag and config entirely: modern supervisors
// (systemd, launchd, runit, s6, containers) always get foreground; SysV and
// BSD rc.d always get a daemonized process. Manual starts respect the
// --daemon flag first, then the config setting, defaulting to foreground.
func ShouldDaemonize(isServiceStart, daemonFlag, configDaemonize bool) bool {
	if isServiceStart {
		switch DetectServiceManager() {
		case "systemd", "launchd", "runit", "s6", "docker", "container":
			return false
		case "sysv", "rcd":
			return true
		default:
			return false
		}
	}

	if daemonFlag {
		return true
	}
	return configDaemonize
}

// FilterDaemonFlag removes --daemon and -d from args so a re-exec'd daemon
// child cannot fork again in an infinite loop.
func FilterDaemonFlag(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "--daemon" && arg != "-d" {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
