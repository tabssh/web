//go:build windows

// Windows process liveness and identity checks, per AI.md PART 8 (PID
// File).
package pid

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// stillActive is the Windows STILL_ACTIVE exit code (STATUS_PENDING),
// returned by GetExitCodeProcess for a running process.
const stillActive = 259

// isProcessRunning checks whether a process with the given PID exists. On
// Windows, os.FindProcess succeeds for any PID value, so a real handle is
// obtained via OpenProcess and the exit code queried.
func isProcessRunning(pid int) bool {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)

	var exitCode uint32
	err = windows.GetExitCodeProcess(handle, &exitCode)
	return err == nil && exitCode == stillActive
}

// isOurProcess verifies the process image is the tabssh server binary.
// Exact case-insensitive matching is required: substring matching would
// also match tabssh-cli.exe.
func isOurProcess(pid int) bool {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer windows.CloseHandle(handle)

	var buf [windows.MAX_PATH]uint16
	var size uint32 = windows.MAX_PATH
	err = windows.QueryFullProcessImageName(handle, 0, &buf[0], &size)
	if err != nil {
		return false
	}
	exePath := windows.UTF16ToString(buf[:size])
	base := filepath.Base(exePath)
	return strings.EqualFold(base, "tabssh.exe") || strings.EqualFold(base, "tabssh")
}
