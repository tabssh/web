// Command dispatch for --status, --shell, --service, --maintenance, and
// --update per AI.md PART 8. Operations whose internals belong to later
// PARTs (service install/uninstall/disable -> PART 24/25; maintenance
// backup/restore/update/setup -> PART 17/22/23; update check/yes/branch
// -> PART 23) report a canonical error and exit non-zero rather than
// shipping stub behavior.
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/tabssh/web/src/common/version"
	"github.com/tabssh/web/src/config"
	"github.com/tabssh/web/src/paths"
	"github.com/tabssh/web/src/pid"
	"github.com/tabssh/web/src/runenv"
)

// runStatus implements --status: exit 0 when healthy/running, 1 otherwise.
func runStatus(opts *options, binName string) int {
	pidPath := resolvePIDPath(opts)
	running, p, err := pid.CheckPIDFile(pidPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
		return 1
	}
	if running {
		fmt.Printf("%s is running (PID %d)\n", binName, p)
		return 0
	}
	fmt.Printf("%s is not running\n", binName)
	return 1
}

// resolvePIDPath resolves the PID file path with the same flag > env >
// default priority used at startup.
func resolvePIDPath(opts *options) string {
	defaults := paths.Resolve()
	return runenv.GetPIDFile(opts.pidFile, defaults)
}

// runShell implements the --shell command family.
func runShell(args []string, binName string) int {
	if len(args) == 0 {
		printShellHelp(os.Stdout, binName)
		return 0
	}
	switch args[0] {
	case "help", "--help", "-h":
		printShellHelp(os.Stdout, binName)
		return 0
	case "completions":
		shell := detectShell(argAt(args, 1))
		if !isSupportedShell(shell) {
			fmt.Fprintf(os.Stderr, "ERROR: VALIDATION_FAILED: unsupported shell %q (supported: bash, zsh, fish, sh, dash, ksh, powershell, pwsh)\n", shell)
			return exitUsage
		}
		fmt.Print(completionScript(shell, binName))
		return 0
	case "init":
		shell := detectShell(argAt(args, 1))
		if !isSupportedShell(shell) {
			fmt.Fprintf(os.Stderr, "ERROR: VALIDATION_FAILED: unsupported shell %q (supported: bash, zsh, fish, sh, dash, ksh, powershell, pwsh)\n", shell)
			return exitUsage
		}
		fmt.Print(initCommand(shell, binName))
		return 0
	}
	fmt.Fprintf(os.Stderr, "ERROR: BAD_REQUEST: unknown shell command %q (run --shell help)\n", args[0])
	return exitUsage
}

// runService implements the --service command family. start/stop/
// restart/reload dispatch to the detected service manager; install/
// uninstall/disable require the PART 24/25 service subsystem.
func runService(args []string, opts *options, binName string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return serviceHelp(opts)
	}
	switch args[0] {
	case "start", "stop", "restart", "reload":
		return serviceControl(args[0])
	case "--install", "--uninstall", "--disable":
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: service %s is not available in this build (service unit management ships with the PART 24/25 service subsystem)\n", args[0])
		return 1
	}
	fmt.Fprintf(os.Stderr, "ERROR: BAD_REQUEST: unknown service command %q (run --service help)\n", args[0])
	return exitUsage
}

// serviceHelp prints the service help with live status detection.
func serviceHelp(opts *options) int {
	running, p, _ := pid.CheckPIDFile(resolvePIDPath(opts))
	installed, autoStart := detectServiceInstalled()
	printServiceHelp(os.Stdout, installed, running, autoStart, p)
	return 0
}

// detectServiceInstalled reports whether a service unit exists and is
// enabled for the compiled project name under the detected manager.
func detectServiceInstalled() (installed, autoStart bool) {
	name := version.ProjectName
	switch runenv.DetectServiceManager() {
	case "systemd":
		if err := exec.Command("systemctl", "cat", name).Run(); err == nil {
			installed = true
		}
		if err := exec.Command("systemctl", "is-enabled", "--quiet", name).Run(); err == nil {
			autoStart = true
		}
	case "sysv", "rcd":
		if _, err := os.Stat("/etc/init.d/" + name); err == nil {
			installed = true
			autoStart = true
		}
	case "runit":
		if _, err := os.Stat("/etc/service/" + name); err == nil {
			installed = true
			autoStart = true
		}
	case "launchd":
		if err := exec.Command("launchctl", "list", name).Run(); err == nil {
			installed = true
			autoStart = true
		}
	}
	return installed, autoStart
}

