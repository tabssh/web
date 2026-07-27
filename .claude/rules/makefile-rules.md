# Makefile Rules (Local Dev Only) (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

**Scope note: the Makefile is for local development only — it is NOT the CI/CD build path. CI/CD workflows are a separate concern.**

## CRITICAL - NEVER DO
- NEVER add more than the six core Makefile targets (dev, local, build, test, release, docker) as the primary target set
- NEVER copy or symlink binaries into `binaries/` — build output must be produced there directly by the build process, not copied/linked after the fact
- NEVER use the Makefile as a substitute for or duplicate of the CI/CD pipeline — it is for local development use
- NEVER tag a release version without the required `v`-prefix convention

## CRITICAL - ALWAYS DO
- ALWAYS implement exactly six core targets: `dev`, `local`, `build`, `test`, `release`, `docker`
- ALWAYS build via Docker using `casjaysdev/go:latest` as the Go toolchain image
- ALWAYS apply the `v`-prefix rule to version tags via the documented shell function
- ALWAYS place final build artifacts directly in `binaries/`

## Key rules summary
- The Makefile's six targets map to distinct concerns: `dev` (local dev run), `local` (local non-docker build/run), `build` (compile via Docker toolchain), `test` (run tests), `release` (produce versioned release artifacts), `docker` (build docker image)
- Version tags always carry a `v` prefix (e.g. `v1.2.3`); a shell function in the Makefile normalizes/enforces this
- Release artifact types and naming are documented in a dedicated table in PART 26 — follow that table exactly for artifact naming and types
- Go builds inside the Makefile use the `casjaysdev/go:latest` toolchain image, consistent with the project-wide Docker build hierarchy rule

Reference: AI.md PART 26 — Makefile
