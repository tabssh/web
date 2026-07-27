# TabSSH Web

## Project description

TabSSH Web is the self-hosted web sibling of TabSSH Android and TabSSH Desktop. It
brings TabSSH's browser-style tabbed SSH experience to any web browser and, at the
same time, acts as the hosted companion service the native apps have been hinting
at: an end-to-end-encrypted sync endpoint for the shared TABSSH_SYNC_V2 format and
a device-pairing broker for the existing QR pairing flow.

TabSSH Web is the FULL TabSSH application in a browser, not a thin sync
companion. Everything the Android and desktop apps do is here: tabbed
SSH/Telnet terminal sessions, SFTP, port forwarding and tunnels, hypervisor
management with serial/VNC consoles, cloud-provider host import, direct VNC
hosts, background host monitoring with performance dashboards, snippets,
macros, workspaces, and the shared theme set — plus the hosting role only a
server can play: the sync endpoint and pairing broker for the native apps.
The entity set is the same and byte-compatible, so one encrypted data set
roams across all three. Like its siblings it is a single self-contained
binary, works on first run with zero config, sends no telemetry, gates no
features, and never stores a secret in plaintext.

TabSSH Web is multi-user with organization support and is positioned as a
self-hosted replacement for commercial/hosted SSH managers such as Termius and
Termix: teams get shared connection vaults, org-scoped roles, cross-device sync,
and a full browser terminal — without a third-party cloud holding their keys.

Target users: developers and sysadmins who already use TabSSH on Android or
desktop and want browser access from machines where nothing can be installed;
teams and organizations self-hosting one TabSSH Web instance as their shared
sync point and team SSH manager; users leaving Termius/Termix who want the same
workflow without subscriptions, feature gating, or third-party key custody.

## Project variables

project_name:     tabssh
project_org:      tabssh
# FROZEN — set once at first-time setup, never edit
internal_name:    tabssh
# FROZEN — set once at first-time setup, never edit
internal_org:     tabssh
app_name:         TabSSH Web
official_site:    tabssh.github.io
maintainer_name:  casjay
maintainer_email: casjay@yahoo.com
license:          MIT
repo:             https://github.com/tabssh/web
android_sibling:  ../android
desktop_sibling:  ../desktop

## Business logic

### Product scope & non-goals

**In scope:**

- **Web terminal**: browser-style tabbed SSH and Telnet sessions (multiple
  concurrent tabs, same-host allowed), full xterm-256color and 24-bit
  true-color emulation, UTF-8, configurable scrollback with
  find-in-scrollback, per-tab live titles, split panes, broadcast/cluster
  input, per-profile tmux/screen/zellij auto-attach, post-connect scripts,
  per-connection env vars
