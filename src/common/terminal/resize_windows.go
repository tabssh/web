//go:build windows

// Polling-based terminal resize watching for Windows, per AI.md PART 7
// (Resize Handling). Windows has no SIGWINCH, so the size is polled.
package terminal

import (
	"context"
	"time"
)

// WatchResize invokes onResize with the new terminal size whenever the size
// changes, polling every 500ms until ctx is cancelled. It blocks and is
// intended to run in its own goroutine.
func WatchResize(ctx context.Context, onResize func(TerminalSize)) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastCols, lastRows int

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			size := GetTerminalSize()
			if size.Cols != lastCols || size.Rows != lastRows {
				lastCols, lastRows = size.Cols, size.Rows
				onResize(size)
			}
		}
	}
}
