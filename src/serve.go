// PHASE 5 server startup per AI.md PART 8 (Startup Sequence steps 6-21):
// context detection, one-time path resolution (flag > env > default),
// directory creation, config load-or-create with the port priority
// chain, daemonization, PID file, HTTP server with signal handling, and
// the startup banner. Root privilege-drop/system-user creation (steps
// 8a-8h) ships with the PART 24/25 service subsystem; DB, scheduler, and
// Tor initialization ship with PART 10/19/32.
package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/tabssh/web/src/common/banner"
	"github.com/tabssh/web/src/common/display"
	"github.com/tabssh/web/src/common/version"
	"github.com/tabssh/web/src/config"
	"github.com/tabssh/web/src/mode"
	"github.com/tabssh/web/src/paths"
	"github.com/tabssh/web/src/pid"
	"github.com/tabssh/web/src/runenv"
	"github.com/tabssh/web/src/server"
	signalpkg "github.com/tabssh/web/src/signal"
)

// serve runs the server startup sequence and blocks until shutdown.
func serve(opts *options) int {
	elevated := runenv.StartedElevated()

	// Resolve every directory ONCE at startup (flag > env > default),
	// locked to the starting EUID.
	defaults := paths.Resolve()
	configDir := runenv.GetConfigDir(opts.configDir, defaults)
	dataDir := runenv.GetDataDir(opts.dataDir, defaults)
	cacheDir := runenv.GetCacheDir(opts.cacheDir, defaults)
	logDir := runenv.GetLogDir(opts.logDir, defaults)
	backupDir := runenv.GetBackupDir(opts.backupDir, dataDir)
	pidPath := runenv.GetPIDFile(opts.pidFile, defaults)
	dbDir := runenv.GetDatabaseDir(dataDir)

	// Create every runtime directory, including the fixed subdirectory
	// layout for SSL, security, and Tor.
	dirs := []string{
		configDir,
		filepath.Join(configDir, "ssl"),
		filepath.Join(configDir, "tor"),
		dataDir,
		dbDir,
		filepath.Join(dataDir, "security"),
		filepath.Join(dataDir, "tor"),
		filepath.Join(dataDir, "tor", "site"),
		cacheDir,
		logDir,
		backupDir,
	}
	for _, d := range dirs {
		if err := runenv.EnsureDir(d, elevated); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
			return 1
		}
	}

	// Load or create server.yml (first run picks a random 64xxx port and
	// persists it).
	configPath := filepath.Join(configDir, paths.ConfigFileName)
	cfg, firstRun, err := config.LoadOrCreate(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: loading config: %v\n", err)
		return 2
	}

	// Apply mode with the flag > env > config > default priority: FromEnv
	// applies MODE/DEBUG, then flags override, then config fills the gap.
	mode.FromEnv()
	if opts.mode != "" {
		mode.SetAppMode(opts.mode)
	} else if os.Getenv("MODE") == "" && cfg.Server.Mode != "" {
		mode.SetAppMode(cfg.Server.Mode)
	}
	if opts.debug {
		mode.SetDebugEnabled(true)
	}

	// Port priority: --port > {PROJECT_NAME}_PORT env > config value.
	port := cfg.Server.Port
	if v := os.Getenv("TABSSH_PORT"); v != "" {
		port = v
	}
	if opts.port != "" {
		port = opts.port
	}
	address := cfg.Server.Address
	if opts.address != "" {
		address = opts.address
	}
	if address == "" {
		address = "0.0.0.0"
	}

	// Daemonize when requested by flag/config and appropriate for the
	// detected service manager context.
	if runenv.ShouldDaemonize(false, opts.daemon, bool(cfg.Server.Daemonize)) {
		if err := runenv.Daemonize(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
			return 1
		}
	}

	// PID file: stale detection + write (container-aware no-op inside).
	if bool(cfg.Server.PIDFile) && !runenv.IsContainer() {
		if err := runenv.EnsurePIDFileDir(pidPath, elevated); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
			return 1
		}
		if err := pid.WritePIDFile(pidPath); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: CONFLICT: %v\n", err)
			return 1
		}
	} else {
		pidPath = ""
	}

	// HTTP server with the PART 8 request-ID middleware applied; routes
	// register on srv.Mux() when PART 13/14 lands.
	srv := server.New(net.JoinHostPort(address, port))

	// Signal handling: graceful shutdown, SIGUSR1/SIGUSR2, SIGHUP ignore.
	signalpkg.Setup(srv.HTTPServer(), pidPath)

	printBanner(opts, cfg, port, firstRun)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(os.Stderr, "ERROR: SERVER_ERROR: %v\n", err)
		if pidPath != "" {
			pid.RemovePIDFile(pidPath)
		}
		return 1
	}

	// Graceful shutdown path: the signal handler completes cleanup and
	// exits the process; block until it does.
	select {}
}

// printBanner renders the startup banner with the resolved URLs.
func printBanner(opts *options, cfg *config.Config, port string, firstRun bool) {
	env := display.DetectDisplayEnv()
	unicode := !env.IsDumbTerminal() && display.ColorEnabled(opts.color, "", env)

	host := cfg.Server.FQDN
	if host == "" {
		host = "localhost"
	}
	url := "http://" + host
	if port != "" && port != "80" && port != "443" {
		url += ":" + port
	}
	if opts.baseURL != "" && opts.baseURL != "/" {
		url += "/" + strings.Trim(opts.baseURL, "/")
	}

	banner.PrintStartupBanner(banner.BannerConfig{
		AppName: "TabSSH Web",
		Version: version.Version,
		AppMode: mode.GetAppModeString(),
		Debug:   mode.IsDebugEnabled(),
		URLs:    []string{url},
		// The one-time setup token lives in server.db (PART 17); until
		// that subsystem ships there is no token to display.
		ShowSetup: false,
		Unicode:   unicode,
	})
	if firstRun {
		fmt.Println("First run: default configuration created at " + filepath.Join(runenv.GetConfigDir(opts.configDir, paths.Resolve()), paths.ConfigFileName))
	}
}
