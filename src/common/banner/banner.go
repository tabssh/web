// Package banner prints the responsive startup banner shared by TabSSH Web
// binaries, per AI.md PART 7 (Banner Package) and PART 8 (startup output).
// The banner adapts to the terminal size mode: full, compact, minimal, or
// micro.
package banner

import (
	"fmt"
	"strings"

	"github.com/tabssh/web/src/common/terminal"
)

// BannerConfig configures the startup banner content.
type BannerConfig struct {
	AppName string
	Version string
	// AppMode is "production" or "development", with a debug suffix added
	// separately via Debug.
	AppMode string
	Debug   bool
	URLs    []string
	// ShowSetup shows the one-time setup token (server only, first run).
	ShowSetup  bool
	SetupToken string
	// Unicode selects Unicode box drawing; false falls back to ASCII.
	Unicode bool
}

// PrintStartupBanner prints the startup banner sized for the current
// terminal.
func PrintStartupBanner(cfg BannerConfig) {
	size := terminal.GetTerminalSize()

	switch {
	case size.Mode >= terminal.SizeModeStandard:
		fmt.Print(renderFull(cfg, size))
	case size.Mode >= terminal.SizeModeCompact:
		fmt.Print(renderCompact(cfg))
	case size.Mode >= terminal.SizeModeMinimal:
		fmt.Print(renderMinimal(cfg))
	default:
		fmt.Print(renderMicro(cfg))
	}
}

// modeLine builds the mode line shown after the header and before URLs.
func modeLine(cfg BannerConfig) string {
	mode := cfg.AppMode
	if cfg.Debug {
		mode += " [debugging]"
	}
	return "Mode: " + mode
}

// renderFull renders the full banner with ASCII art and a bordered layout
// for Standard and larger terminals.
func renderFull(cfg BannerConfig, size terminal.TerminalSize) string {
	sym := terminal.GetSymbols(cfg.Unicode)
	width := size.Cols
	if width > 78 {
		width = 78
	}

	var b strings.Builder
	b.WriteString(ASCIIArt(cfg.AppName, cfg.Unicode, width))
	b.WriteString(fmt.Sprintf("%s %s\n", cfg.AppName, cfg.Version))
	b.WriteString(modeLine(cfg) + "\n")
	for _, u := range cfg.URLs {
		b.WriteString(fmt.Sprintf("  %s %s\n", sym.Arrow, u))
	}
	if cfg.ShowSetup && cfg.SetupToken != "" {
		b.WriteString(setupBlock(cfg, width))
	}
	return b.String()
}

// renderCompact renders a borderless multi-line banner for Compact
// terminals.
func renderCompact(cfg BannerConfig) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s\n", cfg.AppName, cfg.Version))
	b.WriteString(modeLine(cfg) + "\n")
	for _, u := range cfg.URLs {
		b.WriteString("  " + u + "\n")
	}
	if cfg.ShowSetup && cfg.SetupToken != "" {
		b.WriteString("Setup token: " + cfg.SetupToken + "\n")
		b.WriteString("Shown once - save it now.\n")
	}
	return b.String()
}

// renderMinimal renders a short banner for Minimal terminals.
func renderMinimal(cfg BannerConfig) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s %s (%s)\n", cfg.AppName, cfg.Version, cfg.AppMode))
	for _, u := range cfg.URLs {
		b.WriteString(u + "\n")
	}
	if cfg.ShowSetup && cfg.SetupToken != "" {
		b.WriteString("Setup: " + cfg.SetupToken + "\n")
	}
	return b.String()
}

// renderMicro renders the smallest possible banner for Micro terminals,
// including phone SSH sessions.
func renderMicro(cfg BannerConfig) string {
	var b strings.Builder
	b.WriteString(cfg.AppName + " " + cfg.Version + "\n")
	if len(cfg.URLs) > 0 {
		b.WriteString(cfg.URLs[0] + "\n")
	}
	if cfg.ShowSetup && cfg.SetupToken != "" {
		b.WriteString(cfg.SetupToken + "\n")
	}
	return b.String()
}

// setupBlock renders the highlighted first-run setup token block.
func setupBlock(cfg BannerConfig, width int) string {
	sym := terminal.GetSymbols(cfg.Unicode)
	line := strings.Repeat(sym.Horiz, width-2)
	var b strings.Builder
	b.WriteString(sym.TopLeft + line + sym.TopRight + "\n")
	b.WriteString(padBox("First-run setup token (shown once):", width, sym.Vert))
	b.WriteString(padBox("  "+cfg.SetupToken, width, sym.Vert))
	b.WriteString(padBox("Save it now - it cannot be recovered.", width, sym.Vert))
	b.WriteString(sym.BotLeft + line + sym.BotRight + "\n")
	return b.String()
}

// padBox pads a line of text into a bordered box row of the given width.
func padBox(text string, width int, vert string) string {
	inner := width - 4
	if len(text) > inner {
		text = text[:inner]
	}
	return vert + " " + text + strings.Repeat(" ", inner-len(text)) + " " + vert + "\n"
}