- **SSH auth parity**: password, public key (RSA/ECDSA/Ed25519), keyboard-interactive;
  key import (OpenSSH, PEM, PKCS#8, PuTTY) and in-app key generation; SHA-256
  fingerprints; OpenSSH user certificates; ProxyJump/jump-host chains; agent
  forwarding; keepalive and auto-reconnect matching sibling rules (never
  retry auth failures); port knocking pre-connect sequences
- **Host key trust**: TOFU verification with MITM detection and the sibling apps'
  three trust levels (UNKNOWN / ACCEPTED / VERIFIED), tracked per user
- **SFTP**: file browser (upload/download/rename/chmod/delete) and a remote text
  editor with a binary-file guard, matching sibling behavior
- **Profile management**: connections, identities, stored keys, connection groups,
  workspaces, snippets (with `{?variable}` placeholders and run-time prompt
  UI), macros, and themes — the same entity set the siblings sync, editable
  from the browser
- **Import / export**: `~/.ssh/config` import; bulk import from CSV, JSON,
  PuTTY .reg, and Terraform `.tf` files, mapping fields onto connection
  profiles exactly as the siblings do; portable encrypted backup/restore
  archive (sibling BACKUP_VERSION 3, forward-migrating)
- **Vault security**: tiered credential access matching the siblings
  (never / session-only / encrypted), vault auto-lock with idle TTL,
  clipboard auto-clear for sensitive copies, in-memory credential zeroing
- **Accessibility**: WCAG 2.1 AA, full keyboard navigation, high-contrast
  and large-text modes, screen-reader labels on all interactive elements;
  mobile-responsive from day one
- **Port forwarding / tunnels**: local, remote, and dynamic SOCKS5 forwards
  terminating on the server, with tunnels that outlive terminal tabs —
  the server-hosted equivalent of the siblings' tunnel manager
- **Hypervisor management**: Proxmox VE, XCP-ng/XenServer, Xen Orchestra,
  VMware ESXi/vCenter, Oracle Cloud OCI, and libvirt/QEMU over SSH (no
  exposed libvirt TCP daemon); VM power/snapshot operations; serial, VNC,
  and SPICE consoles opening as tabs next to terminal sessions; reusable
  hypervisor credential accounts; per-profile TLS certificate pinning
  (TOFU, off by default for self-signed hypervisor certs), matching
  Android's behavior
- **Cloud host import**: DigitalOcean, Hetzner, Linode, Vultr, AWS EC2,
  Google Compute, Azure, OCI — enumerate instances, create connection
  profiles, start/stop/restart, one-click SSH connect
- **VNC hosts**: direct VNC connections (no SSH) as tabs, with reusable VNC
  identities, matching the siblings' RFB client behavior
- **Host monitoring**: two-tier like the siblings — always-on availability
  checks (TCP reachability, no credentials) for hosts the user opts in,
  plus live-session performance metrics (CPU/mem/disk/load) with per-host
  monitor slots, thresholds, down/recovery notifications, and a multi-host
  dashboard
- **Session recording**: opt-in per user, transcripts stored encrypted in
  the user's vault, replayable and format-compatible with the siblings;
  never enabled silently (see Security decisions)
- **Themes**: the shared built-in terminal theme set, byte-compatible with the
  sibling theme format; dark/light/auto; custom theme import/export with
  contrast validation
- **E2E-encrypted sync service**: store, list, and serve TABSSH_SYNC_V2 blobs per
  account; blobs are opaque ciphertext produced with the user's sync passphrase;
  the server never possesses a decryption key; per-entity versioning, tombstones,
  and conflict semantics are honored by clients, not re-implemented server-side
- **Device pairing broker**: short-lived rendezvous for the existing QR pairing
  envelope (opaque encrypted payload, 6-digit code, 60-second TTL), letting
  Android and desktop pair through this server instead of camera-to-screen only
- **Multi-user**: per-user isolated vaults and sync storage; admin-managed
  accounts (see AI.md PART 34 for the multi-user pattern)
- **Organizations**: org accounts with member roles and org-scoped shared
  vaults — shared connections, identities, keys, snippets, and themes visible
  to members per their org role (see AI.md PART 35 for the organizations
  pattern); the Termius/Termix team-vault use case

**Non-goals:**

- Mosh — its UDP roaming client is inherently native; browser sessions
  reconnect via the server instead
- X11 forwarding — the browser has no local display server to forward to
- Local-device forwards — a "local" forward binds on the server, never on
  the browser's machine; true local-machine binds remain a native-app feature
- Silent/admin-initiated session recording — recording is user-opt-in only
  (see Security decisions)
- Being a third-party-operated public SaaS — this is a self-hosted product; an
  instance serves the users and organizations its operator admits

### Roles & permissions

- **Guest (unauthenticated)**: no access beyond login/setup and health
- **User**: owns a private vault (profiles, identities, keys, themes, snippets,
  macros, groups, workspaces, host-key trust store) and private sync-blob
  storage; opens terminal/SFTP sessions only with their own profiles; can never
  see another user's vault, blobs, or sessions
- **Org member**: sees and uses the org's shared vault entities per org role;
  personal vault stays private regardless of org membership
- **Org admin**: manages org membership, org roles, and the org shared vault;
  no access to members' personal vaults, blobs, or sessions
- **Org owner**: org admin plus org lifecycle (rename, transfer, delete)
- **Instance admin**: manages accounts, organizations, server settings, quotas,
  and instance health; can NOT read any user's or org's vault contents, sync
  blobs (ciphertext only, and no keys exist server-side), or live session
  content
- First account created at first run is the instance admin (zero-config first
  run)

### Data model & sensitivity

Entities mirror the sibling apps' shared sync entity set (Android Room ~v37 /
desktop SQLite) so sync remains lossless; web-only entities are additive.

- **User account**: email/username, password hash, role, app-lock settings —
  sensitive (credentials); password stored as hash only
