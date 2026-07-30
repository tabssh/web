//go:build windows

// Windows display detection, per AI.md PART 7 (Platform-Specific Display
// Detection).
package display

import (
	"os"

	"golang.org/x/sys/windows"
)

// detectPlatformDisplay performs Windows display detection. Windows always
// has a display unless the process runs as a service in session 0.
func (e *DisplayEnv) detectPlatformDisplay() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	hwnd, _, _ := getConsoleWindow.Call()

	// Session 0 is the service session with no interactive desktop.
	var sessionID uint32
	windows.ProcessIdToSessionId(windows.GetCurrentProcessId(), &sessionID)

	if sessionID == 0 {
		e.HasDisplay = false
		e.DisplayType = "none"
		return
	}

	// A SESSIONNAME env var indicates a remote-desktop session.
	if os.Getenv("SESSIONNAME") == "RDP-Tcp#0" || os.Getenv("SESSIONNAME") != "" {
		e.HasDisplay = true
		e.DisplayType = "windows-rdp"
		return
	}

	e.HasDisplay = hwnd != 0
	if e.HasDisplay {
		e.DisplayType = "windows"
	} else {
		e.DisplayType = "none"
	}
}
