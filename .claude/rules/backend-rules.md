# Error Handling, Database, Security/Logging & Tor Rules (PART 9, 10, 11, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use `mattn/go-sqlite3` (CGO) — SQLite driver is `modernc.org/sqlite`; never `lib/pq` (use pgx); never `dgrijalva/jwt-go`; never old `go-redis/redis` import path
- Never hash passwords with bcrypt/scrypt/MD5/SHA — Argon2id only (time=3, memory=64*1024 KB, threads=4, keyLen=32, saltLen=16)
- Never store or log raw API/session tokens — SHA-256 hash them
- Never build SQL by string concatenation — parameterized queries only
- Never show stack traces or internal details to public users in production — generic errors outward, detailed errors in logs

## CRITICAL - ALWAYS DO

- DB schema follows PART 5's table summary: `srv_*` server tables / `usr_*` user tables (remote DB) or server.db / users.db (SQLite)
- Cluster mode: remote DB is source of truth, server.yml is cache+backup, read-only fallback on DB outage
- Self-healing per PART 5 maintenance-mode rules for DB/disk critical errors
- Security middleware order per PART 5 (PathSecurity early, Auth after allowlist/blocklist/ratelimit/geoip)
- Structured logging via `log/slog` with request IDs; debug-level DB/cache logging only when debug enabled
- Tor hidden service support per PART 32 (onion address is per-node cluster state)

## Status

The NEVER/ALWAYS items above are derived from PART 0-6 content only. PART 9, 10, 11, 32 bodies have not yet been read — read them directly and expand this file's sections before or during the first task that touches this domain.

For complete details, see AI.md PART 9, 10, 11, 32
