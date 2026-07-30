// Core display-environment detection logic, per AI.md PART 7 (Display
// Environment Detection and TERM=dumb Handling).
package display

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// DisplayEnv is the detected display environment.
type DisplayEnv struct {
	Mode DisplayMode
	// HasDisplay is true when an X11, Wayland, Windows, or macOS display exists.
	HasDisplay bool
	// DisplayType is "x11", "wayland", "windows", "windows-rdp", "macos", or "none".
	DisplayType string
	// IsTerminal is true when stdout is a TTY.
	IsTerminal bool
	// IsSSH is true when running over SSH.
	IsSSH bool
	// IsMosh is true when running over mosh.
	IsMosh bool
	// IsScreen is true when running inside screen or tmux.
	IsScreen bool
	// TerminalType is the TERM value.
	TerminalType string
	// Cols is the terminal column count (0 if no terminal).
	Cols int
	// Rows is the terminal row count (0 if no terminal).
	Rows int
}

// DetectDisplayEnv auto-detects the display environment.
func DetectDisplayEnv() DisplayEnv {
	env := DisplayEnv{}

	env.IsTerminal = term.IsTerminal(int(os.Stdout.Fd()))
	if env.IsTerminal {
		env.Cols, env.Rows, _ = term.GetSize(int(os.Stdout.Fd()))
	}
	env.TerminalType = os.Getenv("TERM")

	env.IsSSH = os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_TTY") != ""
	env.IsMosh = os.Getenv("MOSH") != "" || strings.Contains(os.Getenv("TERM"), "mosh")
	env.IsScreen = os.Getenv("STY") != "" || os.Getenv("TMUX") != ""

	env.detectPlatformDisplay()

	env.Mode = env.autoDetectDisplayMode()

	return env
}

// autoDetectDisplayMode determines the display mode from the environment.
func (e *DisplayEnv) autoDetectDisplayMode() DisplayMode {
	if !e.IsTerminal && !e.HasDisplay {
		return DisplayModeHeadless
	}
	// TERM=dumb forces CLI mode: no TUI, no ANSI escapes.
	if e.TerminalType == "dumb" {
		return DisplayModeCLI
	}
	if e.HasDisplay && !e.IsSSH && !e.IsMosh {
		return DisplayModeGUI
	}
	if e.IsTerminal {
		return DisplayModeTUI
	}
	return DisplayModeCLI
}

// IsDumbTerminal reports whether the terminal is dumb (no ANSI support).
func (e *DisplayEnv) IsDumbTerminal() bool {
	return e.TerminalType == "dumb"
}

// IsAutoDetectDisplayModeGUI reports whether the detected mode is GUI.
func (e DisplayEnv) IsAutoDetectDisplayModeGUI() bool { return e.Mode == DisplayModeGUI }

// IsAutoDetectDisplayModeTUI reports whether the detected mode is TUI.
func (e DisplayEnv) IsAutoDetectDisplayModeTUI() bool { return e.Mode == DisplayModeTUI }

// IsAutoDetectDisplayModeCLI reports whether the detected mode is CLI.
func (e DisplayEnv) IsAutoDetectDisplayModeCLI() bool { return e.Mode == DisplayModeCLI }

// IsAutoDetectDisplayModeHeadless reports whether the detected mode is headless.
func (e DisplayEnv) IsAutoDetectDisplayModeHeadless() bool { return e.Mode == DisplayModeHeadless }

// CanUseANSI reports whether ANSI escapes (cursor movement, clear screen,
// colors) may be used. NO_COLOR users typically want plain output, so it is
// respected here too.
func CanUseANSI(env *DisplayEnv) bool {
	if env.IsDumbTerminal() {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return env.IsTerminal
}
