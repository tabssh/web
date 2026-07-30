// TabSSH Web server entry point, per AI.md PART 8 (Server Binary CLI).
// Dispatch order matches the startup sequence: PHASE 1 immediate-exit
// flags, PHASE 2 --service, PHASE 3 --maintenance, PHASE 4 --update,
// PHASE 5 server startup.
package main

import (
	"fmt"
	"os"

	"github.com/tabssh/web/src/common/version"
)

// Build-time identity, injected via -ldflags "-X main.ProjectName=...
// -X main.Version=... -X main.Build=..." per AI.md PART 7.
var (
	ProjectName = "tabssh"
	Version     = ""
	Build       = ""
)

func main() {
	// Propagate build-time identity to the shared version package.
	if ProjectName != "" {
		version.ProjectName = ProjectName
	}
	if Version != "" {
		version.Version = Version
	}
	if Build != "" {
		version.Build = Build
	}

	os.Exit(run(os.Args[1:]))
}

// run parses arguments and dispatches; it returns the process exit code.
func run(args []string) int {
	binName := version.BinaryName()

	opts, err := parseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\nRun '%s --help' for usage.\n", binName, err, binName)
		return exitUsage
	}
	if err := opts.validate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", binName, err)
		return exitUsage
	}

	// PHASE 1: immediate-exit flags. Help and version never require any
	// privileges and exit immediately.
	if opts.showHelp {
		printMainHelp(os.Stdout, binName)
		return 0
	}
	if opts.showVersion {
		printVersion(os.Stdout, binName)
		return 0
	}
	if opts.shellCmd != nil {
		return runShell(opts.shellCmd, binName)
	}
	if opts.status {
		return runStatus(opts, binName)
	}

	// PHASE 2: service subcommands.
	if opts.serviceCmd != nil {
		return runService(opts.serviceCmd, opts, binName)
	}

	// PHASE 3: maintenance subcommands.
	if opts.maintenanceCmd != nil {
		return runMaintenance(opts.maintenanceCmd, opts, binName)
	}

	// PHASE 4: update subcommands.
	if opts.updateCmd != nil {
		return runUpdate(opts.updateCmd, binName)
	}

	// PHASE 5: normal server startup.
	return serve(opts)
}
