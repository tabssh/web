package banner

import (
	"strings"
	"testing"

	"github.com/tabssh/web/src/common/terminal"
)

func testConfig() BannerConfig {
	return BannerConfig{
		AppName: "TabSSH Web",
		Version: "1.0.0",
		AppMode: "production",
		URLs:    []string{"http://example.com:64123/", "http://abc.onion/"},
	}
}

func TestRenderersIncludeCoreFields(t *testing.T) {
	cfg := testConfig()
	tests := []struct {
		name   string
		render func() string
	}{
		{"full", func() string { return renderFull(cfg, terminal.TerminalSize{Cols: 80, Rows: 24}) }},
		{"compact", func() string { return renderCompact(cfg) }},
		{"minimal", func() string { return renderMinimal(cfg) }},
		{"micro", func() string { return renderMicro(cfg) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.render()
			if !strings.Contains(out, "TabSSH Web") {
				t.Errorf("%s banner missing app name:\n%s", tt.name, out)
			}
			if !strings.Contains(out, "1.0.0") {
				t.Errorf("%s banner missing version:\n%s", tt.name, out)
			}
			if !strings.Contains(out, "http://example.com:64123/") {
				t.Errorf("%s banner missing primary URL:\n%s", tt.name, out)
			}
		})
	}
}

func TestSetupTokenShownOnlyWhenRequested(t *testing.T) {
	cfg := testConfig()
	cfg.SetupToken = "deadbeefdeadbeefdeadbeefdeadbeef"

	tests := []struct {
		name      string
		showSetup bool
		want      bool
	}{
		{"hidden without ShowSetup", false, false},
		{"shown with ShowSetup", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.ShowSetup = tt.showSetup
			for _, out := range []string{
				renderFull(cfg, terminal.TerminalSize{Cols: 80, Rows: 24}),
				renderCompact(cfg),
				renderMinimal(cfg),
				renderMicro(cfg),
			} {
				got := strings.Contains(out, cfg.SetupToken)
				if got != tt.want {
					t.Errorf("token visibility = %v, want %v in:\n%s", got, tt.want, out)
				}
			}
		})
	}
}

func TestModeLine(t *testing.T) {
	tests := []struct {
		name  string
		mode  string
		debug bool
		want  string
	}{
		{"production", "production", false, "Mode: production"},
		{"debugging", "development", true, "Mode: development [debugging]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := BannerConfig{AppMode: tt.mode, Debug: tt.debug}
			if got := modeLine(cfg); got != tt.want {
				t.Errorf("modeLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestASCIIArt(t *testing.T) {
	tests := []struct {
		name    string
		unicode bool
		char    string
	}{
		{"unicode rule", true, "─"},
		{"ascii rule", false, "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := ASCIIArt("tabssh", tt.unicode, 40)
			if !strings.Contains(out, "TABSSH") {
				t.Errorf("ASCIIArt missing title: %q", out)
			}
			if !strings.Contains(out, tt.char) {
				t.Errorf("ASCIIArt missing rule char %q: %q", tt.char, out)
			}
			if !strings.HasSuffix(out, "\n") {
				t.Errorf("ASCIIArt must end with newline: %q", out)
			}
		})
	}
}

func TestPadBoxWidth(t *testing.T) {
	line := padBox("hello", 40, "|")
	if len(strings.TrimSuffix(line, "\n")) != 40 {
		t.Errorf("padBox row width = %d, want 40: %q", len(line)-1, line)
	}
}
