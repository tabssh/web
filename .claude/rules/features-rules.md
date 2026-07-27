# Features Rules (Email, Scheduler, GeoIP, Metrics, Backup, Update) (PART 18-23)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- NEVER require an external scheduler (cron, systemd timers, external job queues) — the built-in Go time/ticker-based scheduler is the only scheduling mechanism
- NEVER use `geoip2-golang` for GeoIP lookups — use `github.com/oschwald/maxminddb-golang` against the `sapics/ip-location-db` MMDB files
- NEVER expose the `/metrics` endpoint publicly — it is internal-only (localhost/internal network access control)
- NEVER use unbounded/high-cardinality label values in Prometheus metrics (e.g. raw user IDs, raw IPs, raw paths) — cardinality must stay bounded and safe
- NEVER run backup retention cleanup before the new backup has been verified — verification must happen before old backups are pruned
- NEVER skip encryption on backups — backups are always encrypted with AES-256-GCM, keys derived via Argon2id
- NEVER disable or bypass backup/restore rules when compliance mode is enabled — compliance mode enforcement is mandatory
- NEVER perform the self-update binary replacement without verifying the SHA256 checksum of the downloaded binary first
- NEVER use the Windows in-place rename-over-running-binary approach used on Unix — Windows requires rename + reboot-delete of the old binary since it cannot be overwritten while running
- NEVER skip the WebUI-vs-Email notification decision matrix — some events are WebUI-only, some are email-eligible, per the documented matrix

## CRITICAL - ALWAYS DO
- ALWAYS support SMTP auto-detection/auto-configuration for email, with embedded default templates that can be overridden per-instance
- ALWAYS store custom email template overrides in `{config_dir}/template/email/`, falling back to embedded defaults when no override exists
- ALWAYS persist scheduler job state in the database and make the scheduler cluster-aware (uses DB-backed locking so only one node in a cluster runs a given job at a time)
- ALWAYS use Prometheus-compatible metric naming conventions and format for the `/metrics` endpoint
- ALWAYS verify a backup archive's integrity before it can be relied upon for retention cleanup or restore
- ALWAYS support cumulative update channels: stable, beta, and daily — updates apply cumulatively within/between channels as documented
- ALWAYS verify the SHA256 checksum of any downloaded update binary before applying it
- ALWAYS use the platform-correct binary replacement strategy for self-update: atomic rename on Unix; rename + delete-old-on-reboot on Windows

## Key rules summary
- Email system: embedded defaults + custom override directory pattern (`{config_dir}/template/email/`) lets admins customize without losing upgrade compatibility
- Notifications: WebUI notifications and Email notifications are distinct channels governed by a decision matrix — not every event goes to both
- Scheduler: Go `time`/`ticker` based, no external cron dependency, DB-persisted job state, cluster-safe via locking
- GeoIP: data source is `sapics/ip-location-db`; the Go library is `github.com/oschwald/maxminddb-golang`, not `geoip2-golang`
- Metrics: Prometheus-compatible `/metrics` endpoint, internal-access-only, naming conventions matter, cardinality safety is a hard requirement
- Backup/restore: AES-256-GCM encryption + Argon2id key derivation; verify-before-cleanup lifecycle; compliance-mode enforcement; cluster-aware backup rules
- Update command: platform-specific binary replacement (atomic rename on Unix vs rename+reboot-delete on Windows), mandatory SHA256 checksum verification, cumulative stable/beta/daily channels

Reference: AI.md PART 18-23 — Email & Notifications, Scheduler, GeoIP, Metrics, Backup & Restore, Update Command
