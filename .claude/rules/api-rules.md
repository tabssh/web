# Health, API Structure & TLS Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never invent ad-hoc error response shapes — use the canonical error body from PART 14: `{"ok": false, "error": "CODE", "message": "...", "details": {...}}`
- Never put operational metadata (retry timing, maintenance flags) in the body — use standard headers (`Retry-After`, `X-Maintenance-*`)
- Never skip `/server/healthz` maintenance-state reporting (status/mode/checks per PART 5's schema)
- Never accept boolean API params outside `config.ParseBool` semantics

## CRITICAL - ALWAYS DO

- Version the API per PART 13 (`/api/{api_version}/...`) and expose health/version endpoints per PART 13's schema
- Follow PART 14's API structure for routes, request/response formats, and status codes (standard HTTP semantics, RFC 9110)
- SSL/TLS per PART 15: Let's Encrypt HTTP-01 on port 80, TLS-ALPN-01 + auto-SSL on port 443; SSL dirs `{config_dir}/ssl/` (letsencrypt/, local/)
- Maintenance mode: writes return 503 with canonical body + `Retry-After: 30`

## Status

The NEVER/ALWAYS items above are derived from PART 0-6 content only. PART 13, 14, 15 bodies have not yet been read — read them directly and expand this file's sections before or during the first task that touches this domain.

For complete details, see AI.md PART 13, 14, 15
