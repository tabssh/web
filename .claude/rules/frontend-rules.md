# Frontend & Admin Panel Rules (PART 16, 17)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

**Non-negotiable: these rules are pulled directly from AI.md. If this file and AI.md ever disagree, AI.md wins — regenerate this file, do not edit it by hand.**

## CRITICAL - NEVER DO
- NEVER implement client-side rendering frameworks (React/Vue), client-side routing, client-side data fetching for initial page load, or business logic in JavaScript — all frontend HTML MUST use Go `html/template`.
- NEVER require JavaScript for a core feature — JS is enhancement only ("The server does the work. The client displays the result").
- NEVER use `template.HTML` on untrusted/user-supplied content except through a sanitizer; `markdownToHTML` must disable raw HTML passthrough, sanitize with an allow-list, escape code fences, and use safe link attributes (`rel="noopener noreferrer nofollow ugc"`).
- NEVER use inline `onclick` or other inline event-handler attributes (CSP blocks them) — use a `data-action` attribute plus a delegated listener in `app.js` instead.
- NEVER use `alert()`, `confirm()`, or `prompt()` — always use the native `<dialog>` element.
- NEVER use deprecated HTML elements or attributes.
- NEVER hardcode colors — use the shared CSS custom-property theme system; never create separate stylesheets per layout (no `admin-dark.css`/`public-light.css`) — public and admin layouts differ in structure only, not theme classes/variables; the theme class goes on `<html>`, never `<body>`.
- NEVER link to the admin path from any public route, and never hint at its existence anywhere in the public site — admin access must be intentional (type the path manually).
- NEVER store PWA session state in JS-held tokens, localStorage, or IndexedDB — session lives only in an HttpOnly+Secure+SameSite cookie; logout must be a plain `POST /logout` form that works with zero JS.
- NEVER include admin pages (`/server/{admin_path}/*`), auth pages, or API endpoints in the generated sitemap.xml.
- NEVER allow remote user-controlled SVG images (logo/favicon/OG image fetch) unless the project explicitly sanitizes and rasterizes them before storage/display; remote images follow the same active-content rules as uploads (SSRF prevention: block private/loopback/link-local IPs and `.local`/`.internal` hostnames, HTTPS-only in production, size limits, MIME allow-list, redirect validation).
- NEVER render site-verification meta tags that are empty, fail validation, contain invalid/potentially-XSS characters, or exceed the provider's max length.
- NEVER allow the footer's `custom_html` to include `<script>`, `<iframe>`, `<object>`, `<embed>`, `<form>`, `<input>`, `<button>`, `<style>`, `<link>`, `<meta>`, `<base>`, event handlers, `javascript:` URLs, or the `style` attribute.
- NEVER use `!important` in CSS except for print styles; never write inline `style` attributes.
- NEVER apply CSRF validation to Bearer/API-token requests, public endpoints, read-only methods, or non-browser callers — CSRF protection applies only to cookie-authenticated browser forms.
- NEVER load analytics or third-party scripts, and never load embedded third-party content (show a placeholder instead), after the user declines cookie consent.
- NEVER use generic placeholder content on the standard pages (`/server/about`, `/server/privacy`, `/server/contact`, `/server/help`, `/server/terms`) — content MUST come from IDEA.md; never use text like "Your application name here", generic feature lists, or example.com API examples.
- NEVER point users at raw `/.well-known/security.txt` from `/server/contact`, and never render `server.contact.admin.email` publicly.
- NEVER create any admin route other than `/server/{admin_path}/{admin_username}/*` (the admin's own account) or `/server/{admin_path}/config/*` (everything server-related) as a direct child of the admin path — e.g. `/server/{admin_path}/settings` is invalid, it must be `/server/{admin_path}/config/settings`.
- NEVER let `--debug` bypass admin authentication.
- NEVER store admin credentials in the config file — they live in `users.db` (admins table), stated explicitly and repeated verbatim in the source.
- NEVER expose a list endpoint or per-admin GET endpoint for other Server Admins, and never let one Server Admin see another's username, email, password, API token, 2FA secret, or session data (Server Admin Privacy) — an admin may see only the total admin count and their own full details.
- NEVER treat a Server Admin as a privileged regular user — Server Admins and regular users are completely separate account types in different database tables; a Server Admin hitting `/users/*` routes is treated as a guest and redirected to `/server/{admin_path}`.
- NEVER let the Primary Admin be deleted except via `--maintenance setup`.
- NEVER relax security rules before initial setup completes — pre-setup public surfaces follow the same public-endpoint, sanitization, rate-limit, TLS, auth, and secret-handling rules as the fully configured app.

## CRITICAL - ALWAYS DO
- ALWAYS build mobile-first responsive CSS, semantic HTML, and ensure every feature works with JavaScript disabled (progressive enhancement layered on top).
- ALWAYS meet WCAG 2.1 AA accessibility, with both required themes (dark default, light) meeting the 4.5:1 minimum contrast ratio and theme switching happening without a page reload.
- ALWAYS embed all templates and static assets into the binary via `//go:embed`; only externally-fetched/downloaded assets (GeoIP DBs, blocklists, CVE DBs, SSL certs) stay external.
- ALWAYS use exactly one JavaScript file (`static/js/app.js`) — no frameworks, no bundlers, no transpilers, no npm/node for the frontend build.
- ALWAYS keep the PWA session cookie HttpOnly + Secure + SameSite, and version the service-worker cache as `{app_name}-cache-v{major}.{minor}.{patch}`.
- ALWAYS use the unified response envelope: success `{"ok":true,"data":{}}`; error `{"ok":false,"error":"CODE","message":"...","details":{}}` with the standard error-code table (BAD_REQUEST, VALIDATION_FAILED, UNAUTHORIZED, TOKEN_EXPIRED, TOKEN_INVALID, 2FA_REQUIRED, 2FA_INVALID, FORBIDDEN, CSRF_FAILED, ACCOUNT_LOCKED, NOT_FOUND, METHOD_NOT_ALLOWED, CONFLICT, RATE_LIMITED, SERVER_ERROR, MAINTENANCE).
- ALWAYS end every non-HTML response with a single trailing newline; use Unix line endings and UTF-8 with no trailing spaces for text responses.
- ALWAYS resolve the CORS allow-list in order: 1) explicit `web.cors` config, 2) `DOMAIN` env entries, 3) reverse-proxy-learned hosts via `X-Forwarded-Host` gated on trusted proxies, 4) default `*` fallback — and only send credentials when the resolved list is explicit, never with `*`.
- ALWAYS style error pages with the site theme system (`error.tmpl` extends `public.tmpl`) — never render a plain/unstyled error page.
- ALWAYS render the cookie consent banner when no `cookie_consent` cookie is present, and skip it entirely once one is; the banner must work fully without JS via plain POST forms to `/server/consent`.
- ALWAYS provide the five standard pages (`/server/about`, `/server/privacy`, `/server/contact`, `/server/help`, `/server/terms`) with real IDEA.md-sourced content, and render `/server/contact` with exactly two informational sections (Security Issues pointing to `/server/security`, and Abuse Reports).
- ALWAYS give every project a full admin panel, completely isolated from the public site, rooted at a configurable `server.admin_path` (default `admin`).
- ALWAYS restrict admin path values to `[a-z0-9-]`, length 2-32, no leading/trailing hyphens, and never allow it to collide with reserved paths (`api`, `health`, `healthz`, `metrics`, `version`, `.well-known`, `about`, `privacy`, `contact`, `help`, `terms`, `docs`, `auth`, `security`, `static`, `assets`).
- ALWAYS test admin routes with real automated tests: verify unauthenticated access is blocked, use the setup token to create an admin, test login, verify session-based access, and verify invalid credentials are rejected.
- ALWAYS require App-works-before-setup behavior: first run auto-creates default `server.yml`, empty `server.db`, auto-detects SMTP, picks a random port in the 64xxx range, generates a one-time setup token, prints a console banner, and starts serving immediately.
- ALWAYS make the one-time setup token single-use, 32 hex chars (128-bit random), displayed exactly once; losing it means the database must be reset.
- ALWAYS apply the same security requirements (password complexity, TOTP, Passkey/WebAuthn, recovery keys, session timeout, API token security, audit logging, rate limiting) to every Server Admin, primary or additional — no exceptions, even for simple apps without regular users.
- ALWAYS keep the admin session cookie (`admin_session`, stored in `server.db` admin_sessions, `admins` table credentials, 30-day default duration) fully separate from the user session cookie (`user_session`, `users.db`, `users` table, 7-day default).

