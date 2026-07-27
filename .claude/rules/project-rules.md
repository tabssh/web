# Project Structure & Licensing Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Use any license other than MIT · Skip `LICENSE.md` at project root
- Put third-party license attribution anywhere other than `LICENSE.md` once dependencies exist
- Use `github.com/mattn/go-sqlite3` (CGO) — always use `modernc.org/sqlite`
- Hardcode a specific Go version anywhere (docs, Docker, CI) — always latest stable, `casjaysdev/go:latest` unpinned
- Skip a required build platform — all 8: linux/darwin/windows/freebsd × amd64/arm64
- Create root-level `config/`, `data/`, `logs/`, `tmp/`, `temp/`, `build/`, `dist/`, `out/`, `vendor/`, `node_modules/`, `lib/`, `utils/`, `common/`
- Commit `binaries/`, `releases/`, or `volumes/` — always gitignored
- Use a config file name other than `server.yml` on any platform
- Use Docker-only paths (`/data/**`, `/config/**`) on native (non-container) platforms

## CRITICAL - ALWAYS DO
- MIT License in `LICENSE.md`, copyright holder `tabssh`, current year, exact boilerplate text from AI.md PART 2
- Attribute all third-party dependency licenses in `LICENSE.md` once dependencies are vendored
- All Go libraries must be pure Go and work with `CGO_ENABLED=0`
- Use `modernc.org/sqlite` (driver name `sqlite`, aliases `sqlite2`/`sqlite3`) for SQLite
- Follow the OS-specific path convention exactly (see table below)
- Config file is always `server.yml` at every OS path
- Docker/container paths: `/config/{project_name}/`, `/data/{project_name}/` (container-only, never used natively)
- `.gitignore` header: `# gitignore created on MM/DD/YY at HH:MM` then literal `ignoredirmessage` on line 2
- `.dockerignore` excludes `.git/`, CI workflow dirs, `volumes/`, `binaries/`, `releases/`, `tests/`, `docs/`

## Key rules summary
- Required Go libraries — DB: `modernc.org/sqlite` (sqlite), `github.com/jackc/pgx/v5/stdlib` (postgres), `github.com/go-sql-driver/mysql` (mysql), `github.com/microsoft/go-mssqldb` (mssql); Cache: `github.com/redis/go-redis/v9`; Core: `gopkg.in/yaml.v3`, `github.com/google/uuid`, `golang.org/x/crypto/argon2`, `golang.org/x/crypto/bcrypt` (verify-only, rehash to Argon2id); Auth (required even without end users): `github.com/pquerna/otp` (TOTP), `github.com/go-webauthn/webauthn` (Passkeys), `github.com/golang-jwt/jwt/v5` (JWT), `github.com/coreos/go-oidc/v3` (OIDC), `golang.org/x/oauth2`, `github.com/go-ldap/ldap/v3`, `github.com/crewjam/saml`, `github.com/gorilla/sessions`; Network: `github.com/go-chi/chi/v5` (router), `github.com/cretz/bine` (Tor), `github.com/gorilla/websocket`, `github.com/rs/cors`; Utilities: stdlib `embed`, `github.com/go-co-op/gocron/v2` (scheduler), `golang.org/x/time/rate`, `github.com/go-playground/validator/v10`
- Directory structure: `src/` holds all source (singular Go package dirs: `handler/`, `model/`, `middleware/`, etc.); `docker/` holds `Dockerfile`, `Dockerfile.dev`, `docker-compose.yml`, `docker-compose.dev.yml`, `docker-compose.test.yml`, `rootfs/` (committed, build-time overlay); `scripts/`, `tests/` (`run_tests.sh`, `docker.sh`, `incus.sh`), `docs/` (MkDocs/ReadTheDocs only); `.github/workflows/` holds `ci.yml`, `release.yml`, `beta.yml`, `daily.yml`, `docker.yml`; `binaries/`, `releases/`, `volumes/` are gitignored, never committed
- Root files: `README.md`, `LICENSE.md`, `AI.md`, `IDEA.md`, `SPEC.md`, `CLAUDE.md`, `TODO.AI.md`, `Jenkinsfile`, `release.txt`, `site.txt`, `.gitignore`, `.dockerignore`, `mkdocs.yml`, `.readthedocs.yaml`, `Makefile`, `go.mod`, `go.sum`
- OS-specific config paths (internal_org=tabssh, internal_name=tabssh):

| OS | Config (privileged) | Config (user) |
|----|---------------------|----------------|
| Linux | `/etc/tabssh/tabssh/server.yml` | `~/.config/tabssh/tabssh/server.yml` |
| macOS | `/Library/Application Support/tabssh/tabssh/server.yml` | `~/Library/Application Support/tabssh/tabssh/server.yml` |
| BSD | `/usr/local/etc/tabssh/tabssh/server.yml` | `~/.config/tabssh/tabssh/server.yml` |
| Windows | `%ProgramData%\tabssh\tabssh\server.yml` | `%AppData%\tabssh\tabssh\server.yml` |
| Docker | `/config/tabssh/server.yml` | — |

- Data/cache/log/backup/ssl/security/db directories follow the same `{internal_org}/{internal_name}` pattern per OS (see AI.md PART 4 for exact paths per type)

Reference: AI.md PART 2, PART 3, PART 4 — License & Attribution, Project Structure, OS-Specific Paths
