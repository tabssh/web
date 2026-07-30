// Color and emoji output policy, per AI.md PART 8 (NO_COLOR). Priority for
// color: --color CLI flag > config output.color > NO_COLOR env var > auto.
// NO_COLOR (https://no-color.org/) disables colors AND emojis, but never
// bold/underline or box-drawing characters.
package display

import "os"

// ColorEnabled resolves whether colored output is enabled. flagVal is the
// --color CLI flag value ("auto", "yes", "no", or "" when not given);
// configVal is the config file output.color value with the same options.
// When both defer to auto, NO_COLOR wins, then TTY plus TERM!=dumb.
func ColorEnabled(flagVal, configVal string, env DisplayEnv) bool {
	switch flagVal {
	case "yes":
		return true
	case "no":
		return false
	}
	switch configVal {
	case "yes":
		return true
	case "no":
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return env.IsTerminal && !env.IsDumbTerminal()
}

// EmojiEnabled resolves whether emoji output is enabled. An explicit config
// override (output.emoji) wins; otherwise NO_COLOR disables emojis, dumb
// terminals disable emojis, and the default is enabled.
func EmojiEnabled(configOverride *bool, env DisplayEnv) bool {
	if configOverride != nil {
		return *configOverride
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if env.IsDumbTerminal() {
		return false
	}
	return true
}
