// ASCII art generation for the startup banner, per AI.md PART 7 (Banner
// Package). Art is only shown at SizeModeStandard and above; the caller is
// responsible for that check.
package banner

import (
	"strings"

	"github.com/tabssh/web/src/common/terminal"
)

// ASCIIArt renders the application name inside a horizontal rule header
// sized to width. Unicode selects box-drawing characters; false uses plain
// ASCII suitable for dumb terminals.
func ASCIIArt(name string, unicode bool, width int) string {
	sym := terminal.GetSymbols(unicode)
	if width < 10 {
		width = 10
	}
	title := " " + strings.ToUpper(name) + " "
	if len(title) > width {
		title = title[:width]
	}
	pad := width - len(title)
	left := pad / 2
	right := pad - left

	var b strings.Builder
	b.WriteString(strings.Repeat(sym.Horiz, left))
	b.WriteString(title)
	b.WriteString(strings.Repeat(sym.Horiz, right))
	b.WriteString("\n")
	return b.String()
}