- **Vault entities** (per user, stored encrypted with a key derived from the
  user's vault passphrase; server holds ciphertext only):
  - ConnectionProfile: name, host, port, username, protocol, auth type,
    identity/key reference, group, color tag, env vars, post-connect script
  - Identity: reusable credential profile (username, auth type, key/password ref)
  - StoredKey: key metadata (type, fingerprint, comment); private key material
    ciphertext
  - HostKey: hostname:port, key type, public key, fingerprint, trust level
  - TrustedCertificate: pinned TLS certs for hypervisor/cloud/VNC endpoints
    (fingerprint, trust level), TOFU like host keys
  - HypervisorAccount / HypervisorProfile: provider, endpoint, region,
    credential references
  - CloudAccount: provider, enabled flag, last refresh metadata; API token
    ciphertext
  - VncHost / VncIdentity: host, port, display, security type, credential ref
  - MonitorSlot: per-host thresholds, intervals, alert settings (host:port
    disclosed to the server scheduler when monitoring is enabled — see
    Security decisions)
  - Tunnel definitions: type (local/remote/dynamic), bind and target specs
  - Session recordings: encrypted transcripts, opt-in only
  - Theme, ConnectionGroup, Workspace, Snippet, Macro
- **Organization**: name, owner, members with org roles — metadata sensitivity
- **Org shared vault**: same entity types as the personal vault, org-scoped;
  encrypted with an org vault key shared among members client-side, server
  holds ciphertext only
- **Sync blob**: opaque TABSSH_SYNC_V2 ciphertext + owner + timestamps — high
  sensitivity but server-unreadable by design
- **Pairing envelope**: opaque encrypted payload + code + expiry — ephemeral,
  deleted on claim or TTL expiry, never persisted past its 60-second lifetime
- **Session state**: live terminal/SFTP traffic — highest sensitivity, transient,
  in memory only, never written to disk
- **Audit events**: auth and administrative actions only (no command or session
  content), per AI.md PART 11

### Trust boundaries & external services

- **Browser client (after vault unlock)**: trusted with plaintext vault data and
  session content; vault passphrase never leaves the client for storage purposes
- **The server itself**: trusted relay for live SSH traffic — the product
  assumption is that the operator self-hosts and trusts their own instance;
  failure mode: if the server is down, native apps continue fully offline
  (BYO-storage sync remains their fallback)
- **Target SSH/Telnet/VNC hosts**: untrusted until TOFU acceptance; host key
  or certificate change after trust is treated as potential MITM and blocks
  connection pending explicit user re-approval
- **Hypervisor and cloud-provider APIs** (Proxmox, XCP-ng, Xen Orchestra,
  VMware, OCI, libvirt; DigitalOcean, Hetzner, Linode, Vultr, AWS, GCP,
  Azure): reached outbound over TLS with per-endpoint TOFU/pinned
  certificates; their responses are untrusted input (strict parsing, no
  execution); credentials are supplied from the client-decrypted vault at
  call time and held in memory only; failure mode: provider unreachable →
  the affected panel degrades, terminal/SFTP/sync remain fully functional
- **Sibling apps (Android/desktop)**: trusted peers only via the shared sync
  passphrase and pairing envelope encryption; the server treats their payloads
  as opaque
- **External services**: none required. Optional outbound SMTP for notifications
  (AI.md PART 18); failure mode: notification loss only, never functional loss.
  No vendor SDKs, no analytics services

### Threat model & abuse cases

**Primary assets**: SSH private keys and passwords, live session plaintext, sync
blobs, user account credentials, host-key trust store.

**Trusted vs untrusted inputs**: authenticated users' own vault operations are
trusted after unlock; everything else is untrusted — login attempts, uploaded
sync blobs (opaque, size-capped), pairing payloads, key/profile imports, all
terminal input/output, and every value returned by target SSH hosts.

**Attacker goals**: steal credentials or keys, read or tamper with sessions,
harvest sync blobs for offline cracking, take over accounts, or abuse the
gateway as an attack proxy.

**Abuse cases and required defenses:**

- **Credential stuffing / brute force on login**: rate limiting, lockout, and
  enumeration-safe identical auth errors (AI.md PART 11)
- **Vault passphrase offline cracking via stolen blobs**: strong KDF is mandated
  by the shared format (PBKDF2-HMAC-SHA256, 100,000 iterations); blob access
  requires authentication; downloads are rate-limited
- **Gateway abuse as SSH scanning/spam proxy**: per-user concurrent-session and
  connection-rate limits; admin-configurable destination allow/deny lists
- **SSRF via user-supplied host:port and API endpoints**: applies to SSH,
  Telnet, VNC, tunnels, hypervisor endpoints, and cloud API base URLs alike —
  outbound connections to loopback, link-local, and cloud metadata ranges are
  denied by default (admin override for genuine LAN use — see Security
  decisions)
- **Cloud/hypervisor credential theft**: provider tokens live only as vault
  ciphertext at rest and in server memory during an operation; never logged,
  never in audit events, zeroed after use
- **Pairing-code guessing**: 60-second TTL, single claim, attempt rate limiting;
  payload remains useless without the pairing passphrase
- **Malicious import files** (keys, profiles, themes): strict parsing, size
  caps, no execution of imported content
- **Privilege escalation across users/orgs**: strict per-user and per-org
  isolation; org role checks on every shared-vault operation; IDOR guards on
  all object references (AI.md PART 11); a departing member's org access ends
  immediately on removal
- **Session hijacking**: standard session security, CSRF protection, and
  security headers per AI.md PART 11
- **Storage exhaustion**: per-user quotas on sync blobs and vault size

**Explicit non-goal**: defending a user against their own server operator for
live sessions — see below.

### Security decisions & exceptions

- **Live-session relay**: unlike sync (E2E, server-blind), interactive
  SSH/Telnet/VNC sessions, tunnels, and hypervisor/cloud API calls
  necessarily transit the server, which briefly holds session plaintext and
  credentials in memory. Accepted because the product is self-hosted — the
  operator is the user or their team. Mitigations: never persisted, no silent
  server-side recording, memory zeroed on session close
- **Session recording is user-opt-in only**: transcripts are encrypted into
  the recording user's vault; admins cannot enable recording for others —
  silent recording is a non-goal, not a config option
- **Monitoring metadata disclosure**: enabling availability monitoring for a
  host discloses that host:port (not credentials) to the server's scheduler
  in plaintext config — accepted and documented; performance checks run only
  through live, user-initiated sessions, matching the siblings' two-tier
  model
- **Outbound SSH is the product function**: the gateway intentionally opens
  user-directed outbound connections (inherently SSRF-adjacent). Mitigated by
  the default internal-range denylist and admin-controlled overrides for
  legitimate LAN/homelab targets
- **Server-blind sync**: sync blobs and pairing envelopes are stored and relayed
  without any server-held key — losing the vault/sync passphrase means the data
  is unrecoverable; this is intentional and matches the siblings
- **No telemetry, no feature gating, opt-in only for anything that phones out** —
  matching the sibling apps' invariants

**Business rules / invariants:**

- **Single ecosystem, Android as reference**: the Android app is currently
  the most feature- and code-complete member; web tracks Android's feature
  surface (minus device-local features like Mosh, X11, Tasker, on-screen
  keyboards), and desktop converges on the same surface as it is built —
  divergence from Android behavior is a bug unless documented here
- Sync wire format is byte-compatible with TABSSH_SYNC_V2 as shipped by the
  desktop and Android apps; the web project never forks the format
- Themes, backup archives (BACKUP_VERSION 3), and QR pairing payloads remain
  interchangeable with the siblings — Android's implementations are
  authoritative for these formats
- Protocol compatibility matches the siblings: SSH2 (RFC 4251-4254), Telnet
  (RFC 854 + ECHO/SGA/TERMINAL-TYPE/NAWS), VNC (RFB 3.8), SHA-256 host-key
  fingerprints with Android's emoji visual fingerprint format
- All eight cloud providers expose the same feature surface — list, live
  state, power control, SSH connect; no provider gets a reduced experience;
  no vendor SDKs, REST APIs only
- No secret (password, private key, token) is ever stored in plaintext — at
  rest everything sensitive is ciphertext; account passwords are hashed
- Host-key trust is per user, TOFU, three levels, MITM-change blocks
- First run works with zero config; fully functional on an offline LAN
- Server never re-implements client-side merge: conflict resolution stays in
  the clients that own the plaintext

**Endpoints (WHAT, not paths — see AI.md PART 14 for HOW):**

- Account: first-run admin setup, login/logout, password change; API token
  management (create / rotate / revoke; `usr_`/`org_` tokens are how the
  sibling apps authenticate for sync — see SPEC.md cross-app sync contract)
- Vault: CRUD for connections, identities, keys, themes, snippets, macros,
  groups, workspaces; key generation; host-key trust management
- Sessions: open interactive terminal session (SSH/Telnet), resize, close;
  SFTP browse, transfer, edit, chmod; VNC and hypervisor serial/VNC console
  sessions; opt-in recording start/stop and replay
- Tunnels: create / list / stop local, remote, and dynamic forwards
- Hypervisors: manage accounts/profiles, list VMs, power/snapshot ops, open
  console
- Cloud: manage provider accounts, refresh/enumerate instances, instance
  power ops, create profile from instance
- Monitoring: manage monitor slots, read availability state and metrics,
  multi-host dashboard feed
- Sync: upload / download / list / delete sync blobs (opaque)
- Pairing: create pairing session, claim by code
- Organizations: create org, manage members and roles, manage org shared vault
- Admin: user and organization management, instance settings, quotas
- Health and metrics per AI.md PARTs 13 and 21

**Data sources:**

- Sync blobs and backup archives produced by the Android and desktop apps
- Built-in terminal theme set embedded at build, shared with the siblings
- Hypervisor and cloud-provider APIs, queried on demand with user-supplied
  credentials (no vendor SDKs — plain REST/XML-RPC like the siblings)
