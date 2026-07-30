// CLI flag parsing for the server binary, per AI.md PART 8 (Server Binary
// Commands). Every flag supports both --flag=value and --flag value; only
// -h and -v have short forms. Unknown flags are usage errors (exit 64).
package main

import (
	"fmt"
	"strings"
)

// exitUsage is the standardized exit code for bad arguments.
const exitUsage = 64

// options holds the parsed command line.
type options struct {
	showHelp    bool
	showVersion bool
	status      bool
	daemon      bool
	debug       bool
	mode        string
	configDir   string
	dataDir     string
	cacheDir    string
	logDir      string
	backupDir   string
	pidFile     string
	address     string
	port        string
	baseURL     string
	color       string
	lang        string
	// Command flags: nil when absent, otherwise the subcommand words.
	shellCmd       []string
	serviceCmd     []string
	maintenanceCmd []string
	updateCmd      []string
}

// valueFlags maps long value-taking flags to their options field setter.
func (o *options) valueFlagSetter(name string) func(string) {
	switch name {
	case "--mode":
		return func(v string) { o.mode = v }
	case "--config":
		return func(v string) { o.configDir = v }
	case "--data":
		return func(v string) { o.dataDir = v }
	case "--cache":
		return func(v string) { o.cacheDir = v }
	case "--log":
		return func(v string) { o.logDir = v }
	case "--backup":
		return func(v string) { o.backupDir = v }
	case "--pid":
		return func(v string) { o.pidFile = v }
	case "--address":
		return func(v string) { o.address = v }
	case "--port":
		return func(v string) { o.port = v }
	case "--baseurl":
		return func(v string) { o.baseURL = v }
	case "--color":
		return func(v string) { o.color = v }
	case "--lang":
		return func(v string) { o.lang = v }
	}
	return nil
}

// parseArgs parses the server command line. Command flags (--shell,
// --service, --maintenance, --update) consume all remaining arguments as
// their subcommand words, matching the documented invocation patterns.
func parseArgs(args []string) (*options, error) {
	o := &options{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		name := arg
		inlineVal := ""
		hasInline := false
		if idx := strings.Index(arg, "="); idx > 0 && strings.HasPrefix(arg, "--") {
			name = arg[:idx]
			inlineVal = arg[idx+1:]
			hasInline = true
		}

		switch name {
		case "-h", "--help":
			o.showHelp = true
			continue
		case "-v", "--version":
			o.showVersion = true
			continue
		case "--status":
			o.status = true
			continue
		case "--daemon":
			o.daemon = true
			continue
		case "--debug":
			o.debug = true
			continue
		case "--shell", "--service", "--maintenance", "--update":
			rest := []string{}
			if hasInline {
				rest = append(rest, inlineVal)
			}
			rest = append(rest, args[i+1:]...)
			switch name {
			case "--shell":
				o.shellCmd = rest
			case "--service":
				o.serviceCmd = rest
			case "--maintenance":
				o.maintenanceCmd = rest
			case "--update":
				o.updateCmd = rest
			}
			return o, nil
		}

		if set := o.valueFlagSetter(name); set != nil {
			if hasInline {
				set(inlineVal)
				continue
			}
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag %s requires a value", name)
			}
			i++
			set(args[i])
			continue
		}

		return nil, fmt.Errorf("unknown flag: %s", arg)
	}
	return o, nil
}

// validate checks parsed flag values that have a fixed value set.
func (o *options) validate() error {
	switch o.mode {
	case "", "production", "prod", "development", "dev", "devel", "debug":
	default:
		return fmt.Errorf("invalid --mode %q (use production or development)", o.mode)
	}
	switch o.color {
	case "", "auto", "yes", "no":
	default:
		return fmt.Errorf("invalid --color %q (use auto, yes, or no)", o.color)
	}
	if o.baseURL != "" && !strings.HasPrefix(o.baseURL, "/") {
		return fmt.Errorf("invalid --baseurl %q (must start with /)", o.baseURL)
	}
	return nil
}
