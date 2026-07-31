# AI Behavior Rules (PART 0, 1)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never modify `AI.md` — it is the read-only spec of truth
- Never guess project variables (project_name, project_org, internal_name, etc.) — infer from git remote/path or ask
- Never skip Session Initialization steps at the start of a task
- Never create forbidden root files: `SUMMARY.md`, `COMPLIANCE.md`, `NOTES.md`, `CHANGELOG.md`, `AUDIT.md`/`REPORT.md`/`ANALYSIS.md`, root-level `CONTRIBUTING.md`/`CODE_OF_CONDUCT.md`/`SECURITY.md`/`PULL_REQUEST_TEMPLATE.md`
- Never create root `Dockerfile`, root `docker-compose.yml`, `*.example.*`/`*.sample.*`, committed `server.yml`/`cli.yml`, or `.env*` files
- Never create forbidden directories at root: `config/`, `data/`, `logs/`, `tmp/`/`temp/`, `test-data/`, `build/`/`dist/`/`out/`, `vendor/`, `node_modules/`, `lib/`/`libs/`, `utils/`/`common/` (note: `src/data/` for embedded static JSON IS allowed)
- Never auto-resolve merge conflicts or bypass a PreToolUse hook block
- Never attribute actions to an invented third-party role ("operator decision", etc.)
- Never leave a flagged-but-unfixed issue only in conversation — log it or fix it
- Never commit — agents edit/report; the coordinator writes COMMIT_MESS and runs `gitcommit`

## CRITICAL - ALWAYS DO

- Read `AI.md` (spec of truth) and `IDEA.md` (business logic, mutable) before acting; `SPEC.md` wins over `AI.md` on conflict
- Follow the 8-step Session Initialization sequence at task start (read spec files, verify variables, check TODO.AI.md, etc.)
- Only create files/directories the spec explicitly mandates — never invent structure
- Check the project's own spec (AI.md/IDEA.md/SPEC.md) before asking the user anything already answered there
- Log any out-of-scope issue found mid-task to `TODO.AI.md` before moving on
- Verify against ground truth before declaring work done
- Keep `.claude/rules/*.md` (14 files, per PART 0's table) and `CLAUDE.md` (loader) in sync with `AI.md`

## Key Reference: `.claude/rules/*.md` File Table

| File | PARTs |
|------|-------|
| ai-rules.md | 0, 1 |
| project-rules.md | 2, 3, 4 |
| config-rules.md | 5, 6, 12 |
| binary-rules.md | 7, 8, 33 |
| backend-rules.md | 9, 10, 11, 32 |
| api-rules.md | 13, 14, 15 |
| frontend-rules.md | 16, 17 |
| features-rules.md | 18-23 |
| service-rules.md | 24, 25 |
| makefile-rules.md | 26 |
| docker-rules.md | 27 |
| cicd-rules.md | 28 |
| testing-rules.md | 29, 30, 31 |
| optional-rules.md | 34-36 |

Each rule file must contain, in order: header, NON-NEGOTIABLE warning, `## CRITICAL - NEVER DO`, `## CRITICAL - ALWAYS DO`, a key-rules summary, and a "see AI.md PART X" reference line.

For complete details, see AI.md PART 0, 1
