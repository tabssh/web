# Backend Rules (PART 9, 10, 11, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Never use `DROP COLUMN`, `DROP TABLE`, or `DELETE` for schema updates — schema changes are additive only, never destructive
- Never rename a database column directly — add the new column with a default, have the app read-new/write-both, then deprecate (never drop) the old one
- Never return stack traces, full Go error chains, SQL queries/EXPLAIN output, or internal IP/hostnames to production responses — those are Tier 3 (debug-only) or Tier 1 (never public) data
- Never expose DB credentials, internal IPs/hostnames, API/session/CSRF tokens, recovery keys, MFA secrets, other users' PII, filesystem paths, or account-existence signals on any publicly reachable endpoint, even in debug mode
- Never cast user-controlled content to `template.HTML`; never render user-supplied HTML/SVG/XML inline by default — force attachment download instead
- Never execute user-supplied content server-side; never shell out with raw user content/filenames/refs; never run hooks/build steps/interpreters on untrusted content unless explicitly sandboxed
- Never log or persist plaintext passwords/hashes, full API token or session token values, recovery keys, TOTP secrets, private keys, credit card numbers, or full email addresses in audit logs
- Never use bcrypt/MD5/SHA-* for password hashing — Argon2id only
- Never emit `Expect-CT`, `Public-Key-Pins`/HPKP, or `Feature-Policy` headers (deprecated/superseded)
- Never use the default Tor ports 9050/9051 or the system Tor daemon — the app starts its own dedicated Tor process, always via a runtime-detected `127.0.0.1:auto` control port, never a hardcoded/saved port
- Never let the server fail to start because of a Tor error — Tor is optional; missing Tor is logged INFO, not an error
- Never overwrite an existing `torrc` on normal startup — only regenerate it on an explicit admin-driven config change
- Never treat agents as cluster nodes — agents never share the database and never receive `app_secrets`
- Never return API tokens, secrets, or encryption keys in API responses (even to admins) — only "configured" / "rotated N days ago" / a fingerprint hash
- Never allow the Allowlist (trusted IPs) to bypass admin authentication, API token validation, CSRF protection, path security, or SSL/TLS — it only bypasses IP blocklists, rate limiting, GeoIP blocking, and account lockout
- Never disable breach detection — thresholds are configurable, detection itself cannot be turned off

