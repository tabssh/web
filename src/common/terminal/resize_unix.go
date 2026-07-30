//go:build !windows

// SIGWINCH-based terminal resize watching for Unix platforms, per AI.md
// PART 7 (Resize Handling).
package terminal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WatchResize invokes onResize with the new terminal size every time the
// terminal is resized (SIGWINCH), until ctx is cancelled. It blocks and is
// intended to run in its own goroutine.
func WatchResize(ctx context.Context, onResize func(TerminalSize)) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	defer signal.Stop(sigCh)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigCh:
			onResize(GetTerminalSize())
		}
	}
}
