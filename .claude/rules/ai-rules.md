# AI Assistant & Critical Rules (PART 0, 1)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Guess or assume - READ THE SPEC or ASK
- Implement without reading relevant PART first
- Modify AI.md PART content (read-only spec)
- Add features not in spec without asking
- Use "I think" or "probably" - KNOW from spec or ASK
- Ask multiple plain-text questions in separate messages - use AskUserQuestion wizard instead
- Use generic placeholder content ("Your app name", "Feature 1")
- Create `/server/about` or `/server/help` with placeholder text
- Leave TODO comments in code - implement fully or don't implement
- Create stub functions or "future" placeholders
- Partial implementations - every feature must be 100% complete
- "I'll come back to this later" - there is no later, do it NOW
- Use bcrypt (use Argon2id) · Use CGO (CGO_ENABLED=0 always) · Create premium tiers
- Use external cron (internal scheduler, PART 19) · Use Makefile in CI/CD (explicit commands only)
- Skip build platforms (all 8: linux/darwin/windows/freebsd × amd64/arm64)
- Client-side rendering (React/Vue) - server-side Go templates only
- Require JavaScript for core features - progressive enhancement only
- Put Dockerfile in root - always `docker/Dockerfile`
- Plain `git commit` / `git push` - use `gitcommit` wrapper only
- Subagent writing `.git/COMMIT_MESS` or calling `gitcommit` - parent instance only
- Read an image larger than 1000x1000 directly - resize first (see AI.md PART 0 "Large Image Handling")
- Use a non-conforming IDEA.md without migration

## CRITICAL - ALWAYS DO
- Read relevant PART before implementing ANY feature
- Search AI.md before asking questions (answer is likely there)
- Follow spec EXACTLY - no "improvements" without approval
- Update IDEA.md when features change
- Keep all docs in sync with code
- When unsure, ASK - never guess or assume
- Use AskUserQuestion wizard - one question at a time, options + custom input
- Source `/server/about` and `/server/help` content from IDEA.md
- Implement features 100% complete - no stubs, no TODOs, no "future"
- ONE thing at a time - finish current task completely before starting another
- Return after cross-references - a "See PART X" jump never replaces the rest of the section you were reading
- NEVER install Go on the local machine - ALL builds/tests use Docker (`casjaysdev/go:latest`) or Incus
- Commit often - small, focused commits; findings-based work = one commit per finding, feature work = one commit per feature

## Key rules summary
- Terminology: `server` = main binary `tabssh` (service) · `client` = CLI binary `tabssh-cli` (required) · `agent` = optional binary `tabssh-agent` · Server Admin = app administrator (NOT OS root) · Regular User = end-user · Cluster Node = another instance of this app · Managed Node = external resource this app controls/monitors
- Key pre-answered decisions: password hash is always Argon2id (never bcrypt, PART 11) · Dockerfile always at `docker/Dockerfile`, never root (PART 27) · CGO never enabled, `CGO_ENABLED=0` always (PART 7) · no premium/paid feature tiers, all features free (PART 1) · no external cron, use built-in scheduler (PART 19) · no client-side rendering frameworks, server-side Go templates only (PART 16)
- Container-only development: local machine has no Go installed - never run `go` directly on host; building/unit tests use Docker `casjaysdev/go:latest`, quick smoke testing `alpine:latest`, full OS/integration testing Incus `debian:latest`; use Makefile targets `make dev`, `make local`, `make build`, `make test`
- Security-first design principles: Never Trust Input · Defense in Depth · Least Privilege · Fail Secure · Secure by Default · Internet-Facing Baseline (assume hostile public network unless IDEA.md says otherwise) · Suggest MFA, never force it
- Input handling: validate type/length/format/range before processing; sanitize for context (HTML-encode, SQL-parameterize); allowlist unknown input; trim whitespace except on passwords (reject, don't trim); never `SELECT *`; parameterized queries only; HTML-escape output + CSP; CSRF tokens on state-changing forms; never shell out with user input; `filepath.Clean()` + reject `..` for paths; consistent timing/vague errors/rate limiting against enumeration
- AI must never reduce friction by disabling/loosening/bypassing authn/authz, TLS/secure cookies, CSRF/CSP/CORS, rate limiting/lockouts, input validation/output sanitization, or least-privilege rules; any intentionally weaker compatibility mode must be explicit, documented in IDEA.md, and labeled a security tradeoff
- Naming/file rules: directories lowercase singular (Go: `handler/`, `model/`) · Go files snake_case · config files lowercase dot-extension · docs UPPERCASE.md · scripts snake_case · binaries `{project_name}-{os}-{arch}`
- Never create: `SUMMARY.md`, `COMPLIANCE.md`, `NOTES.md`, `CHANGELOG.md`, `AUDIT.md`/`REPORT.md`/`ANALYSIS.md`, root `CONTRIBUTING.md`/`CODE_OF_CONDUCT.md`/`SECURITY.md`/`PULL_REQUEST_TEMPLATE.md` (belong in `.github/`), root `Dockerfile`/`docker-compose.yml` (belong in `docker/`), `*.example.*`/`*.sample.*`, committed `server.yml`/`cli.yml`, `.env*`
- Never create directories: root `config/`/`data/`/`logs/`/`tmp/`/`temp/`/`test-data/`, `build/`/`dist/`/`out/` (use gitignored `binaries/`), `vendor/`, `node_modules/`, `lib/`/`libs/`/`utils/`/`common/` (use specific package names); `src/data/` for embedded static JSON is allowed, only root-level `data/` is forbidden
- Allowed root files (exhaustive): `AI.md`, `IDEA.md`, `CLAUDE.md`, `SPEC.md` (optional), `PLAN.md`/`PLAN.AI.md` (optional), `TODO.md`/`TODO.AI.md` (optional), `CLAUDE.local.md` (optional, gitignored), `README.md`, `LICENSE.md`, `Makefile`, `go.mod`, `go.sum`, `release.txt`, `site.txt` (optional), `.gitignore`, `.dockerignore`, `.gitattributes` (optional), `Jenkinsfile`, `mkdocs.yml`, `.readthedocs.yaml`, `.editorconfig` (optional); anything else, ask before creating
- `.claude/settings.json`, `.claude/CLAUDE.md`, `.claude/agents/`, `.claude/hooks/`, `.claude/commands/`, `.claude/plans/`, `.claude/rules/` are committed team config; `.claude/settings.local.json`, `.claude/*.lock`, `.claude/backups/`, `.claude/cache/`, `.claude/history.jsonl` are gitignored
- Tool access: `git status`/`diff`/`log`/`branch`/`add` allowed directly; plain `git commit`/`git push` prohibited (use `gitcommit` wrapper); subagents never write `.git/COMMIT_MESS` or call `gitcommit` - parent instance only

Reference: AI.md PART 0, PART 1 — AI Assistant Rules, Critical Rules
