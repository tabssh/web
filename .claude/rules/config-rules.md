# Configuration & Application Modes Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Inline YAML comments — ALL comments go ABOVE the setting, never on the same line (exception: GitHub Actions SHA-pin version annotations)
- Use `strconv.ParseBool()` directly — always use `config.ParseBool()` / `config.IsTruthy()`
- Use `server.yaml` — config file is always `server.yml` (auto-migrate `server.yaml` → `server.yml` on startup if found)
- Ship a 1000-line config — sane built-in defaults, everything configurable but commented/defaulted
- Disable authentication or security checks in debug/development mode — debug affects verbosity/diagnostics ONLY, never auth
- Bind to a privileged port (<1024) without either an escalated service install or fallback to >1024
- Hardcode a fixed default port — first run picks a random unused port in 64000-64999 and persists it to `server.yml`
- Skip path normalization/validation on configuration values, HTTP request paths, file paths, or API path parameters

## CRITICAL - ALWAYS DO
- Mode priority: `--mode` CLI flag > `MODE` env var > default `production`
- Debug priority: `--debug` CLI flag > `DEBUG` env var (truthy) > `--mode debug`/`MODE=debug` alias > default `false`
- `--mode debug` / `MODE=debug` expands to development mode + debug on; an explicit `--debug`/`DEBUG` still wins (`MODE=debug DEBUG=false` = development, debug off)
- Support all four operational states: Production, Production+Debug, Development, Development+Debug
- Production: info-level logging, debug/pprof endpoints disabled (404), generic error messages, template/static caching enabled, rate limiting enforced, all security headers on
- Development: debug-level logging, hot reload (caching disabled), relaxed rate limiting/CORS — but `/debug/*` and pprof still require the explicit `--debug` flag
- Debug flag (any mode): enables `/debug/*`, `/debug/pprof/*`, `/debug/vars`, full request/response logging, DB and cache query logging — admin authentication is NEVER bypassed
- Boolean parsing must accept the full truthy/falsy word set (yes/no, on/off, enable/disable, y/n, si/non, oui/niet, etc.), case-insensitive; empty string uses the default; invalid value is an error
- Normalize and validate every path (config values like `admin_path`/`static_path`, HTTP request paths, file paths, API path parameters) using `SafePath`/`normalizePath`/`validatePath` — reject `..`, invalid characters, oversized segments
- `PathSecurityMiddleware` runs third in the middleware chain (after URL normalization and request-ID attachment, before security headers/allowlist/blocklist/rate-limit/geoip/auth/logging)
- Single-instance mode (SQLite): `server.yml` is the source of truth, admin panel writes to it directly
- Cluster mode (remote DB): database is the source of truth, `server.yml` is cache+backup; if DB unavailable, fall back to read-only mode using cached config, synced every 5 minutes and on every config change
- Only two errors are truly critical — database connection failure and file write failure (disk full/permissions); everything else triggers self-healing retries, not maintenance mode
- Maintenance mode: reject writes with HTTP 503 + canonical error body, `Retry-After` and `X-Maintenance-*` as headers (not body fields); admin panel stays accessible with fix guidance; exit automatically once self-healing succeeds

## Key rules summary
- Config file location: root/privileged → `/etc/tabssh/tabssh/server.yml` (Linux); regular user → `~/.config/tabssh/tabssh/server.yml`; Docker → `/config/tabssh/server.yml`
- Port rules: default is a random unused port in 64000-64999 on first run, then persisted permanently; port 80 auto-enables Let's Encrypt HTTP-01; port 443 auto-enables TLS-ALPN-01 + SSL; port 0 lets the OS assign a port; dual-port format `"8090,8443"` (HTTP,HTTPS) is supported
- Privileged port binding (<1024): service install (`sudo {project_name} --service --install`) escalates once, binds the port as root, then drops to a dedicated `{internal_name}` system user; user-mode runs (no sudo) are restricted to ports >1024
- Sensitive CLI operations (`--maintenance setup`, `--maintenance restore`, `--maintenance mode`) require real authorization (first-run, root/admin, valid token, or admin credentials) — not just filesystem access
- Environment variables: runtime (always checked) — `NO_COLOR`, `TERM`, `DOMAIN`, `MODE`, `DATABASE_DRIVER`, `DATABASE_URL`, `SMTP_HOST`/`SMTP_PORT`/`SMTP_USERNAME`/`SMTP_PASSWORD`/`SMTP_FROM_NAME`/`SMTP_FROM_EMAIL`/`SMTP_TLS`; init-only (first run only, then ignored) — `CONFIG_DIR`, `DATA_DIR`, `LOG_DIR`, `DATABASE_DIR`, `BACKUP_DIR`, `PORT`, `LISTEN`, `APPLICATION_NAME`, `APPLICATION_TAGLINE`
- URL variable resolution prefers reverse-proxy headers first: `{fqdn}` via `X-Forwarded-*` → `DOMAIN` → `os.Hostname()` → `$HOSTNAME` → global IP → `localhost`; `{proto}` via `X-Forwarded-Proto`/`X-Forwarded-Ssl`/`X-Url-Scheme` → TLS detection → `http`
- Mode shortcuts: `--mode dev`/`devel`/`development` → development; `--mode prod`/`production` → production; `--mode debug` → development + debug on
- Full database schema for config storage includes `config`/`srv_config`, `config_meta`, `admin_sessions`, `rate_limits`, `audit_log`, `scheduler_tasks`, `scheduler_history`, `backups` (server tables) and `admins`, `users`, `api_keys`, `password_resets`, `email_verifications`, `totp_secrets`, `passkeys`, `trusted_devices`, `user_sessions`, `custom_domains`, `custom_domain_audit` (user tables)

Reference: AI.md PART 5, PART 6, PART 12 — Configuration, Application Modes, Server Configuration