## Key rules summary
- HTML5/CSS-first priority order: HTML5 elements first, CSS second, JavaScript only as a last resort — "JavaScript is the exception, not the rule. Every JS line must be justified."
- Template directory structure separates `layout/`, `partial/`, and `page/`; public pages use `public.tmpl` (clean, top nav, no admin links/hints), admin pages use `admin.tmpl` (sidebar, dashboard-style).
- Static CSS load order: `common.css` → `components.css` → `public.css`/`admin.css`; CSS follows BEM-like naming and mobile-first breakpoints.
- Layout width rules: viewport ≥720px uses 90% width (5% margins); <720px uses 98% width (1% margins); footer always centered and pinned to the bottom via flex layout; mobile nav uses a CSS-only checkbox-hack hamburger, no JS.
- Toasts/modals: modals auto-close on action, trap focus, close on Escape/backdrop click; toasts show at most 5 at once, auto-dismiss by type, pause on hover, positioned top-right.
- PWA requirements: manifest.json with maskable icons (80% safe zone), install-prompt handling (`beforeinstallprompt`, `isInstalledPWA()`), push notifications only requested on user action (never on page load), background sync with retry backoff, offline caching (cache-first static/fonts, network-first HTML, network-only+queue for API).
- Branding/white-labeling is cosmetic only — never changes directory names, system username, log filenames, config paths, binary name, API routes, service names, or container names.
- Admin route hierarchy: only `{admin_username}` (self) and `config` (everything server-related) are valid direct children of `/server/{admin_path}/`; admin path changes reload gracefully (no downtime), port changes require a full restart.
- Admin panel layout: header (logo/title left, search center, status+admin name+logout right), collapsible grouped sidebar (Dashboard, Server, Security, Network, Users, Cluster, Help), dashboard widgets (Status, Uptime, Requests, Errors, System Resources, Quick Actions, Recent Activity, Scheduled Tasks, Alerts).
- Setup Wizard has 6 steps: Create Admin Account, API Token, Server Configuration, Security Settings, Optional Services, Complete.
- Admin invite links are single-use by default (`max_uses=1`, configurable to `0` = unlimited), 7-day default expiry; OIDC/LDAP/SAML admin sync matches on immutable `external_id`+provider, not username/email, with fallback to cached local credentials if the provider is down.
- Server Admin appearance/notification preferences default to dark theme; Security Alerts notifications are always on and cannot be disabled.
- Blocklists subsystem supports P2P, CIDR, DAT, plain-IP, and gzip-compressed sources with auto-format detection, radix-tree lookup, and atomic hot-swap on update; middleware order is allowlist check (bypasses blocklist/ratelimit/geoip, not auth) before blocklist check, before rate limiting, before GeoIP.
- Optional Agent Management (when the project includes an agent component) uses scoped routes (`admin`/`user`/`org`) with owner tokens (`adm_`/`usr_`/`org_`) and matching agent tokens (`adm_agt_`/`usr_agt_`/`org_agt_`), each scope sharing the same endpoint pattern (list/create/get/update/delete/regenerate-token/metrics plus agent-side register/heartbeat/report).

Reference: AI.md PART 16, 17 — Web Frontend; Admin Panel
