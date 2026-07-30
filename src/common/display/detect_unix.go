//go:build !windows

// Unix (Linux/BSD/macOS) display detection, per AI.md PART 7
// (Platform-Specific Display Detection).
package display

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// detectPlatformDisplay performs Unix/macOS display detection. Wayland is
// checked before X11 (preferred on Linux); on macOS the display is available
// unless running over SSH or as a LaunchDaemon with no GUI session.
func (e *DisplayEnv) detectPlatformDisplay() {
	if waylandDisplay := os.Getenv("WAYLAND_DISPLAY"); waylandDisplay != "" {
		e.HasDisplay = true
		e.DisplayType = "wayland"
		return
	}

	if display := os.Getenv("DISPLAY"); display != "" {
		e.HasDisplay = true
		e.DisplayType = "x11"
		return
	}

	if runtime.GOOS == "darwin" {
		if !e.IsSSH && os.Getenv("__CFBundleIdentifier") != "" {
			e.HasDisplay = true
			e.DisplayType = "macos"
			return
		}
		// Check whether the WindowServer (Aqua session) is accessible.
		cmd := exec.Command("launchctl", "managername")
		if output, err := cmd.Output(); err == nil {
			if strings.Contains(string(output), "Aqua") {
				e.HasDisplay = true
				e.DisplayType = "macos"
				return
			}
		}
	}

	e.HasDisplay = false
	e.DisplayType = "none"
}
