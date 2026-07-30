# TODO.AI.md

Tracking list for work remaining after the PART 0-6 bootstrap pass. AI.md is
the source of truth for WHAT must be done — this file only tracks it.

## [x] Finish PART 5 config package (path-security + storage)
Done — src/config: path.go (normalizePath, validatePathSegment, validatePath,
SafePath, SafeFilePath, PathSecurityMiddleware — 3rd in the 10-layer chain),
config.go (server.yml load/save/defaults, 64xxx first-run port, yaml
migration), store.go (ConfigDatabase hook, single-instance vs cluster
source-of-truth, 5-min sync, read-only fallback), maintenance.go
(self-healing retry loop, 503 middleware), types.go; full unit tests.
Remaining 9 middleware layers belong to PART 11/16/17; ConfigDatabase
concrete impl belongs to PART 10.

## [x] Database engine for PART 10 — resolved (user confirmed: in AI.md)
Read: AI.md PART 10
`modernc.org/sqlite` is mandatory for single-instance mode. Cluster mode
supports BOTH PostgreSQL and MySQL as remote engines, exactly as AI.md
PART 10 prescribes (per-engine ALTER/locking/deadlock patterns included).
No separate project decision needed.

## [ ] Binary requirements and CLI (server binary)
Read: AI.md PART 7
Read: AI.md PART 8

## [ ] Error handling & caching
Read: AI.md PART 9

## [ ] Database & cluster mode
Read: AI.md PART 10

## [ ] Security & logging
Read: AI.md PART 11

## [ ] Server configuration (full server.yml schema, admin-editable settings)
Read: AI.md PART 12

## [ ] Health & versioning endpoints
Read: AI.md PART 13

## [ ] API structure (routes, versioning, error envelope)
Read: AI.md PART 14

## [ ] SSL/TLS & Let's Encrypt
Read: AI.md PART 15

## [ ] Web frontend
Read: AI.md PART 16

## [ ] Admin panel
Read: AI.md PART 17

## [ ] Email & notifications
Read: AI.md PART 18

## [ ] Scheduler
Read: AI.md PART 19

## [ ] GeoIP
Read: AI.md PART 20

## [ ] Metrics
Read: AI.md PART 21

## [ ] Backup & restore
Read: AI.md PART 22

## [ ] Update command
Read: AI.md PART 23

## [ ] Privilege escalation & service integration
Read: AI.md PART 24

## [ ] Service support (systemd units, install/uninstall)
Read: AI.md PART 25

## [ ] Makefile (dev/local/build/test targets)
Read: AI.md PART 26

## [ ] Docker (Dockerfile, docker-compose, rootfs, entrypoint)
Read: AI.md PART 27

## [ ] CI/CD workflows
Read: AI.md PART 28

## [ ] Testing & development harness (tests/run_tests.sh, docker.sh, incus.sh)
Read: AI.md PART 29
(Bootstrap test debt cleared — src/paths, src/config ParseBool, and src/mode
unit tests were added with the PART 5 commit.)

## [ ] ReadTheDocs / MkDocs documentation
Read: AI.md PART 30

## [ ] I18n & a11y
Read: AI.md PART 31

## [ ] Tor hidden service support
Read: AI.md PART 32

## [ ] Client & agent (tabssh-cli required, tabssh-agent optional)
Read: AI.md PART 33

## [ ] Multi-user (ACTIVE per SPEC.md)
Read: AI.md PART 34

## [ ] Organizations (ACTIVE per SPEC.md)
Read: AI.md PART 35

## [ ] Cross-app sync contract (usr_/org_ tokens, TABSSH_SYNC_V2, CAS revisions)
Read: SPEC.md
Additive to PART 34/35; do not duplicate — implement alongside those PARTs.

## [x] git init and initial commit
Done — repo initialized on main, remote https://github.com/tabssh/web
auto-created by gitcommit; bootstrap commit 0bc82743149e, scope alignment
commit f0eff4c5f0e4.

Note: PART 36 (Custom Domains) is explicitly NOT ACTIVE per SPEC.md — do not
scaffold any custom-domain code, tables, routes, or config for this project.
