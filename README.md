# TabSSH Web

[![License](https://img.shields.io/github/license/tabssh/web)](LICENSE.md)

## About

TabSSH Web is the server component of the TabSSH ecosystem — a self-hosted
web application for managing SSH connection profiles, identities, and
organization-shared vaults, with server-blind end-to-end encrypted sync to
the TabSSH desktop and Android clients.

## Official Site

https://tabssh.github.io

## Features

- Multi-user accounts with organization and team support
- Server-blind, end-to-end encrypted sync of connection vaults
  (`TABSSH_SYNC_V2`) shared with the TabSSH desktop and Android clients
- Organization-shared vaults with role-based access (member / admin / owner)
- Full web UI, installable PWA, and JSON API for every feature
- Single self-contained binary, SQLite by default, zero-config first run

## Production

### Docker (Recommended)

```bash
docker run -d \
  --name tabssh-web \
  -p 8080:8080 \
  -v tabssh-config:/config/tabssh \
  -v tabssh-data:/data/tabssh \
  ghcr.io/tabssh/web:latest
```

### Docker Compose

```bash
docker compose -f docker/docker-compose.yml up -d
```

### Binary

Download the latest release for your platform from the
[releases page](https://github.com/tabssh/web/releases), then:

```bash
./tabssh --config /etc/tabssh/tabssh/server.yml
```

## Client

See the [TabSSH CLI](https://github.com/tabssh/web) (`tabssh-cli`) for
scripting and automation access to this server's API.

### Install

```bash
go install github.com/tabssh/web/src/cmd/tabssh-cli@latest
```

### Configure

```bash
tabssh-cli config set server https://your-tabssh-instance.example.com
```

### Usage

```bash
tabssh-cli auth login
tabssh-cli vault list
```

## Configuration

Configuration lives in `server.yml` (single-instance mode) or the database
(cluster mode, with `server.yml` as a local cache/backup). Key options:

| Option | Env Var | Default | Description |
|--------|---------|---------|-------------|
| `mode` | `MODE` | `production` | `production` or `development` |
| `debug` | `DEBUG` | `false` | Verbose logging and diagnostics |
| `listen_addr` | `LISTEN_ADDR` | `:8080` | HTTP listen address |
| `database.driver` | `DB_DRIVER` | `sqlite` | `sqlite` or a supported cluster DB |

See the admin panel's Settings page for the full configuration reference.

## API

Full API documentation is served by the running instance itself:

- Swagger UI: `/server/docs/swagger` (JSON spec also at `/api/swagger`)
- GraphiQL UI: `/server/docs/graphql` (endpoint also at `/api/graphql`)
- Health check: `/server/healthz` (also `/api/{api_version}/server/healthz`)

## Other

### Troubleshooting

If the server enters maintenance mode (database or config write failure),
check the admin panel banner and server logs for the specific cause — the
server self-heals and retries automatically once the underlying issue
clears.

## Development

**For contributors only — end users should use the Production section above.**

### Prerequisites

- Docker (all builds and tests run in containers — no local Go toolchain
  is required or supported)

### Build

```bash
make dev      # quick build to a tempdir
make local    # versioned build to binaries/
make build    # full cross-platform release build
make test     # run unit tests
```

### Project Structure

```
src/            application source
docker/         Dockerfile and container assets
scripts/        developer and release scripts
tests/          integration and container test harnesses
docs/           MkDocs documentation source
```

## Disclaimer

This software is provided "as is" without warranty of any kind. Use at your own risk.

- **No Warranty**: The authors are not responsible for any damages, data loss, or issues arising from use of this software
- **Not Professional Advice**: This software does not constitute legal, financial, medical, or other professional advice
- **Third-Party Services**: If this software connects to external APIs or services, their terms of service apply separately
- **Security**: While we strive to follow security best practices, no software is guaranteed to be free of vulnerabilities

## License

MIT License — see [LICENSE.md](LICENSE.md) for details.
