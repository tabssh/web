// Help output for the server binary, matching the exact templates in
// AI.md PART 8 (Standard Server Flags help template and the Service,
// Maintenance, Shell, and Update help outputs). The actual binary name
// (possibly renamed) is shown; the compiled project name is used only
// for User-Agent and paths.
package main

import (
	"fmt"
	"io"

	"github.com/tabssh/web/src/common/version"
)

// projectDescription is the one-line description shown in --help,
// sourced from IDEA.md.
const projectDescription = "Browser-based tabbed SSH client with encrypted sync and device pairing"

// printMainHelp writes the top-level --help output.
func printMainHelp(w io.Writer, binName string) {
	fmt.Fprintf(w, "%s %s - %s\n", binName, version.Version, projectDescription)
	fmt.Fprintf(w, `
Usage:
  %[1]s [flags]

Information:
-h, --help                             - Show help (--help for any command shows its help)
-v, --version                          - Show version
--status                               - Show server status and health

Shell Integration:
--shell completions [SHELL]            - Print shell completions
--shell init [SHELL]                   - Print shell init command
--shell help                           - Show shell help

Server Configuration:
--mode {production|development}        - Application mode (default: production)
--config DIR                           - Config directory
--data DIR                             - Data directory
--cache DIR                            - Cache directory
--log DIR                              - Log directory
--backup DIR                           - Backup directory
--pid FILE                             - PID file path
--address ADDR                         - Listen address (default: 0.0.0.0)
--port PORT                            - Listen port (default: random 64xxx, 80 in container)
--baseurl PATH                         - URL path prefix (default: /)
--daemon                               - Run as daemon (detach from terminal)
--debug                                - Enable debug mode
--color {auto|yes|no}                  - Color output (default: auto)
--lang CODE                            - Language for output (default: auto)

Service Management:
--service CMD                          - Service management (run --service help for details)
--maintenance CMD                      - Maintenance operations (run --maintenance help for details)
--update [CMD]                         - Check/perform updates (run --update help for details)

Run '%[1]s <command> help' for detailed help on any command.
`, binName)
}

// printVersion writes the --version output.
func printVersion(w io.Writer, binName string) {
	if version.Build != "" {
		fmt.Fprintf(w, "%s %s (%s)\n", binName, version.Version, version.Build)
		return
	}
	fmt.Fprintf(w, "%s %s\n", binName, version.Version)
}

// printServiceHelp writes the --service help output, including the
// current service status block.
func printServiceHelp(w io.Writer, installed, running, autoStart bool, pid int) {
	fmt.Fprint(w, `Service management commands:

start                                 - Start the service
stop                                  - Stop the service
restart                               - Restart the service
reload                                - Reload configuration without restart
--install                              - Install, enable, and start service
--disable                              - Stop and disable service (keeps data)
--uninstall                            - Stop, disable, and remove everything (keeps binary)

Current status:
`)
	svc := "not installed"
	if installed {
		svc = "installed"
	}
	state := "stopped"
	if running {
		state = "running"
	} else if installed && !autoStart {
		state = "disabled"
	}
	auto := "disabled"
	if autoStart {
		auto = "enabled"
	}
	fmt.Fprintf(w, "  Service:    %s\n", svc)
	fmt.Fprintf(w, "  State:      %s\n", state)
	fmt.Fprintf(w, "  Auto-start: %s\n", auto)
	if running && pid > 0 {
		fmt.Fprintf(w, "  PID:        %d\n", pid)
	}
}

// printMaintenanceHelp writes the --maintenance help output.
func printMaintenanceHelp(w io.Writer, binName string) {
	fmt.Fprintf(w, `Maintenance commands:

backup [file]                         - Create backup of all data
                                        Default: {backup_dir}/%[1]s-{timestamp}.tar.gz
restore <file>                        - Restore from backup file
                                        Stops server, restores data, restarts server
update [cmd]                          - Manage updates
                                        check         - Check for available updates
                                        yes           - Download and install update
                                        branch <name> - Switch update branch (stable|beta|daily)
mode <mode>                           - Set application mode
                                        production    - Normal operation (default)
                                        development   - Relaxed dev defaults (does NOT enable debug; use --debug/DEBUG=true)
setup                                 - Run interactive setup wizard
                                        Creates primary Server Admin, configures server

Examples:
  %[1]s --maintenance backup
  %[1]s --maintenance backup /path/to/backup.tar.gz
  %[1]s --maintenance restore /path/to/backup.tar.gz
  %[1]s --maintenance update check
  %[1]s --maintenance update yes
  %[1]s --maintenance mode development
  %[1]s --maintenance setup
`, binName)
}

// printShellHelp writes the --shell help output.
func printShellHelp(w io.Writer, binName string) {
	fmt.Fprintf(w, `Shell integration commands:

completions [SHELL]                   - Print shell completion script
                                        Auto-detects shell if SHELL omitted
                                        Supported: bash, zsh, fish, sh, dash, ksh, powershell, pwsh
init [SHELL]                          - Print shell init command for eval
                                        Auto-detects shell if SHELL omitted

Usage:
  # Add to shell profile for persistent completions
  # bash
  %[1]s --shell init >> ~/.bashrc
  # zsh
  %[1]s --shell init >> ~/.zshrc
  # fish
  %[1]s --shell init >> ~/.config/fish/config.fish

  # Or eval directly for current session
  eval "$(%[1]s --shell init)"

  # Generate completion script only
  %[1]s --shell completions bash > /etc/bash_completion.d/%[1]s
`, binName)
}

// printUpdateHelp writes the --update help output.
func printUpdateHelp(w io.Writer, binName, branch string) {
	fmt.Fprintf(w, `Update management:

check                                 - Check for available updates
                                        Compares current version with latest release
yes                                   - Download and install update
                                        Downloads latest release, replaces binary, restarts
branch <name>                         - Switch update branch
                                        stable - Stable releases (default)
                                        beta   - Beta/preview releases
                                        daily  - Daily builds (development)

Examples:
  %[1]s --update check
  %[1]s --update yes
  %[1]s --update branch beta

Current:
  Version:  %[2]s
  Branch:   %[3]s
`, binName, version.Version, branch)
}
