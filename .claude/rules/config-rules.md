# Configuration & Modes Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use inline YAML comments — comments go ABOVE the setting (only exception: GitHub Actions SHA-pin `# vX.Y.Z` annotations stay inline)
- Never use `strconv.ParseBool` — always `config.ParseBool()` / `config.IsTruthy()` for env vars, config values, CLI flags, API params, form inputs
- Never use `server.yaml` — always `server.yml`; auto-migrate `.yaml` → `.yml` on startup
- Never let debug mode bypass authentication or any security check — in ANY mode, including production
- Never expose `/debug/*` or pprof endpoints without `--debug`/`DEBUG=true` (they return 404 otherwise)
- Never skip path normalization/validation — all config paths, HTTP paths, file paths, and API path params go through `SafePath`/`PathSecurityMiddleware`
- Never re-select a random port after first run — port is saved to `server.yml` and persists
- Never put maintenance metadata in ad-hoc body fields — use `details` + `Retry-After`/`X-Maintenance-*` headers

## CRITICAL - ALWAYS DO

- Mode priority: `--mode` flag > `MODE` env > default `production`; Debug priority: `--debug` flag > `DEBUG` env > `--mode debug` alias > default off
- `--mode debug` = development + debug on; explicit `DEBUG=false` still wins over the alias
- Four states: production, production+debug, development, development+debug — behavior tables in PART 6
- Config source of truth: single instance = `server.yml`; cluster mode = database (server.yml becomes cache+backup, read-only fallback when DB down)
- Maintenance mode on the ONLY two critical errors (DB connection, file write) — self-heal with 30s retries, admin panel stays up with fix guidance, writes return 503
- Default port: random unused 64000-64999 on first run, saved to config; 80/443 enable Let's Encrypt challenges
- Middleware order: URLNormalize → RequestID → PathSecurity → SecurityHeaders → Allowlist → Blocklist → RateLimit → GeoIP → Auth → Logging
- Runtime env vars always checked (`NO_COLOR`, `TERM`, `DOMAIN`, `MODE`, `DEBUG`, `DATABASE_*`, `SMTP_*`); init-only vars (`CONFIG_DIR`, `DATA_DIR`, `PORT`, …) used once on first run then ignored
- Boolean parsing accepts the full truthy/falsy table (yes/no, on/off, enable/disable, si/non, …), case-insensitive; empty = default; invalid = error

## Key Rules Summary

| Topic | Rule |
|-------|------|
| Config file | `server.yml` (`{config_dir}/server.yml` per PART 4 paths) |
| DB tables | `srv_*` prefix (server) / `usr_*` prefix (users) in remote DB; server.db / users.db in SQLite |
| Config sync | DB → server.yml on every change, on startup, and every 5 min |
| Debug console line | `🔒 Running in mode: production [debugging]` |
| Server Configuration (PART 12) | Full server.yml schema, defaults, and admin-editable settings — read PART 12 before touching server config structure |

For complete details, see AI.md PART 5, 6, 12
