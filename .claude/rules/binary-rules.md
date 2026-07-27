# Binary Rules (PART 7, 8, 33)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Never build on the host or outside `casjaysdev/go:latest` (or project toolchain image) — CI/CD runs `go build` directly inside the image, not via `actions/setup-go`.
- Never use `strconv.ParseBool()` for CLI/Agent boolean flags or config — always `config.ParseBool()` / `config.IsTruthy()`.
- Never construct URLs by `fmt.Sprintf` with raw user input (usernames, org names, IDs, search terms, filenames) — always `urlutil.BuildAPIURL()` / `EncodePathSegment()` / `EncodeQueryValue()`.
- Never let the CLI's User-Agent header reflect a renamed binary — it must always be the compiled `ProjectName`, e.g. `{project_name}-cli/{version}`, even if the executable is renamed (`mypaste`, `pb`, etc.).
- Never implement `--tui`/`--cli`/`--gui`/`--mode tui|cli|gui` flags on the CLI — display mode is always auto-detected, never explicitly selectable.
- Never require root/admin privileges for `--help` or `--version` at any command level (main, subcommand, nested) — no escalation checks, no `sudo`, exit immediately with help text.
- Never use short flag forms for anything except `-h`/`--help` and `-v`/`--version` — every other flag is long-form only.
- Never give the CLI a `config` command, a `tui` command, or `--port`/`--address` flags (those don't apply to a client).
- Never give the Agent `--port` or `--address` flags — the agent never listens for inbound HTTP connections.
- Never give the Agent admin flags (`--admin*`), user operations, or interactive modes (TUI/GUI) — agent does only its designated job, headless.
- Never run the Agent inside a Docker container or Kubernetes pod — agents must run directly on the target system (systemd/Windows service/launchd) for full system access and accurate metrics.
- Never use a setup token for Agent onboarding — Agent registration is always initiated from the server admin/user/org panel via a generated one-liner (`--server URL --token TOKEN`), never the server's own `/setup` setup-token flow.
- Never use OS system directories for CLI runtime state — no `/etc/...`, `/var/lib/...`, `/var/log/...`, `C:\ProgramData\` for the CLI, even when the CLI binary is installed system-wide or invoked by root; CLI state is always user-scope.
- Never resolve or cache `$HOME`/user directories for the Agent after a privilege drop — Agent runs in system context (`/etc`, `/var/lib`, `/var/log`, `/var/cache`), not user context.
- Never clear or overwrite an existing valid `server.primary`/token config value with an empty or unvalidated one — only save new values when current is empty and the new value validates.
- Never ship separate shell-completion files — completions (`--shell completions`/`--shell init`) must be generated live from the running binary, not distributed as static files.
- Never fabricate or bypass the CLI's Auth priority order (`--token` flag > `{PROJECT_NAME}_TOKEN` env var > `auth.token` in cli.yml) — env var token must never be persisted to the config file.
- Never expect users to type the project's real name for CLI operation — internal identifiers (config path, User-Agent) must always work regardless of the binary's actual filename.
- Never hardcode host, IP, or port anywhere in project code — `localhost`, `127.0.0.1`, `0.0.0.0`, `[::1]` are never hardcoded, always detected dynamically.
- Never skip display-environment / `TERM=dumb` / `NO_COLOR` detection in any binary (server, CLI, agent).
- Never leave a stale PID file undetected after a crash/`kill -9`.

## CRITICAL - ALWAYS DO
- Always ship a CLI client for every project — the CLI is REQUIRED, non-optional, no include/skip decision (only the Agent is a per-project determination).
- Always compile the true project name into the binary via `-ldflags "-X main.ProjectName=... -X main.Version=..."` and use it for User-Agent and config directory naming; use `filepath.Base(os.Args[0])` only for display purposes (`--help`, `--version`, error messages).
- Always support both `--flag=value` and `--flag value` syntax on every flag, on every binary.
- Always provide a config-file equivalent for every flag (`defaults.*` in cli.yml / matching keys in agent.yml) — the config file setting is the source of truth for the flag's default.
- Always implement automatic cluster failover for both CLI and Agent: try `server.primary`, silently fall back to `server.cluster` nodes, periodically call `/api/autodiscover` to refresh `server.primary`/`server.cluster`, and update the config file asynchronously/non-blocking.
- Always validate before saving any server/token config value; only save valid values; never clear a valid existing value with an empty/invalid one (`saveIfEmpty` pattern).
- Always build shell completions into the binary itself for ALL binaries (server, agent, client) via `--shell completions [SHELL]` and `--shell init [SHELL]`, auto-detecting from `$SHELL` when omitted; support bash, zsh, fish, sh, dash, ksh, powershell, pwsh.
- Always make CLI runtime state fully user-scoped (XDG dirs on Linux/macOS, `%APPDATA%`/`%LOCALAPPDATA%` on Windows) regardless of whether the binary is installed system-wide or invoked as root.
- Always run `EnsureDirs()`/`EnsureFile()` with correct permissions on CLI startup — Unix: `chmod 0700` dirs / `0600` files; Windows: rely on ACL inheritance (no-op).
- Always detect and respond to terminal size changes (SIGWINCH on Unix, polling on Windows) in the CLI TUI, and reflow layout using the `SizeMode` breakpoints (Micro/Minimal/Compact/Standard/Wide/Ultrawide/Massive) from `common/terminal`.
- Always support small/phone-SSH terminal sizes down to Micro (<40 cols/<10 rows) with graceful degradation — this is a required use case, not optional.
- Always URL-encode all user input placed into URLs using `EncodePathSegment()` (path segments) or `EncodeQueryValue()` (query values) — never raw string concatenation/`Sprintf`.
- Always make the Agent's directory structure and privilege model match the server's: same `{config_dir}/agent.yml` alongside `server.yml`, same `{data_dir}`, `{data_dir}/db/agent.db`, root/admin privileges required for full system access.
- Always generate the Agent's one-line connect command from the server (`{project_name}-agent --server {url} --token {token}`) and have the agent auto-register, auto-configure, and self-install as a service (if root) on first connect.
- Always scope Agent tokens by owner type with the correct prefix: `adm_agt_` (admin-issued), `usr_agt_` (user-issued), `org_agt_` (org-issued).
- Always exit CLI commands with the standardized exit codes: `0` success, `1` general error, `2` config error, `3` connection error, `4` authentication error, `5` not found, `64` usage error (bad arguments).
- Always show the actual (possibly renamed) binary name in `--help`/`--version`/error text while using the compiled/original project name for the User-Agent header and config file path — these must never be conflated.
- Always make `make build` build the CLI and Agent binaries automatically whenever `src/client/` / `src/agent/` exist, alongside the server, from the same source tree and `src/common/` shared packages.
- Always detect display environment (GUI/TUI/CLI/Headless) and adapt output in every binary; `TERM=dumb` disables all ANSI escapes and forces CLI mode; respect `NO_COLOR` (https://no-color.org/) everywhere.
- Always create missing directories automatically for every directory flag (`--config`, `--data`, `--cache`, `--log`, `--backup`, `--pid`).
- Always detect stale PID files (crash/`kill -9` recovery) with container-aware skip logic.
- Always resolve `{proto}`, `{fqdn}`, `{port}` dynamically from request context (reverse-proxy headers preferred: `X-Forwarded-Proto`/`X-Forwarded-Ssl`/`X-Forwarded-Port` before fallbacks) and use them consistently across templates, API docs, and email links.
- Always build CGO_ENABLED=0 for all 8 release platform targets (linux/darwin/windows/freebsd × amd64/arm64).

## Key rules summary
- Binary naming: server `{project_name}`, CLI `{project_name}-cli` (REQUIRED), agent `{project_name}-agent-{os}-{arch}` (OPTIONAL, e.g. `monitor-agent-linux-amd64`).
- `--config` flag resolution: bare name → `{config_dir}/{name}.yml`; explicit `.yml`/`.yaml` used as-is; absolute/`~` paths used verbatim; no extension auto-detects `.yml` then `.yaml`, defaulting to `.yml` for new configs.
- Config precedence (both CLI and Agent): CLI flag > environment variable (`{PROJECT_NAME}_{SECTION}_{KEY}` / `{PROJECT_NAME}_AGENT_{KEY}`) > config file > compiled default.
- Terminal size breakpoints (`src/common/terminal/size.go`): Micro <40 cols/<10 rows, Minimal 40-59/10-15, Compact 60-79/16-23, Standard 80-119/24-39, Wide 120-199/40-59, Ultrawide 200-399/60-79, Massive ≥400/≥80 — each has a defined `MaxTableColumns()` and help-text density.
- Required TUI libraries: `charmbracelet/bubbletea`, `charmbracelet/bubbles`, `charmbracelet/lipgloss`, `golang.org/x/term`.
- Agent communication patterns: Send Only (agent → server, e.g. metrics/log shippers — lower risk), Receive Only (server → agent, e.g. config pull — medium risk), Bidirectional (full command-and-control, e.g. CI/CD runners, remote management — highest risk, requires mTLS/tokens, authorization, audit logging, input validation, sandboxing).
- Agent vs Client execution context: Client runs user-scope (`~/`) with no privilege escalation; Agent runs system-scope (`/etc`, `/var/lib`, `/var/log`, `/var/cache` on Linux) requiring root/admin.
- Client = full remote administration surface (`--admin*` flags, TUI/CLI/GUI modes); Agent = purpose-specific worker only (no admin, no user ops, no interactive modes).
- Output formats: `json`, `table` (box-drawing), `plain` — selectable via `--output`, defaulting per `output.format` in cli.yml.
- Universal flags on every binary: `--help`/`-h`, `--version`/`-v`, `--color {auto|yes|no}`, `--lang CODE` — only `-h`/`-v` get short forms.
- Boolean parsing is shared between server/CLI/agent: truthy = `true, yes, on, 1, enable, enabled`; falsey = `false, no, off, 0, disable, disabled, none`.
- Server startup sequence differs from Agent only at steps 6-8: server connects to its own database and listens for HTTP; agent connects to the parent server and starts its reporter/collector.
- Server binary CLI must expose standard flags (`--mode`, `--debug`, `--service --install/start/stop/restart/status`, directory flags, signal handling) per PART 8; never display a route without its full URL (`{proto}://{fqdn}:{port}/path`), and strip default ports `:80`/`:443` from displayed/generated URLs.

Reference: AI.md PART 7, 8, 33 — Binary Requirements, Server Binary CLI, Client & Agent
