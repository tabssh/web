# TabSSH Web — Project-Specific Rule Overrides

This file overrides rules from AI.md and the global CLAUDE.md.
**SPEC.md wins over everything.** Only add entries here when a rule must actively
differ from the template or global default, or to declare project-wide contracts
the sibling apps (`../android`, `../desktop`) implement against.

## Activated optional PARTs

- **PART 34 (Multi-User): ACTIVE** — per-user accounts, isolated vaults, user
  API tokens
- **PART 35 (Organizations): ACTIVE** — org accounts, membership roles, org
  shared vaults, org API tokens
- **PART 36 (Custom Domains): NOT ACTIVE**

## Active overrides

None. All AI.md and global CLAUDE.md rules apply as written. The sections below
are additive contracts, not rule contradictions.

## Cross-app sync contract (implemented by android + desktop)

This section is the authoritative contract for how the sibling apps sync with a
self-hosted TabSSH Web server. The apps read THIS file for the sync behavior;
route shapes and server internals follow AI.md (PART 14 API structure, PART 11
tokens/security).

### Authentication — tokens, never username/password

- Apps authenticate with an **API token only** (`Authorization: Bearer {token}`);
  username/password is for the web UI, never for app sync
- Token format per AI.md PART 11: `usr_{32 alphanumeric}` maps to a user vault,
  `org_{32 alphanumeric}` maps to an org shared vault — the token alone selects
  which vault syncs; apps never send a user/org identifier separately
- Tokens are created, rotated, and revoked in the web UI (per AI.md: show-once,
  SHA-256 stored, multiple named tokens per owner); sync tokens SHOULD be
  created with scope `read-write`, expiry per user choice (default never)
- Token delivery to an app: pasted manually, or carried inside the existing
  encrypted QR pairing envelope (pairing payload gains an optional
  `server_url` + `sync_token` pair; both stripped fields remain
  passphrase-encrypted in transit as today)
- Apps store the token in their platform secret store (Android Keystore / OS
  keychain), same rules as passwords — never plaintext on disk
- Rotation: a rotated token invalidates the old value immediately; apps surface
  a re-auth prompt on HTTP 401 and never retry-loop a dead token

### Sync model — server-blind blob store

- The wire payload is the existing **TABSSH_SYNC_V2** encrypted package,
  unchanged (PBKDF2-HMAC-SHA256 100k → AES-GCM; GZIP JSON `SyncDataPackage`).
  The server stores and serves ciphertext only; the sync passphrase never
  leaves the apps; merge logic stays client-side
- One **canonical package per vault** (user token → personal vault package,
  org token → org vault package) plus a server-assigned monotonically
  increasing **revision** and content hash
- Sync cycle (each app, unchanged 30s debounce after local change):
  1. **Pull**: fetch canonical package + revision (skip download when revision
     matches the last-seen revision)
  2. **Merge**: decrypt, three-way merge against the local base exactly as in
     today's file-based sync (per-entity versioning, `modifiedAt`,
     `syncDeviceId`, tombstones)
  3. **Push**: upload the merged package conditionally on the pulled revision
     (compare-and-swap). If the server revision moved meanwhile, the push is
     rejected and the app repeats pull → merge → push
- Deletes travel as tombstones inside the package, as today; the server never
  interprets package contents
- Per-entity sync toggles remain a client-side concern (excluded entities are
  simply absent from the package the app builds)

### Server unavailable — offline-first (home-LAN case)

The server is **additive**: app function never depends on it.

- Unreachable server (host down, phone off-LAN, DNS gone): the app continues
  fully offline on its local database; local edits accumulate against the last
  synced base
- Retry: exponential backoff, resuming immediately on network-change events;
  no user-blocking errors — sync state is a status indicator, not a dialog
- On reconnect the normal pull → merge → push cycle reconciles everything;
  no special catch-up protocol exists or is needed (three-way merge already
  handles divergence)
- **Coexistence**: BYO-storage sync (Nextcloud/syncthing/rclone/SAF) remains
  supported and may run alongside server sync; both funnel through the same
  merge engine, so order does not matter
- Apps accept a user-configured server URL; LAN use with a self-signed or
  private-CA certificate is supported via the apps' existing TOFU certificate
  pinning (pin on first use, hard-fail on change until re-approved)

### Non-negotiables

- TLS required for the sync endpoint; no plaintext HTTP fallback
- Server never holds a decryption key, never logs token values (token ID hash
  only), never persists session or package plaintext
- A future wire-format bump is a coordinated change across all three repos —
  the web project never forks TABSSH_SYNC_V2 unilaterally

## Notes

- `internal_name: tabssh` is frozen — on-disk paths (config, data dirs) and
  service names use `tabssh`, never `web`
- The repository is `tabssh/web`; `project_name: tabssh` per IDEA.md
  `## Project variables`
