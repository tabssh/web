# Project Structure & License Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never add a GPL/AGPL/LGPL-licensed dependency (project license is MIT — incompatible)
- Never leave `LICENSE.md`'s embedded-licenses section stale — every real `go.mod` dependency must be listed with version/license/copyright
- Never use `mattn/go-sqlite3`, `lib/pq`, `ooni/go-libtor`, `dgrijalva/jwt-go`, `gorilla/mux`, or the old `go-redis/redis` import path (forbidden libraries)
- Never enable CGO — always `CGO_ENABLED=0`; use `modernc.org/sqlite`, never a CGO SQLite driver
- Never hash passwords with anything but Argon2id (OWASP 2023 params: time=3, memory=64*1024 KB, threads=4, keyLen=32, saltLen=16); API/session tokens use SHA-256
- Never hardcode a path — always resolve via the OS-specific path table (PART 4) through `src/paths` / `src/runenv`
- Never use `server.yaml` — config file is always `server.yml` (auto-migrate `.yaml` → `.yml` on startup if found)
- Never add a root file that isn't in the Allowed Root Files table
- Never gitignore `.claude/rules/`, `.claude/agents/`, `.claude/hooks/`, `.claude/commands/`, `.claude/plans/`, or `.claude/CLAUDE.md` — those are committed; only `.claude/settings.local.json`, `.claude/*.lock`, `.claude/backups/`, `.claude/cache/`, `.claude/history.jsonl` are gitignored

## CRITICAL - ALWAYS DO

- MIT license; `LICENSE.md` structure = MIT text + `## Embedded Licenses` section listing every third-party dependency (library/version/license/copyright)
- README.md must include the license badge and a Disclaimer section per PART 1's template
- `.gitignore` line 1 = `# gitignore created on MM/DD/YY at HH:MM` (generated once, never updated), line 2 = literal `ignoredirmessage`
- Follow the canonical directory tree from PART 3 (`docs/`, `scripts/`, `tests/`, `docker/`, `volumes/` (gitignored), `binaries/`/`releases/` (gitignored), `.claude/rules/*.md`)
- Resolve OS-specific paths per PART 4's tables (Linux/macOS/BSD/Windows/Docker, privileged vs. user)
- Config file location: root → `/etc/{internal_org}/{internal_name}/server.yml`, user → `~/.config/{internal_org}/{internal_name}/server.yml`; Docker → `/config/{project_name}/server.yml`
- Use required Go libraries per PART 3's tables (auth: TOTP/WebAuthn/JWT/OIDC/OAuth2/LDAP/SAML/sessions are mandatory even without end-users)

## Verified Compliant (this project)

- `LICENSE.md`: MIT text intact, embedded-licenses table lists all 5 go.mod deps (uuid, x/net, x/sys, x/term, yaml.v3), all BSD-3-Clause/MIT/Apache-2.0 — compliant
- `README.md`: license badge present (line 3), Disclaimer section present (line ~170), License section present
- `.gitignore`: correct header format, correct nuanced `.claude/`/`.cursor/`/`.aider/`/`.ai/`/`.windsurf/` ignore rules (does not blanket-ignore committed `.claude/` subdirs)
- `src/paths/paths.go`, `src/runenv/*.go`: implement OS-specific path resolution matching PART 4's Linux/macOS/BSD/Windows/Docker tables
- `src/config/bool.go`: `ParseBool`/`MustParseBool`/`IsTruthy`/`IsFalsy` match PART 5's boolean-handling spec exactly
- `go.mod`: no forbidden libraries; `CGO_ENABLED=0` compatible (modernc.org/sqlite not yet added — no DB driver dependency present yet, not a violation at this stage)

For complete details, see AI.md PART 2, 3, 4
