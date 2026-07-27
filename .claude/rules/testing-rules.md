# Testing, Documentation & I18N/A11Y Rules (PART 29, 30, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Never run Go binaries or builds directly on the host — always inside Docker (`casjaysdev/go:latest`) or Incus (preferred, full systemd)
- Never run `reboot`, `systemctl` state changes, `iptables`, `mount`, package installs, or network reconfiguration on the host during tests — always inside a container, VM, chroot, or network namespace
- Never use the project directory for test/runtime data — all runtime/test data goes in a temp directory
- Never use bare `/tmp` — always the `/tmp/{project_org}/{internal_name}-XXXXXX/` structure
- Never commit `./volumes/` or generated runtime config files (e.g. `server.yml`) into the source repo
- Never let the AI touch `docker-compose.yml` or `docker-compose.dev.yml` — only `docker-compose.test.yml`, invoked via `tests/` scripts, is allowed
- Never put Go unit tests (`*_test.go`) under `tests/` — they stay next to the code they test
- Never treat Phase 2 binary validation (`./tests/*.sh`) as a `make test` gate — Phase 1 (toolchain gate, `*_test.go`, ≥60% coverage) is what `make test` runs; Phase 2 is manual/developer-initiated, 100% endpoint coverage
- Never leave a new user-facing string untranslated — every new string gets an i18n key across all 7 required languages
- Never hardcode `lang="en"` on `<html>` — always `lang="{{.Lang}}" dir="{{.Dir}}"`
- Never ship a UI without WCAG 2.1 AA compliance (skip links, ARIA patterns, focus management, 4.5:1 contrast minimum)
- Never add ReadTheDocs/MkDocs content outside `docs/`

## CRITICAL - ALWAYS DO
- All builds/tests run inside Docker or Incus, never on host
- `tests/run_tests.sh` (auto-detect), `tests/docker.sh`, `tests/incus.sh` (full OS with systemd, preferred) exist as the required integration-test entry points
- `make test` (Phase 1) must pass — Go unit tests via `*_test.go`, ≥60% coverage — before every commit
- ReadTheDocs/MkDocs Material docs live in `docs/`: `index.md`, `installation.md`, `configuration.md`, `api.md`, `cli.md`, `admin.md`, `security.md`, `integrations.md`, `development.md`, `stylesheets/` (dark + light themes), `requirements.txt`
- All 7 required languages are kept in sync: en, es, zh, fr, ar, de, ja, with a documented fallback chain
- Locale strings are shared via `go:embed`, using CLDR plural rules
- Arabic renders RTL (`dir="rtl"`)
- New user-facing text gets a translation key added to `en.json` (base) at minimum, with all other locale files updated to match

## Key rules summary
- Two-phase testing strategy: Phase 1 = Toolchain Gate (`*_test.go`, `make test`, ≥60% coverage, CI-enforced); Phase 2 = Binary Validation (`./tests/*.sh`, 100% endpoint coverage, manual/developer-initiated, not a commit gate)
- Container-only development extends to test execution identically — no exceptions for "just running quickly on host"
- I18N fallback chain and `go:embed` shared locale files keep translations centralized rather than duplicated per binary
- A11Y baseline: skip links, ARIA landmark/roles, visible focus indicators, 4.5:1 contrast ratio, full keyboard navigation

Reference: AI.md PART 29, PART 30, PART 31 — Testing & Development, ReadTheDocs Documentation, I18N & A11Y
