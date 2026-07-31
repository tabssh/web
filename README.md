# TabSSH Web

[![License](https://img.shields.io/github/license/tabssh/web)](LICENSE.md)

## About

TabSSH Web is the full TabSSH application in your browser — the web sibling
of TabSSH Android and TabSSH Desktop, and a self-hosted replacement for
services like Termius and Termix. One self-contained binary gives you tabbed
SSH/Telnet terminal sessions, SFTP, tunnels, hypervisor and cloud-provider
management, VNC, and host monitoring, plus the hosting role only a server
can play: server-blind end-to-end encrypted sync (`TABSSH_SYNC_V2`) and
device pairing for the native TabSSH apps.

## Official Site

https://tabssh.github.io

## Features

- Browser-style tabbed SSH and Telnet sessions — full xterm-256color and
  true-color terminal, split panes, broadcast input, find-in-scrollback,
  tmux/screen/zellij auto-attach
- SSH auth parity with the native apps: password, public key
  (RSA/ECDSA/Ed25519), keyboard-interactive, OpenSSH certificates,
  ProxyJump chains, agent forwarding, port knocking; in-app key generation
  and import (OpenSSH, PEM, PKCS#8, PuTTY)
- Host key TOFU verification with MITM detection and three trust levels
- SFTP file browser with remote editing and chmod
- Port forwarding: local, remote, and dynamic SOCKS5 tunnels that outlive
  their tabs
- Hypervisor management: Proxmox VE, XCP-ng/XenServer, Xen Orchestra,
  VMware ESXi/vCenter, OCI, libvirt/QEMU over SSH — VM power/snapshot ops
  with serial, VNC, and SPICE consoles as tabs
- Cloud host import: DigitalOcean, Hetzner, Linode, Vultr, AWS EC2, Google
  Compute, Azure, OCI
- Direct VNC hosts and reusable VNC identities
- Two-tier host monitoring: always-on availability checks plus live-session
  performance metrics with thresholds, alerts, and a multi-host dashboard
- Connection profiles, identities, groups, workspaces, snippets, macros,
  and the shared 23-theme terminal theme set — byte-compatible with the
  Android and desktop apps
- Import from `~/.ssh/config`, CSV, JSON, PuTTY .reg, and Terraform `.tf`;
  portable encrypted backup/restore interchangeable with the native apps
- Vault security: tiered credential access, auto-lock idle TTL, clipboard
  auto-clear, in-memory credential zeroing
- Accessible by design: WCAG 2.1 AA, full keyboard navigation,
  high-contrast and large-text modes, mobile-responsive UI
- Multi-user accounts with organization and team support
- Server-blind, end-to-end encrypted sync of connection vaults
  (`TABSSH_SYNC_V2`) shared with the TabSSH desktop and Android clients,
  with token-based app authentication and QR device pairing
- Organization-shared vaults with role-based access (member / admin / owner)
- Opt-in encrypted session recording and replay
- Full web UI, installable PWA, and JSON API for every feature
- Single self-contained binary, SQLite by default, zero-config first run,
  no telemetry, no feature gating

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
(cluster mode, with `server.yml` as a local cache/backup).

### Runtime environment variables

Honored on every start:

| Env Var | Default | Description |
|---------|---------|-------------|
| `MODE` | `production` | `production` or `development` |
| `DEBUG` | `false` | Verbose logging and diagnostics (accepts `1/0`, `yes/no`, `true/false`) |
| `DOMAIN` | auto | FQDN override (highest priority for hostname resolution) |
| `TABSSH_PORT` | config value | HTTP port override (below `--port`, above config) |
| `NO_COLOR` | unset | Disable ANSI color output when set and non-empty |
| `TERM` | terminal | `TERM=dumb` disables ANSI escapes and forces plain output |

### Init-only environment variables (first run only)

Applied once when `server.yml` is created, then ignored on later starts:

| Env Var | Maps to | Description |
|---------|---------|-------------|
| `PORT` | `server.port` | Initial HTTP port (otherwise a random 64xxx port is picked) |
| `LISTEN` | `server.address` | Initial listen address |
| `APPLICATION_NAME` | `server.branding.title` | Initial branding title |
| `APPLICATION_TAGLINE` | `server.branding.tagline` | Initial branding tagline |
| `CONFIG_DIR` | config directory | Overrides the platform config directory |
| `DATA_DIR` | data directory | Overrides the platform data directory |
| `CACHE_DIR` | cache directory | Overrides the platform cache directory |
| `LOG_DIR` | log directory | Overrides the platform log directory |
| `DATABASE_DIR` | database directory | Overrides the SQLite database directory |
| `BACKUP_DIR` | backup directory | Overrides the backup directory |
| `PID_FILE` | PID file path | Overrides the PID file location |

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
