package main

import (
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(*options) bool
	}{
		{"empty", nil, false, func(o *options) bool { return !o.showHelp && !o.showVersion }},
		{"short help", []string{"-h"}, false, func(o *options) bool { return o.showHelp }},
		{"long help", []string{"--help"}, false, func(o *options) bool { return o.showHelp }},
		{"short version", []string{"-v"}, false, func(o *options) bool { return o.showVersion }},
		{"status", []string{"--status"}, false, func(o *options) bool { return o.status }},
		{"flag space value", []string{"--port", "8080"}, false, func(o *options) bool { return o.port == "8080" }},
		{"flag equals value", []string{"--port=8080"}, false, func(o *options) bool { return o.port == "8080" }},
		{"mode equals", []string{"--mode=development"}, false, func(o *options) bool { return o.mode == "development" }},
		{"multiple flags", []string{"--config", "/tmp/c", "--debug", "--color=no"}, false, func(o *options) bool {
			return o.configDir == "/tmp/c" && o.debug && o.color == "no"
		}},
		{"shell consumes rest", []string{"--shell", "completions", "bash"}, false, func(o *options) bool {
			return len(o.shellCmd) == 2 && o.shellCmd[0] == "completions" && o.shellCmd[1] == "bash"
		}},
		{"service install", []string{"--service", "--install"}, false, func(o *options) bool {
			return len(o.serviceCmd) == 1 && o.serviceCmd[0] == "--install"
		}},
		{"maintenance update branch", []string{"--maintenance", "update", "branch", "stable"}, false, func(o *options) bool {
			return len(o.maintenanceCmd) == 3 && o.maintenanceCmd[2] == "stable"
		}},
		{"update bare", []string{"--update"}, false, func(o *options) bool {
			return o.updateCmd != nil && len(o.updateCmd) == 0
		}},
		{"missing value", []string{"--port"}, true, nil},
		{"unknown flag", []string{"--bogus"}, true, nil},
		{"unknown short flag", []string{"-x"}, true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o, err := parseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseArgs(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if err == nil && tt.check != nil && !tt.check(o) {
				t.Errorf("parseArgs(%v) parsed options do not match: %+v", tt.args, o)
			}
		})
	}
}

func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*options)
		wantErr bool
	}{
		{"defaults ok", func(o *options) {}, false},
		{"mode production", func(o *options) { o.mode = "production" }, false},
		{"mode dev alias", func(o *options) { o.mode = "dev" }, false},
		{"mode invalid", func(o *options) { o.mode = "staging" }, true},
		{"color auto", func(o *options) { o.color = "auto" }, false},
		{"color invalid", func(o *options) { o.color = "maybe" }, true},
		{"baseurl ok", func(o *options) { o.baseURL = "/app" }, false},
		{"baseurl no slash", func(o *options) { o.baseURL = "app" }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &options{}
			tt.mutate(o)
			if err := o.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPrintMainHelp(t *testing.T) {
	var b strings.Builder
	printMainHelp(&b, "renamed-binary")
	out := b.String()
	wantFragments := []string{
		"renamed-binary",
		"Usage:",
		"Information:",
		"Shell Integration:",
		"Server Configuration:",
		"Service Management:",
		"-h, --help",
		"--shell completions [SHELL]",
		"--mode {production|development}",
		"--port PORT",
		"--maintenance CMD",
		"Run 'renamed-binary <command> help' for detailed help on any command.",
	}
	for _, f := range wantFragments {
		if !strings.Contains(out, f) {
			t.Errorf("main help missing %q", f)
		}
	}
	if strings.Contains(out, "tabssh <command>") {
		t.Error("help must show the actual binary name, not the compiled project name")
	}
}

func TestPrintVersion(t *testing.T) {
	var b strings.Builder
	printVersion(&b, "mybin")
	if !strings.HasPrefix(b.String(), "mybin ") {
		t.Errorf("version output %q must start with the binary name", b.String())
	}
}

func TestSubcommandHelpOutputs(t *testing.T) {
	var svc strings.Builder
	printServiceHelp(&svc, true, true, true, 1234)
	for _, f := range []string{"start", "stop", "restart", "reload", "--install", "--disable", "--uninstall", "Service:    installed", "State:      running", "Auto-start: enabled", "PID:        1234"} {
		if !strings.Contains(svc.String(), f) {
			t.Errorf("service help missing %q", f)
		}
	}

	var maint strings.Builder
	printMaintenanceHelp(&maint, "bin")
	for _, f := range []string{"backup [file]", "restore <file>", "update [cmd]", "mode <mode>", "setup", "bin --maintenance backup"} {
		if !strings.Contains(maint.String(), f) {
			t.Errorf("maintenance help missing %q", f)
		}
	}

	var sh strings.Builder
	printShellHelp(&sh, "bin")
	for _, f := range []string{"completions [SHELL]", "init [SHELL]", "bash, zsh, fish, sh, dash, ksh, powershell, pwsh"} {
		if !strings.Contains(sh.String(), f) {
			t.Errorf("shell help missing %q", f)
		}
	}

	var upd strings.Builder
	printUpdateHelp(&upd, "bin", "stable")
	for _, f := range []string{"check", "yes", "branch <name>", "Branch:   stable"} {
		if !strings.Contains(upd.String(), f) {
			t.Errorf("update help missing %q", f)
		}
	}
}

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		env  string
		want string
	}{
		{"explicit wins", "zsh", "/bin/bash", "zsh"},
		{"env basename", "", "/usr/bin/fish", "fish"},
		{"case folded", "BASH", "", "bash"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SHELL", tt.env)
			if got := detectShell(tt.arg); got != tt.want {
				t.Errorf("detectShell(%q) = %q, want %q", tt.arg, got, tt.want)
			}
		})
	}
}

func TestCompletionScripts(t *testing.T) {
	for _, shell := range supportedShells {
		t.Run(shell, func(t *testing.T) {
			script := completionScript(shell, "tabssh")
			if script == "" {
				t.Fatalf("no completion script for %s", shell)
			}
			if !strings.Contains(script, "tabssh") {
				t.Errorf("%s completion does not reference the binary name", shell)
			}
			init := initCommand(shell, "tabssh")
			if init == "" {
				t.Fatalf("no init command for %s", shell)
			}
		})
	}
	if !isSupportedShell("bash") || isSupportedShell("csh") {
		t.Error("isSupportedShell membership incorrect")
	}
}

func TestSanitizeFuncName(t *testing.T) {
	tests := []struct{ in, want string }{
		{"tabssh", "tabssh"},
		{"my-app.2", "my_app_2"},
		{"OK_name", "OK_name"},
	}
	for _, tt := range tests {
		if got := sanitizeFuncName(tt.in); got != tt.want {
			t.Errorf("sanitizeFuncName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
