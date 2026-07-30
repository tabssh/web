package display

import "testing"

func TestAutoDetectDisplayMode(t *testing.T) {
	tests := []struct {
		name string
		env  DisplayEnv
		want DisplayMode
	}{
		{"no terminal no display", DisplayEnv{IsTerminal: false, HasDisplay: false}, DisplayModeHeadless},
		{"dumb terminal forces cli", DisplayEnv{IsTerminal: true, TerminalType: "dumb"}, DisplayModeCLI},
		{"dumb with display forces cli", DisplayEnv{IsTerminal: true, HasDisplay: true, TerminalType: "dumb"}, DisplayModeCLI},
		{"local display", DisplayEnv{IsTerminal: true, HasDisplay: true, TerminalType: "xterm"}, DisplayModeGUI},
		{"display no terminal", DisplayEnv{IsTerminal: false, HasDisplay: true, TerminalType: "xterm"}, DisplayModeGUI},
		{"ssh with display env", DisplayEnv{IsTerminal: true, HasDisplay: true, IsSSH: true, TerminalType: "xterm"}, DisplayModeTUI},
		{"mosh with display env", DisplayEnv{IsTerminal: true, HasDisplay: true, IsMosh: true, TerminalType: "xterm-mosh"}, DisplayModeTUI},
		{"plain terminal", DisplayEnv{IsTerminal: true, TerminalType: "xterm"}, DisplayModeTUI},
		{"ssh no terminal with display", DisplayEnv{IsTerminal: false, HasDisplay: true, IsSSH: true, TerminalType: "xterm"}, DisplayModeCLI},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.autoDetectDisplayMode(); got != tt.want {
				t.Errorf("autoDetectDisplayMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisplayModeString(t *testing.T) {
	tests := []struct {
		mode DisplayMode
		want string
	}{
		{DisplayModeHeadless, "headless"},
		{DisplayModeCLI, "cli"},
		{DisplayModeTUI, "tui"},
		{DisplayModeGUI, "gui"},
	}
	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.want {
			t.Errorf("DisplayMode(%d).String() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestCanUseANSI(t *testing.T) {
	tests := []struct {
		name    string
		env     DisplayEnv
		noColor string
		want    bool
	}{
		{"dumb terminal", DisplayEnv{IsTerminal: true, TerminalType: "dumb"}, "", false},
		{"no_color set", DisplayEnv{IsTerminal: true, TerminalType: "xterm"}, "1", false},
		{"normal terminal", DisplayEnv{IsTerminal: true, TerminalType: "xterm"}, "", true},
		{"not a terminal", DisplayEnv{IsTerminal: false, TerminalType: "xterm"}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tt.noColor)
			if got := CanUseANSI(&tt.env); got != tt.want {
				t.Errorf("CanUseANSI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorEnabled(t *testing.T) {
	tty := DisplayEnv{IsTerminal: true, TerminalType: "xterm"}
	dumb := DisplayEnv{IsTerminal: true, TerminalType: "dumb"}
	tests := []struct {
		name      string
		flagVal   string
		configVal string
		noColor   string
		env       DisplayEnv
		want      bool
	}{
		{"flag yes wins over no_color", "yes", "auto", "1", dumb, true},
		{"flag no wins", "no", "yes", "", tty, false},
		{"config yes wins over no_color", "auto", "yes", "1", tty, true},
		{"config no wins", "", "no", "", tty, false},
		{"no_color disables", "auto", "auto", "1", tty, false},
		{"auto tty enabled", "auto", "auto", "", tty, true},
		{"auto dumb disabled", "auto", "auto", "", dumb, false},
		{"auto non-tty disabled", "", "", "", DisplayEnv{IsTerminal: false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tt.noColor)
			if got := ColorEnabled(tt.flagVal, tt.configVal, tt.env); got != tt.want {
				t.Errorf("ColorEnabled(%q, %q) = %v, want %v", tt.flagVal, tt.configVal, got, tt.want)
			}
		})
	}
}

func TestEmojiEnabled(t *testing.T) {
	yes, no := true, false
	tty := DisplayEnv{IsTerminal: true, TerminalType: "xterm"}
	dumb := DisplayEnv{IsTerminal: true, TerminalType: "dumb"}
	tests := []struct {
		name     string
		override *bool
		noColor  string
		env      DisplayEnv
		want     bool
	}{
		{"config true wins over no_color", &yes, "1", dumb, true},
		{"config false wins", &no, "", tty, false},
		{"no_color disables", nil, "1", tty, false},
		{"dumb disables", nil, "", dumb, false},
		{"default on", nil, "", tty, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("NO_COLOR", tt.noColor)
			if got := EmojiEnabled(tt.override, tt.env); got != tt.want {
				t.Errorf("EmojiEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
