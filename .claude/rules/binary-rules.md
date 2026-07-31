# Binary Requirements, CLI & Client/Agent Rules (PART 7, 8, 33)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never enable CGO — `CGO_ENABLED=0` for every build, every platform
- Never build on the host — Docker/Incus only (container-only development, PART 1)
- Never emit ANSI escapes when `NO_COLOR` is set (non-empty) or `TERM=dumb` (`TERM=dumb` also forces CLI mode)
- Never ship a binary that needs external files to start — single self-contained binary, first run works with zero config
- Never rename the binaries — server = `tabssh`, client = `tabssh-cli` (required companion), agent = `tabssh-agent` (optional)

## CRITICAL - ALWAYS DO

- Target `linux/amd64` + `linux/arm64` (plus other PART 3 platform-table targets) with reproducible container builds
- Implement the full server CLI surface per PART 8 (flags, `--mode`, `--debug`, `--service` install/uninstall, completions, help)
- Respect mode/debug detection priority (see config-rules.md)
- Client and agent binaries follow PART 33's requirements

## Key Rules Summary

| Binary | Name | Role |
|--------|------|------|
| server | `tabssh` | main binary, runs as service |
| client | `tabssh-cli` | REQUIRED companion (CLI/TUI/GUI) |
| agent | `tabssh-agent` | optional, runs on remote machines |

## Status

The NEVER/ALWAYS items above are derived from PART 0-6 content only. PART 7, 8, 33 bodies have not yet been read — read them directly and expand this file's sections before or during the first task that touches this domain.

For complete details, see AI.md PART 7, 8, 33
