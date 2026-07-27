# CI/CD Workflow Rules (PART 28)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- Never use Makefile targets inside CI/CD workflow files — Makefile is local-dev only (PART 26); CI/CD workflows use explicit shell commands only
- Never omit concurrency/auto-cancel on branch-push workflows (`beta.yml`, `daily.yml`, `docker.yml`, or any other branch-push workflow) — a new push to the same ref MUST cancel the older in-progress run
- Never let a tag-release run cancel a run for a *different* tag — `release.yml` concurrency group scopes to the exact tag ref (`v1.2.4` must not cancel `v1.2.3`)
- Never skip the root `Jenkinsfile` — every project MUST have one alongside the git-hosting-native workflows
- Never include a `-musl` suffix in any binary name produced by CI
- Never use `actions/setup-go` (or any pinned Go version action) — always build inside `casjaysdev/go:latest` unpinned, matching the local toolchain image
- Never pin a third-party Action to a tag or branch — always pin to a full commit SHA
- Never skip the post-push CI status check — a failing or still-running pushed build is never "done"

## CRITICAL - ALWAYS DO
- Every project ships CI/CD workflows appropriate for its git hosting platform (GitHub Actions in `.github/workflows/` for GitHub-hosted projects)
- All workflows set the required standard environment variables (build metadata, version info) consistently across jobs
- Both `amd64` and `arm64` runner/agent labels are available for matrix builds
- Third-party Actions are pinned to a full commit SHA, verified via 3-point SHA verification before merge
- Workflow files are validated locally with `act --list -W {file}` before being staged/committed
- Security-only workflows are created first; `ci.yml`/`release.yml` are created last

## Key rules summary
- Standard workflow set mirrors the project directory tree: `ci.yml` (build/test/lint/coverage/security), `release.yml` (stable releases, tag-scoped concurrency), `beta.yml`, `daily.yml`, `docker.yml` (branch-scoped concurrency with auto-cancel)
- Renovate dependency-update PRs get the same SHA 3-point verification before merge as any other third-party Action pin
- CI/CD workflow content itself (job steps, matrix definitions) is implementation work tracked per-project, not invented ad hoc — always cross-check against the project's own `.github/workflows/` conventions file section

Reference: AI.md PART 28 — CI/CD Workflows
