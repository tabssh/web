// Package mode implements TabSSH Web's application-mode and debug-flag
// detection, per AI.md PART 6 (Application Modes). Mode/debug detection
// priority: CLI flag > environment variable > default. Debug mode affects
// verbosity and diagnostics ONLY — it never bypasses authentication or any
// other security check, in any mode, including production.
package mode

import (
	"os"
	"runtime"
	"strings"

	"github.com/tabssh/web/src/config"
)

var (
	currentMode  = Production
	debugEnabled = false
)

// AppMode is the application's operating mode: production or development.
type AppMode int

const (
	// Production is the default mode: minimal logging, caching enabled,
	// rate limiting enforced, debug endpoints disabled.
	Production AppMode = iota
	// Development enables verbose logging and disables caching, but does
	// NOT enable debug endpoints on its own — that requires --debug.
	Development
)

// String returns the lowercase mode name.
func (m AppMode) String() string {
	switch m {
	case Development:
		return "development"
	default:
		return "production"
	}
}

// SetAppMode sets the application mode from a string. "debug" is an alias
// for development mode with debug enabled; an explicit --debug flag or
// DEBUG env var processed afterward still wins over the alias.
func SetAppMode(m string) {
	switch strings.ToLower(m) {
	case "dev", "devel", "development":
		currentMode = Development
	case "debug":
		currentMode = Development
		SetDebugEnabled(true)
	default:
		currentMode = Production
	}
	updateAppModeProfilingSettings()
}

// SetDebugEnabled enables or disables debug mode. Debug mode never bypasses
// authentication or any other security check in any mode.
func SetDebugEnabled(enabled bool) {
	debugEnabled = enabled
	updateAppModeProfilingSettings()
}

// updateAppModeProfilingSettings toggles Go runtime profiling based on the
// current debug flag.
func updateAppModeProfilingSettings() {
	if debugEnabled {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(1)
	} else {
		runtime.SetBlockProfileRate(0)
		runtime.SetMutexProfileFraction(0)
	}
}

// GetCurrentAppMode returns the current application mode.
func GetCurrentAppMode() AppMode {
	return currentMode
}

// IsAppModeDev reports whether the application is in development mode.
func IsAppModeDev() bool {
	return currentMode == Development
}

// IsAppModeProd reports whether the application is in production mode.
func IsAppModeProd() bool {
	return currentMode == Production
}

// IsDebugEnabled reports whether debug mode (--debug or DEBUG=true) is on.
func IsDebugEnabled() bool {
	return debugEnabled
}

// GetAppModeString returns the mode name, with a "[debugging]" suffix when
// debug mode is enabled.
func GetAppModeString() string {
	s := currentMode.String()
	if debugEnabled {
		s += " [debugging]"
	}
	return s
}

// FromEnv sets mode and debug from the MODE and DEBUG environment variables.
// An explicitly set DEBUG env var always wins over the MODE=debug alias:
// MODE=debug DEBUG=false runs development mode with debug off.
func FromEnv() {
	if m := os.Getenv("MODE"); m != "" {
		SetAppMode(m)
	}
	if d, ok := os.LookupEnv("DEBUG"); ok && d != "" {
		SetDebugEnabled(config.IsTruthy(d))
	}
}