## CRITICAL - ALWAYS DO
- Always use `CREATE TABLE IF NOT EXISTS` and `ADD COLUMN IF NOT EXISTS` (or the ignore-duplicate-error equivalent per DB engine) for self-creating, idempotent schema — no separate migrations table
- Always apply a context timeout to every query (Simple SELECT 5s, Complex/JOIN SELECT 15s, INSERT/UPDATE/DELETE 10s, bulk ops 60s, migrations 5min, reports 2min) and use connection pooling sized per environment
- Always retry serializable-isolation transaction failures (Postgres `40001`, MySQL `1213`) with backoff; never retry non-retryable 4xx-class errors
- Always log every error with structured context (timestamp, error, request context) — never log and swallow silently
- Always return the canonical API response envelope (`{"ok": true/false, ...}`) and map errors to the standard error-code table with correct HTTP status
- Always support cluster mode with config sync, session sharing, distributed locks, primary election, and heartbeat/health monitoring when an external cache/shared DB is detected
- Always use an advisory lock plus majority quorum check before rotating cluster-wide secrets; abort the rotation (`cluster.rotation_aborted_no_quorum`) if quorum isn't reached
- Always wipe local `app_secrets`, cluster-shared PGP private key material, and the cluster section of `server.yml` on receiving `403 NODE_REMOVED`, and refuse to serve until re-joined
- Always default to Tier 3 (debug-only) when unsure whether a field is Tier 2 or Tier 3, and to Tier 1 (never public) when unsure between Tier 1 and Tier 3
- Always run the Output Sanitization Pipeline (allow-list fields, redact sensitive query params, strip internal IPs/paths, truncate, strip debug-only fields, constant-time finalize) on every response
- Always hash passwords with Argon2id, API tokens with SHA-256, and encrypt 2FA secrets with AES-256-GCM using `server.security.encryption_key`
- Always send the standard security headers on every response (`X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, CSP, Permissions-Policy, etc.), and add HSTS whenever TLS is enabled
- Always generate `installation_secret`, `cookie_signing_key`, `csrf_token_secret`, and `server.security.encryption_key` on first start before any user-visible operation, and never log or return them in any response
- Always write audit log entries in JSON Lines format with `id` (ULID), `time` (ISO8601+ms UTC), `event`, `category`, `severity`, `actor`, and `result`; keep audit logs append-only with file permissions `0640`
- Always disable all compliance standards (GDPR/CCPA/HIPAA/SOC2/PCI-DSS/ISO27001/FedRAMP/LGPD/PIPEDA/APPI/PDPA) by default, enabling per-project, and apply the strictest overlapping requirement when multiple are enabled
- Always start Tor as its own dedicated process running as the same user the server runs as (after privilege drop), with Tor directories created at `0700` and Tor files (torrc, keys) at `0600`, ownership enforced to the current user on every startup
- Always route hidden-service traffic via `control.AddOnion()` port mapping (.onion:80 → `127.0.0.1:{server_port}`) rather than static `HiddenServiceDir` torrc config, using ED25519 v3 keys

## Key rules summary
- Cache key naming: lowercase, colon-separated, version-prefixed for busting; default TTLs range from 1min (rate-limit counters) to 24h (static content), with API tokens having no expiry
- HTTP cache headers differ by resource class: static assets `public, max-age=31536000, immutable`; HTML pages and private/authenticated responses `no-store`/`private, no-store`; public API `public, max-age=60`
- Standard error codes map to fixed HTTP statuses (e.g. `VALIDATION_FAILED`→400, `TOKEN_EXPIRED`→401, `ACCOUNT_LOCKED`→403, `RATE_LIMITED`→429, `SERVER_ERROR`→500)
- Cluster heartbeat interval is 30s with a 90s timeout (3 missed beats = unresponsive); node states are healthy/degraded/offline/removed, and secret-version drift over 7 days marks a node `stale`
- Primary election picks the lowest node ID; there is no preemption when the old primary returns, and split-brain is resolved by "latest write wins" in the DB
- API token format is `{prefix}_{32-alphanumeric}` with prefixes `adm_`/`usr_`/`org_` (plus agent compounds `adm_agt_`/`usr_agt_`/`org_agt_`); only the SHA-256 hash and an 8-char prefix are stored
- Password policy auto-upgrades under compliance (min length 8→12, uppercase/number/special required, max age 90 days, history 12) — these are minimums; admins may go stricter but never weaker
- Breach detection thresholds trigger tiered automated responses (e.g. 10+ failed logins/5min → block IP; mass data export → queue+alert; privilege escalation → block session, Critical severity)
- Account lockout is separate from IP blocking: 5 fails/15min soft-locks 15min, 10/1hr hard-locks 1hr, 15/24hr locks until admin unlock or password reset
- Well-known files (`robots.txt`, `security.txt` per RFC 9116, `llms.txt`) are served from a reserved, root-owned `.well-known/` namespace never claimable by users/orgs
- Tor config is validated both client- and server-side (bandwidth pattern, circuit/timeout ranges, intro points 3-10); changes require a Tor restart and briefly interrupt the hidden service
- Tor storage paths are fixed and never configurable: `{config_dir}/tor/torrc`, `{data_dir}/tor/site/` (hidden-service keys), `{log_dir}/tor.log`

Reference: AI.md PART 9, 10, 11, 32 — Error Handling & Caching, Database & Cluster, Security & Logging, Tor Hidden Service
