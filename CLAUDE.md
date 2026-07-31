# Project SPEC

Project: TabSSH Web
Role: Efficient loader for AI.md

⚠️ **THIS FILE IS AUTO-LOADED EVERY CONVERSATION. FOLLOW IT EXACTLY.** ⚠️

Purpose:
- This file is a short loader for the most important rules
- `AI.md` is the full source of truth
- `SPEC.md` overrides `AI.md` on conflict (activates PART 34 Multi-User, PART 35 Organizations; PART 36 Custom Domains NOT active)
- For complete details, read the referenced PARTs in `AI.md`

## FIRST TURN - MANDATORY

On EVERY new conversation or after "context compacted" message:
1. **READ** the relevant `.claude/rules/*.md` for your current task
2. **NEVER** assume or guess - verify against AI.md before implementing

## Asking Questions

Default to continuing without asking. Ask only when:
- The project's own spec (AI.md/IDEA.md/SPEC.md) is genuinely silent, contradictory, or missing the value
- The action is destructive/irreversible (force push, drop table, `rm -rf`)
- A real architectural choice exists with no documented default
- Business logic or feature scope is ambiguous

## Before ANY Code Change

1. Confirm the file/path is inside this project
2. Check whether SPEC.md overrides the relevant AI.md rule
3. Read the applicable `.claude/rules/*.md` file
4. Verify against actual project state (don't assume from memory)

## Binary Terminology

- **server** = `tabssh` (main binary, runs as service)
- **client** = `tabssh-cli` (REQUIRED companion, CLI/TUI/GUI)
- **agent** = `tabssh-agent` (optional, runs on remote machines)

## Key Placeholders

- `{project_name}` = tabssh
- `{project_org}` = tabssh
- `{internal_name}` / `{internal_org}` = tabssh (frozen)
- `{admin_path}` = admin (default)

## Account Types (CRITICAL)

- Server Admin: full system control
- Primary Admin: first admin account
- Regular User: standard account (Multi-User active per SPEC.md)
- Cluster nodes vs. managed nodes are distinct concepts — never conflate them

## NEVER Do (Top 19) - VIOLATIONS ARE BUGS

1. bcrypt for passwords → use Argon2id
2. Dockerfile in root → must be `docker/Dockerfile`
3. CGO enabled → always `CGO_ENABLED=0`
4. `mattn/go-sqlite3` → use `modernc.org/sqlite`
5. `gorilla/mux` → forbidden
6. `dgrijalva/jwt-go` → forbidden (unmaintained)
7. Inline YAML comments → always above the setting
8. `server.yaml` → always `server.yml`
9. Root-level `config/`, `data/`, `logs/`, `tmp/`, `vendor/`, `node_modules/` directories
10. `.env*` files committed
11. Root `SUMMARY.md`/`CHANGELOG.md`/`AUDIT.md`/`REPORT.md`
12. Modifying `AI.md`
13. Committing directly (only `gitcommit --dir {dir} all`)
14. Bypassing a PreToolUse hook block
15. GPL/AGPL/LGPL dependency in an MIT project
16. Debug mode bypassing authentication (never, in any mode)
17. Hardcoded paths instead of OS-specific path resolution (PART 4)
18. `strconv.ParseBool` instead of `config.ParseBool`
19. TODOs/stubs/commented-out code committed

## ALWAYS Do - NON-NEGOTIABLE

1. Read AI.md/IDEA.md/SPEC.md before acting
2. Argon2id for passwords, SHA-256 for API/session tokens
3. `CGO_ENABLED=0`, builds only in Docker/Incus, never on host
4. Validate and normalize every path (PathSecurityMiddleware, SafePath)
5. `config.ParseBool`/`IsTruthy` for all boolean parsing
6. MIT license with embedded-licenses table kept in sync with go.mod
7. Mode/debug priority: CLI flag > env var > default
8. Maintenance mode self-healing on critical errors (DB/disk only)
9. Log any out-of-scope issue to TODO.AI.md rather than losing it in chat

## File Locations

- Config: `/etc/tabssh/tabssh/server.yml` (root) or `~/.config/tabssh/tabssh/server.yml` (user); `/config/tabssh/server.yml` (Docker)
- Data: `/var/lib/tabssh/tabssh/` (root) or platform equivalent
- Logs: `/var/log/tabssh/tabssh/`
- Source: `src/`
- Docker: `docker/`

## Where to Find Details

- AI behavior: `.claude/rules/ai-rules.md` (PART 0, 1)
- Project structure/license/paths: `.claude/rules/project-rules.md` (PART 2, 3, 4)
- Full spec: `AI.md` (~63k lines) ← **SOURCE OF TRUTH**
- Business logic/scope: `IDEA.md`
- Overrides/activations: `SPEC.md`

## Current Project State

- Last read AI.md: 2026-07-31 (PART 0-6)
- Current task: PART 0-6 compliance pass (rule files, CLAUDE.md loader, structure/paths/config/mode verification)
- Relevant PARTs: 0, 1, 2, 3, 4, 5, 6
- Remaining: `.claude/rules/` files for PART 7+ (binary, backend, api, frontend, features, service, makefile, docker, cicd, testing, optional, config-rules.md PART 12) not yet created — deferred to a follow-up pass once PART 7+ is read
