# Service Rules (Privilege Escalation & Service Management) (PART 24, 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- NEVER select a UID/GID from the well-known reserved ID ranges when creating the service account — must pick from a safe, unreserved range
- NEVER assume a single privilege-escalation mechanism works across all OSes — detection and handling must be per-OS (Linux, macOS, BSD, Windows each differ)
- NEVER hardcode service manager type — detect and support systemd, OpenRC, SysVinit, runit, rc.d (FreeBSD), launchd (macOS), and Windows Service separately
- NEVER run the Windows service as a full Windows user account when a Virtual Service Account (VSA) is sufficient — VSA is the preferred low-privilege mechanism
- NEVER leave a service account's UID and GID mismatched — they must match when created

## CRITICAL - ALWAYS DO
- ALWAYS detect the current privilege-escalation context correctly per OS before attempting service install/uninstall/disable operations
- ALWAYS create a dedicated system/service user with matching UID/GID selected from safe ranges (avoiding reserved IDs)
- ALWAYS generate the correct service unit/script for the detected service manager (systemd unit, OpenRC init script, SysVinit script, runit run script, rc.d script, launchd plist, or Windows Service registration)
- ALWAYS provide install/uninstall/disable operations with matching help-output text for each service manager
- ALWAYS use a Virtual Service Account on Windows where applicable, rather than a full user account

## Key rules summary
- Escalation detection and service commands are implemented per-platform: Linux, macOS, FreeBSD (rc.d), and Windows each have distinct command sets and requirements documented in PART 24
- Service templates exist for every supported manager: systemd, OpenRC, SysVinit, runit, rc.d, launchd, and Windows Service — each with its own full config/script content (PART 25)
- The UID/GID selection algorithm picks values in safe ranges and explicitly avoids reserved/well-known UID/GID values (see reserved UID/GID table in PART 24)
- Windows integration uses Go code for native Windows Service support, including Virtual Service Account handling

Reference: AI.md PART 24-25 — Privilege Escalation & Service, Service Support
