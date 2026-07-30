// Package version holds TabSSH Web's compiled-in build identity, per AI.md
// PART 7 (Binary Requirements) and PART 13 (Versioning). The values are set
// at build time via -ldflags "-X"; the version resolution order at build
// time is release.txt > git tag > "dev".
package version

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	// ProjectName is the compiled project name. It never changes even if
	// the binary file is renamed; it is used for User-Agent headers and
	// default path construction, never for display of the binary name.
	ProjectName = "tabssh"
	// Version is the build version, injected at build time. It follows
	// SemVer for stable releases, YYYYMMDDHHMMSS-beta for beta, and
	// YYYYMMDDHHMMSS for daily builds. Fallback is "dev".
	Version = "dev"
	// Build is the short VCS revision, injected at build time.
	Build = ""
)

// BinaryName returns the actual (possibly renamed) executable name for use
// in --help, --version, and error messages.
func BinaryName() string {
	return filepath.Base(os.Args[0])
}

// UserAgent returns the server's User-Agent string. It always uses the
// compiled ProjectName, never the on-disk binary name.
func UserAgent() string {
	return ProjectName + "/" + Version
}

// GoVersion returns the Go runtime version the binary was built with.
func GoVersion() string {
	return runtime.Version()
}
