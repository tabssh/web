# Project SPEC

Project: tabssh (TabSSH Web)
Role: Efficient loader for AI.md

⚠️ **THIS FILE IS AUTO-LOADED EVERY CONVERSATION. FOLLOW IT EXACTLY.** ⚠️

Purpose:
- This file is a short loader for the most important rules
- `AI.md` is the full source of truth
- `SPEC.md` overrides AI.md for this project (optional PART activation, sync contract)
- For complete details, read the referenced PARTs in `AI.md`

## FIRST TURN - MANDATORY

On EVERY new conversation or after "context compacted" message:
1. **READ** `SPEC.md` (project overrides) and the relevant `.claude/rules/*.md` for your current task
2. **NEVER** assume or guess - verify against AI.md before implementing

## Asking Questions

- **Default to continuing work** - do not stop just to ask whether you should continue
- **Never guess** - if the answer cannot be determined from `AI.md`, `IDEA.md`, `SPEC.md`, the codebase, or repo state and the missing information materially changes behavior, scope, or safety, ASK the user
- **Question mark = question** - when user ends with `?`, answer/clarify, don't execute

## Before ANY Code Change

1. Have I read the relevant PART in AI.md? (If no → read it)
2. Does this follow AI.md/SPEC.md EXACTLY? (If unsure → check spec)
3. Am I guessing or do I KNOW from the spec? (If guessing → read spec)
4. Would this pass the compliance checklist? (AI.md FINAL section)

**WHEN IN DOUBT: READ THE SPEC. DO NOT GUESS.**

## Binary Terminology
- **server** = `tabssh` (main binary, runs as service)
- **client** = `tabssh-cli` (REQUIRED companion, CLI/TUI/GUI)
- **agent** = `tabssh-agent` (optional, runs on remote machines)

## Key Placeholders
- `project_name` = tabssh · `project_org` = tabssh
- `internal_name` = tabssh (FROZEN) · `internal_org` = tabssh (FROZEN)
- `app_name` = TabSSH Web

## Project-Specific Activation (see SPEC.md)
- PART 34 (Multi-User): ACTIVE
- PART 35 (Organizations): ACTIVE
- PART 36 (Custom Domains): NOT ACTIVE — build nothing for it

## NEVER Do (Top 19) - VIOLATIONS ARE BUGS
1. Use bcrypt → Use Argon2id
2. Put Dockerfile in root → `docker/Dockerfile`
3. Use CGO → CGO_ENABLED=0 always
4. Hardcode dev values → Detect at runtime
5. Use external cron → Internal scheduler (PART 19)
6. Store passwords plaintext → Argon2id (tokens use SHA-256)
7. Create premium tiers → All features free, no paywalls
8. Use Makefile in CI/CD → Explicit commands only
9. Guess or assume values that a command can produce → Run the command
10. Skip platforms → Build all 8 (linux/darwin/windows/freebsd × amd64/arm64)
11. Client-side rendering (React/Vue) → Server-side Go templates
12. Require JavaScript for core features → Progressive enhancement only
13. Let long strings break mobile → Use word-break CSS
14. Skip validation → Server validates EVERYTHING
15. Implement without reading spec → Read relevant PART first
16. Modify AI.md content → READ-ONLY. Project changes go in IDEA.md/SPEC.md
17. Edit `## Project variables` in IDEA.md without confirming with the user
18. Read an image larger than 1000×1000 directly into context → resize first
19. Use a non-conforming IDEA.md without migration

## ALWAYS Do - NON-NEGOTIABLE
1. Read AI.md before implementing ANY feature
2. Server-side processing (server does the work, client displays)
3. Mobile-first responsive CSS
4. All features work without JavaScript
5. Tor hidden service support (auto-enabled if Tor found)
6. Built-in scheduler, GeoIP, metrics, email, backup, update
7. Full admin panel with ALL settings
8. Client binary for ALL projects
9. Commit often — small, focused commits per gitcommit conventions

## File Locations
- Config: `{config_dir}/server.yml`
- Data: `{data_dir}/`
- Logs: `{log_dir}/`
- Source: `src/`
- Docker: `docker/`

## Where to Find Details
- AI behavior: `.claude/rules/ai-rules.md` (PART 0, 1)
- Project structure: `.claude/rules/project-rules.md` (PART 2, 3, 4)
- Config/modes: `.claude/rules/config-rules.md` (PART 5, 6, 12)
- Binaries: `.claude/rules/binary-rules.md` (PART 7, 8, 33)
- Backend: `.claude/rules/backend-rules.md` (PART 9, 10, 11, 32)
- API: `.claude/rules/api-rules.md` (PART 13, 14, 15)
- Frontend/WebUI: `.claude/rules/frontend-rules.md` (PART 16, 17)
- Features: `.claude/rules/features-rules.md` (PART 18-23)
- Service: `.claude/rules/service-rules.md` (PART 24, 25)
- Makefile: `.claude/rules/makefile-rules.md` (PART 26)
- Docker: `.claude/rules/docker-rules.md` (PART 27)
- CI/CD: `.claude/rules/cicd-rules.md` (PART 28)
- Testing/docs/i18n: `.claude/rules/testing-rules.md` (PART 29, 30, 31)
- Optional (multi-user/orgs/domains): `.claude/rules/optional-rules.md` (PART 34-36)
- Full spec: `AI.md` (~63k lines) ← **SOURCE OF TRUTH**
- Project overrides: `SPEC.md` ← wins over AI.md for this project

## Current Project State
- Last read AI.md: 2026-07-27 (PART 0-6 for bootstrap)
- Current task: bootstrap complete, PART 7+ implementation pending — see TODO.AI.md
- Relevant PARTs: 0-6 executed; 7+ tracked in TODO.AI.md