// serviceControl shells out to the detected service manager for
// start/stop/restart/reload and propagates its exit status.
func serviceControl(verb string) int {
	name := version.ProjectName
	var cmd *exec.Cmd
	switch runenv.DetectServiceManager() {
	case "systemd":
		cmd = exec.Command("systemctl", verb, name)
	case "sysv", "rcd":
		cmd = exec.Command("service", name, verb)
	case "runit":
		svVerb := verb
		if verb == "reload" {
			svVerb = "hup"
		}
		cmd = exec.Command("sv", svVerb, name)
	case "launchd":
		switch verb {
		case "start", "stop":
			cmd = exec.Command("launchctl", verb, name)
		default:
			fmt.Fprintf(os.Stderr, "ERROR: BAD_REQUEST: launchd does not support %q; use stop then start\n", verb)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: no supported service manager detected; run the binary directly or use --daemon\n")
		return 1
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
		return 1
	}
	return 0
}

// runMaintenance implements the --maintenance command family. mode is
// fully implemented; backup/restore/update/setup require later-PART
// subsystems (encrypted backup PART 22, updater PART 23, setup wizard
// PART 17) and report a canonical error.
func runMaintenance(args []string, opts *options, binName string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printMaintenanceHelp(os.Stdout, binName)
		return 0
	}
	switch args[0] {
	case "mode":
		return maintenanceMode(argAt(args, 1), opts)
	case "backup", "restore", "setup":
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: maintenance %s is not available in this build (ships with the PART 17/22 maintenance subsystem)\n", args[0])
		return 1
	case "update":
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: maintenance update is not available in this build (ships with the PART 23 update subsystem)\n")
		return 1
	}
	fmt.Fprintf(os.Stderr, "ERROR: BAD_REQUEST: unknown maintenance command %q (run --maintenance help)\n", args[0])
	return exitUsage
}

// maintenanceMode sets the persisted application mode in server.yml.
func maintenanceMode(newMode string, opts *options) int {
	var canonical string
	switch newMode {
	case "production", "prod":
		canonical = "production"
	case "development", "dev", "devel":
		canonical = "development"
	case "":
		fmt.Fprintf(os.Stderr, "ERROR: VALIDATION_FAILED: mode requires an argument (production or development)\n")
		return exitUsage
	default:
		fmt.Fprintf(os.Stderr, "ERROR: VALIDATION_FAILED: invalid mode %q (use production or development)\n", newMode)
		return exitUsage
	}

	configPath := resolveConfigPath(opts)
	cfg, _, err := config.LoadOrCreate(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: loading config: %v\n", err)
		return 2
	}
	cfg.Server.Mode = canonical
	if err := cfg.Save(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: saving config: %v\n", err)
		return 2
	}
	fmt.Printf("OK: application mode set to %s\n", canonical)
	return 0
}

// resolveConfigPath resolves {config_dir}/server.yml with the same
// flag > env > default priority used at startup.
func resolveConfigPath(opts *options) string {
	defaults := paths.Resolve()
	dir := runenv.GetConfigDir(opts.configDir, defaults)
	return dir + string(os.PathSeparator) + "server.yml"
}

// runUpdate implements the --update command family. The updater itself
// (check/yes/branch) ships with PART 23.
func runUpdate(args []string, binName string) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printUpdateHelp(os.Stdout, binName, "stable")
		return 0
	}
	switch args[0] {
	case "check", "yes", "branch":
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: update %s is not available in this build (ships with the PART 23 update subsystem)\n", args[0])
		return 1
	}
	fmt.Fprintf(os.Stderr, "ERROR: BAD_REQUEST: unknown update command %q (run --update help)\n", args[0])
	return exitUsage
}

// argAt returns args[i] or the empty string when out of range.
func argAt(args []string, i int) string {
	if i < len(args) {
		return args[i]
	}
	return ""
}
