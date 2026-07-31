## [ ] Preserve YAML comments on config Save (audit finding)
Read: AI.md PART 5 (YAML Comment Style)

`config.Config.Save` (src/config/config.go) uses `yaml.Marshal`, which drops
every comment. First-run creation goes through `renderDefaultConfig` (comments
intact), but any later `Save` (e.g. cluster `_cache` write via
`Manager.writeCacheLocked`) rewrites the file comment-free, violating the PART 5
"ALL comments ABOVE the setting" rule for the persisted file. Needs a
comment-preserving writer (yaml.Node round-trip or a re-render) before the
first code path that calls `Save` on the human-facing `server.yml` ships. No
caller writes the human config today (PART 10 deferred), so this is latent.

## [ ] Add ReadTimeout/WriteTimeout to HTTP server (audit finding)
Read: AI.md PART 8, PART 11

`server.New` (src/server/server.go) sets `ReadHeaderTimeout` and `IdleTimeout`
but leaves `ReadTimeout` and `WriteTimeout` unset, so a slow-body client can
hold a connection open (Slowloris-style on the body). Add bounded
`ReadTimeout`/`WriteTimeout` when the request-handling PARTs (13/14) wire real
routes and streaming endpoints are known (streaming handlers need per-route
overrides, so set a sane default now and exempt streams later).

## [ ] Resolve symlinks in SafeFilePath (audit finding)
Read: AI.md PART 11 (path traversal)

`SafeFilePath` (src/config/path.go) cleans and sandbox-checks the path but does
not `filepath.EvalSymlinks`, so a symlink inside the sandbox pointing outside it
would pass. No caller passes attacker-controlled paths yet (upload/file routes
are PART 16/17, deferred), so this is latent — add symlink resolution before
the first user-facing file path reaches it.

## [ ] Populate remaining `.claude/rules/*.md` files with real content
Read: AI.md PART 7, 8, 9, 10, 11, 12, 13, 14, 15, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36

`binary-rules.md`, `backend-rules.md`, `api-rules.md`, `features-rules.md`,
`service-rules.md`, `makefile-rules.md`, `docker-rules.md`, `cicd-rules.md`,
`testing-rules.md`, and `optional-rules.md` were created as structural
placeholders during the PART 0-6 bootstrap pass (correct filenames/headers/
PART mapping per PART 0's table) but their NEVER/ALWAYS/summary sections are
not yet populated from the actual PART content — that requires reading each
PART in full. Repopulate each file the first time work touches its domain.

## [ ] Implement PART 27 Docker build (Dockerfile, compose files)
Read: AI.md PART 27

`docker/rootfs/` has the overlay directory skeleton (`config/`, `etc/`,
`usr/local/bin/`) but no `docker/Dockerfile`, `docker-compose.yml`, or
`docker-compose.dev.yml`/`docker-compose.test.yml` exist yet. Toolchain image
must be `casjaysdev/go:latest` per global convention unless AI.md/SPEC.md
says otherwise.

## [ ] Implement PART 28 CI/CD workflows
Read: AI.md PART 28

`.github/` exists but is currently empty (no `workflows/`, no `CODEOWNERS`,
`SECURITY.md`, `ISSUE_TEMPLATE/`, PR template). Jenkinsfile already exists at
root and appears to cover pipeline needs — verify whether GitHub Actions
workflows are additionally required by PART 28 before creating them.

## [ ] Implement PART 26 Makefile
Read: AI.md PART 26

No root `Makefile` exists yet. Local dev target set (build/test/lint/run)
needs to be created per PART 26's pattern once PART 7 binary requirements
and PART 27 Docker toolchain are in place.

## [ ] Implement remaining application layers (PART 7-25, 29-36)
Read: AI.md PART 7-25, 29-36

`src/` currently has: entrypoint (`main.go`, `commands.go`, `flags.go`,
`help.go`, `serve.go`, `completions.go`), and packages `common/`, `config/`,
`handler/` (stub), `middleware/` (stub), `mode/`, `model/` (stub), `paths/`,
`pid/`, `runenv/`, `server/`, `signal/`, `urlutil/`. PARTs 7 (binary
requirements), 8 (server CLI), 9-25 (error handling, DB/cluster, security,
server config, health, API, TLS, web frontend, admin panel, email,
scheduler, GeoIP, metrics, backup, update command, privilege escalation,
service support), and 29-36 (testing, docs, i18n, Tor, client/agent,
multi-user, organizations) are not yet implemented — this is expected, this
bootstrap pass covered PART 0-6 only (structure, config, modes) per the
current task scope. `tests/` and `scripts/` directories are currently empty.
