// Package display implements display-environment detection shared by all
// TabSSH Web binaries, per AI.md PART 7 (Display Environment Detection).
// This file defines the DisplayMode type and its helpers.
package display

// DisplayMode is the UI display mode (NOT the application mode).
type DisplayMode int

const (
	// DisplayModeHeadless means no display and no TTY (daemon, service, cron).
	DisplayModeHeadless DisplayMode = iota
	// DisplayModeCLI means command-line only (piped output or command given).
	DisplayModeCLI
	// DisplayModeTUI means an interactive terminal UI is possible.
	DisplayModeTUI
	// DisplayModeGUI means a native graphical display is available.
	DisplayModeGUI
)

// String returns the lowercase display-mode name.
func (m DisplayMode) String() string {
	switch m {
	case DisplayModeCLI:
		return "cli"
	case DisplayModeTUI:
		return "tui"
	case DisplayModeGUI:
		return "gui"
	default:
		return "headless"
	}
}
