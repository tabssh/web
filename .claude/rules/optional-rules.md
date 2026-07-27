# Optional Rules: Multi-User, Organizations, Custom Domains (PART 34-36)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## Activation status for this project (per SPEC.md)
- **PART 34 (Multi-User): ACTIVE**
- **PART 35 (Organizations): ACTIVE**
- **PART 36 (Custom Domains): NOT ACTIVE** — build ZERO scaffolding for it: no tables, routes, config keys, directories, or mentions in generated code/config. Per AI.md PART 0 "Unused Optional Features Must Not Exist in Code."

## CRITICAL - NEVER DO (PART 34/35, active)
- Reveal account existence in auth error messages — always generic, enumeration-safe errors
- Skip respecting `auto_register` and username-collision rules for OIDC/LDAP/SAML-backed regular users on first login
- Use SVG for user-uploaded avatars or external avatar URLs — raster-only unless explicitly sanitized/rasterized at ingest
- Show recovery keys more than once, or fail to force the user to acknowledge saving them
- Let a departing org member retain org shared-vault access after removal — access ends immediately (SPEC.md)
- Let an org admin or instance admin read a user's or org's vault contents, sync blob plaintext, or live session content — none of these have server-side key access (SPEC.md)
- Skip IDOR guards on any per-user or per-org object reference

## CRITICAL - ALWAYS DO (PART 34/35, active)
- Per-user isolated vaults and sync storage; admin-managed accounts
- Org accounts support member roles (member/admin/owner) and org-scoped shared vaults
- First account created at first run becomes the instance admin (zero-config first run, per IDEA.md)
- Recovery keys generated and displayed once when MFA is enabled; user must explicitly confirm they saved them
- Roles enforced exactly per IDEA.md "Roles & permissions": Guest, User, Org member, Org admin, Org owner, Instance admin — each with the exact access boundaries listed there

## PART 36 (Custom Domains) — inactive, reference only
- If ever activated later, declare the activation in SPEC.md first, then read AI.md PART 36 in full before implementing
- Until activated: do not create custom-domain database tables, API routes, config keys, or UI — any such artifact found in the codebase is a bug
- Reference only (do not build): ownership verification uses a `_verify.{domain}` TXT record (never CNAME/A/AAAA) compared with constant-time comparison; SSL via ACME HTTP-01/TLS-ALPN-01/DNS-01; scheduled verification-retry, SSL-renewal, and cleanup tasks; `custom_domains`/`custom_domain_audit` tables in `users.db`

## Key rules summary
- Usernames and org/team slugs share ONE namespace — a name cannot be both; org slug validation reuses the same reserved-names + collision check as usernames
- Org roles: Owner (full control, delete/transfer org, billing), Admin (manage members/settings/tokens, no delete), Member (view/access only) — Danger Zone and Billing are Owner-only
- Organization/team creation mode (`server.orgs.creation.mode`: open/invite/admin_only/disabled) is a server-level policy, distinct from an org's own `visibility: public/private` setting
- Server Admin (PART 17, always required) and Regular User accounts (PART 34, optional) are strictly separate tables/routes; Server Admin can never view full email/passwords/recovery keys/2FA secrets of a user (zero-knowledge boundary) — matches this project's admin-cannot-read-vault-contents rule
- Private org member profile lookups return 404, not 403, whether for users or for the shared-vault/private-content boundary — never leak existence
- Auth token types (Session/Bearer/Reset/Verify/Invite/Tracking/Partial-session) each have distinct construction/storage/expiry; one-shot tokens are `crypto/rand` 32 bytes, SHA-256 stored, single-use, constant-time compared

Reference: AI.md PART 34, PART 35, PART 36 — Multi-User (ACTIVE), Organizations (ACTIVE), Custom Domains (NOT ACTIVE)
