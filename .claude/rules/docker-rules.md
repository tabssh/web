# Docker Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- NEVER place `Dockerfile` or `docker-compose.yml` in the project root — always under `docker/`
- NEVER modify `ENTRYPOINT` or `CMD` in the Dockerfile — all customization goes through `entrypoint.sh`
- NEVER bake `MODE` or `DEBUG` into the image or into production `docker-compose.yml` — the binary defaults to production; only `docker-compose.dev.yml` and `docker-compose.test.yml` set these
- NEVER include `build:` or `version:` keys in docker-compose files
- NEVER use list-style environment variables (`- KEY=value`) in compose — always map style (`KEY: value`)
- NEVER use `.env`, `.env.example`, or `.env.sample` files — all env vars must be hardcoded with sane inline defaults (`${VAR:-default}`)
- NEVER add `LABEL` blocks to the Dockerfile — all OCI metadata is applied at build time via CI `labels:`/`annotations:` (docker/metadata-action), not baked into the Dockerfile
- NEVER run `docker compose` from the project directory, and NEVER use `--project-directory` pointing at project root, and NEVER mount volumes to `{project_root}/volumes/`
- NEVER commit runtime `./volumes/` content from local runs — only `docker/rootfs/` (build-time overlay) is committed
- NEVER push `:dev` or `:test` image tags to the production registry
- NEVER use Alpine as the AIO (`Dockerfile.aio`) runtime base — Alpine/musl is insufficient for PostgreSQL/Valkey/Tor; use `debian:latest`
- NEVER use `docker/docker-compose.dev.yml` as an AI/automated agent — it is explicitly HUMAN USE ONLY, as is production `docker-compose.yml`
- NEVER expose database or cache ports externally in AIO — only port 80 is exposed; PostgreSQL (5432) and Valkey (6379) are internal-only

## CRITICAL - ALWAYS DO
- ALWAYS use `docker/` for Dockerfile(s)/compose files and `docker/rootfs/` for the build-time container overlay
- ALWAYS use a multi-stage build: `casjaysdev/go:latest` builder stage, `alpine:latest` runtime stage (standard image)
- ALWAYS use `tini` as init (`ENTRYPOINT ["tini", "-p", "SIGTERM", "--", "/usr/local/bin/entrypoint.sh"]`)
- ALWAYS set `STOPSIGNAL SIGRTMIN+3` for graceful shutdown
- ALWAYS set `HEALTHCHECK --start-period=90s --interval=10s --timeout=5s --retries=3` calling `{binary} --status`
- ALWAYS expose internal port 80 only
- ALWAYS keep `entrypoint.sh` minimal: set env, optionally start auxiliary services, exec the binary, handle shutdown signals — the binary itself handles directories, permissions, user/group, and Tor setup
- ALWAYS mount exactly two volumes in compose: `./volumes/config:/config:z` and `./volumes/data:/data:z` (use `:z` in production/test, omit `:z` in dev since it runs from a temp dir)
- ALWAYS include the `x-logging: &default-logging` anchor (`json-file`, `max-size: "5m"`, `max-file: "1"`) and reference it (`logging: *default-logging`) on every service
- ALWAYS name services/containers per convention: main service `{project_name}` / container `{project_name}-app`; db `{project_name}-db`; cache `{project_name}-cache`; search `{project_name}-search`; queue `{project_name}-queue`; proxy `{project_name}-proxy`
- ALWAYS run `docker compose` from an isolated temp directory (`{ostempdir}/{project_org}/{project_name}-XXXXXX/`), copying the compose file and creating `volumes/` fresh — never from the project directory
- ALWAYS bind production ports to the Docker bridge IP (`172.17.0.1:{port}:80`); development binds to all interfaces (`{port}:80`)
- ALWAYS build release images for both `linux/amd64` and `linux/arm64`, pushed only to `{PLATFORM_CONTAINER_REGISTRY}/{project_org}/{internal_name}`
- ALWAYS prefer the project's `tests/run_tests.sh` / `tests/docker.sh` scripts over invoking `docker-compose.test.yml` directly; direct invocation is only a fallback

## Key rules summary
- Container directory layout: `/config/{project_name}/` (server.yml, ssl/, tor/), `/data/{project_name}/` (security/geoip, uploads, cache, tor/), `/data/db/{sqlite,postgres,valkey}/`, `/data/log/{project_name}/`, `/data/backups/{project_name}/`
- SQLite databases are always named `server.db` and `users.db`, always under `/data/db/sqlite/`
- Required OCI labels applied by CI (not Dockerfile `LABEL`): maintainer, vendor, authors, title, base.name, description, licenses, created, version, schema-version, revision, url, source, documentation, vcs-type, `com.github.containers.toolbox=false`
- Image variants: standard (`:latest`, alpine, app-only via `docker/Dockerfile`) vs All-in-One (`:latest-aio`, debian, app+postgres+valkey+tor via `docker/Dockerfile.aio`, managed by supervisord)
- AIO uses PostgreSQL (unix socket only, no TCP) and Valkey (unix socket, AOF persistence) for lower footprint and no external DB/cache port exposure
- Tor is auto-enabled in containers whenever the `tor` binary is present — no `ENABLE_TOR` flag; binary owns all Tor setup (see PART 32)
- Image tags: release tags are `:latest`, `:{version}`, `:{YYMM}`, `:{commit}` (7-char) pushed to the platform registry; dev/local tags (`:dev`, `:test`) are local-only, never pushed
- Entrypoint environment defaults: `TZ=America/New_York`, `MODE=development` (compose-set only, never baked in), `DEBUG=false`, `ADDRESS=0.0.0.0`, `PORT=80`; boolean env vars accept true/yes/1/on/enable case-insensitively
- Container detection: `/.dockerenv`, `/run/.containerenv`, `/dev/lxc`, `container`/`KUBERNETES_SERVICE_HOST` env vars, init systems (tini, dumb-init, s6-svscan, runsv, runsvdir, catatonit), and cgroup contents (docker/kubepods/lxc)

Reference: AI.md PART 27 — Docker
