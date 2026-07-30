//go:build windows

// Windows daemonization stub, per AI.md PART 8 (Daemonization). Windows
// does not support Unix daemonization; a Windows Service is the supported
// mechanism, so --daemon is ignored with a warning.
package runenv

import (
	"fmt"
	"os"
)

// Daemonize warns that --daemon is unsupported on Windows and continues in
// the foreground.
func Daemonize() error {
	fmt.Fprintln(os.Stderr, "Warning: --daemon is not supported on Windows")
	fmt.Fprintln(os.Stderr, "Use --service --install && --service start for Windows Service")
	return nil
}
