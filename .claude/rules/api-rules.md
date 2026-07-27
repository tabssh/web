# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- NEVER expose in `/server/healthz` or any health endpoint: tokens/API keys/passwords, DB connection strings/host/port/user, internal IPs, config/data file paths, admin usernames, SMTP host/admin emails, encryption/session secrets, stack traces/debug info, internal endpoint status (e.g. `/metrics`).
- NEVER add sub-routes to health endpoints — just `/server/healthz`, not `/server/healthz/db` or `/server/healthz/**`.
- NEVER reuse a version number once released; never use a leading zero in a SemVer component; never put a `v` prefix in the version string itself (Git tags DO get the `v` prefix, e.g. `v1.0.0`).
- NEVER keep legacy API endpoints — delete them; no backwards-compatibility shims, redirects, or deprecation periods.
- NEVER hardcode `v1` (or any API version) in route code — always use `APIBasePath()` or `{api_version}`.
- NEVER use singular nouns, uppercase, underscores, trailing slashes, or verbs in route paths (routes must be plural, lowercase, hyphen-separated, no trailing slash, no verbs).
- NEVER implement SPA/client-side routing, client-side rendering (React/Vue-style), client-side data fetching for initial page load, business logic in JS, or make JavaScript required for a core feature.
- NEVER manually edit generated Swagger/OpenAPI or GraphQL schema files — they are auto-generated at build time and must stay in sync with each other and with the API.
- NEVER put Swagger/GraphQL files anywhere except `src/swagger/` and `src/graphql/` — files in the project root are forbidden.
- NEVER serve `.yaml` OpenAPI or add a `.json` suffix to the OpenAPI path — OpenAPI is JSON only.
- NEVER implement redirects from old removed paths (`/openapi`, `/openapi.json`, `/graphql` GET/POST at root) — they are no longer served at all.
- NEVER use forbidden/legacy error response shapes — only the canonical envelope is allowed.
- NEVER skip RFC compliance when the application itself implements an RFC-defined protocol — this is non-negotiable.
- NEVER set `DOMAIN` to an overlay-network address (`.onion`, `.i2p`, `.exit`).
- NEVER use icons, ASCII art, or special characters in log output — logs are always raw text.

## CRITICAL - ALWAYS DO
- ALWAYS end every non-HTML response (and every file) with a single trailing newline.
- ALWAYS order `HealthResponse` fields canonically: Project; Status/PendingRestart/RestartReason; Version/GoVersion/Build; Uptime/Mode/Timestamp; Cluster; Features (public only); Checks; Stats; then app-specific fields.
- ALWAYS start first stable release at SemVer `1.0.0` (never `0.x.x`); use `1.0.0-rc1`-style pre-release suffixes when needed.
- ALWAYS follow version format by channel: Stable = SemVer, Beta = `YYYYMMDDHHMMSS-beta`, Daily = `YYYYMMDDHHMMSS`.
- ALWAYS resolve the version from, in priority order: 1) `release.txt` in project root, 2) Git tag, 3) fallback `dev`.
- ALWAYS require API versioning on every route (`/api/{api_version}/...`).
- ALWAYS implement all three API types for every project: REST (primary), Swagger/OpenAPI, and GraphQL.
- ALWAYS follow the canonical response envelope: success `{"ok": true, "data": {...}}`; error `{"ok": false, "error": "CODE", "message": "...", "details": {...}}`.
- ALWAYS use 2-space indentation for HTML/JSON/YAML/JS/CSS, and tabs for Go code/Makefiles; never leave trailing whitespace.
- ALWAYS honor content negotiation: `.txt` extension always returns plain text; Accept header determines JSON vs text vs HTML; non-interactive/CLI/HTTP-tool clients get plain text by default.
- ALWAYS use client-type detection functions where relevant: `isOurCliClient(r)`, `isTextBrowser(r)`, `isHttpTool(r)`, `isNonInteractiveClient(r)`.
- ALWAYS provide built-in Let's Encrypt support in every project, supporting HTTP-01 (default, port 80), TLS-ALPN-01 (port 443), and DNS-01 (wildcard) challenge types.
- ALWAYS resolve FQDN in this priority order: 1) reverse-proxy headers, 2) `DOMAIN` env var, 3) `os.Hostname()`, 4) `$HOSTNAME` env var, 5) public IPv6, 6) public IPv4, 7) `localhost` fallback.
- ALWAYS strip `:80` and `:443` from any displayed URL.
- ALWAYS look up certificates in this priority order: `/etc/letsencrypt/live/domain/`, `/etc/letsencrypt/live/{fqdn}/`, `{config_dir}/ssl/letsencrypt/{fqdn}/` (app-managed), `{config_dir}/ssl/local/{fqdn}/` (user-provided, no auto-renewal).
- ALWAYS check for app-managed Let's Encrypt certificate renewal daily at 03:00.

## Key rules summary
- Health endpoint family: `/server/healthz`, optional root alias `/healthz` (gated by `server.healthz.root.enabled`), `/api/{api_version}/server/healthz`, and unversioned alias `/api/healthz`. Rule of thumb for health data: "If in doubt, leave it out" — only intentionally public-safe fields belong in health responses.
- Auth token extraction: accepted from any supported header (see PART 8) or `?token=` query param, checked in a defined priority order, first match wins.
- Route Migration Rule: migrate the implementation to the new route first, then delete the superseded route — no parallel-run period.
- Server-Side Processing Philosophy: "The server does the work. The client displays the result." Client JS is allowed only for form validation feedback, theme toggle, copy-to-clipboard, auto-refresh/polling, modal/dropdown UX, keyboard shortcuts, and progressive enhancement.
- External API Compatibility default is feature/behavior compatibility using the project's own routes, not route-for-route cloning of a third-party API.
- HTTP status codes follow RFC 7231: 200/201/204/301/302/400/401/403/404/405/409/422/429/500/503.
- Text response format for CLI/agent-facing output: `OK: {message}` / `ERROR: {code}: {message}`.
- JSON rules: no comments, no trailing commas, double quotes only, no `undefined`, 2-space indent, trailing newline.
- Port configuration (non-negotiable): single port defaults to HTTP, except a single port of 443 which is HTTPS-only; dual ports means first is HTTP, second is HTTPS.
- Startup banner is responsive to terminal width; logs remain plain text always.

Reference: AI.md PART 13, 14, 15 — Health & Versioning; API Structure; SSL/TLS & Let's Encrypt
