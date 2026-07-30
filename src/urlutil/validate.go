// FQDN validation per AI.md PART 8 (FQDN Validation Rules) using
// golang.org/x/net/publicsuffix for proper eTLD handling.
package urlutil

import (
	"net"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// devOnlyTLDs are internal/dev-only TLDs blocked in production mode.
var devOnlyTLDs = map[string]bool{
	"localhost": true, "test": true, "example": true, "invalid": true,
	"local": true, "lan": true, "internal": true, "home": true,
	"localdomain": true, "home.arpa": true, "intranet": true,
	"corp": true, "private": true,
}

// IsValidHost validates a host for the given mode. Production requires a
// valid public suffix (ICANN TLD), no IPs, no dev TLDs; development
// additionally allows localhost, dev TLDs, and the dynamic project TLD.
func IsValidHost(host string, devMode bool, projectName string) bool {
	lower := strings.ToLower(strings.TrimSpace(host))

	// Reject empty.
	if lower == "" {
		return false
	}

	// Reject IP addresses always.
	if net.ParseIP(lower) != nil {
		return false
	}

	// Handle localhost.
	if lower == "localhost" {
		return devMode
	}

	// Must contain at least one dot.
	if !strings.Contains(lower, ".") {
		return false
	}

	// Overlay network TLDs - valid but app-managed (not set via DOMAIN).
	if strings.HasSuffix(lower, ".onion") ||
		strings.HasSuffix(lower, ".i2p") ||
		strings.HasSuffix(lower, ".exit") {
		return true
	}

	// Dynamic project-specific TLD (e.g. app.tabssh) is dev-only.
	if projectName != "" && strings.HasSuffix(lower, "."+strings.ToLower(projectName)) {
		return devMode
	}

	// Get the public suffix (TLD or eTLD like co.uk).
	suffix, icann := publicsuffix.PublicSuffix(lower)

	// Dev-only TLDs are valid only in dev mode.
	if devOnlyTLDs[suffix] {
		return devMode
	}

	// In production, require a valid ICANN TLD.
	if !devMode && !icann {
		return false
	}

	// Verify we have at least eTLD+1 (not just the suffix itself).
	etldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(lower)
	if err != nil {
		return false
	}

	// Host must be at least eTLD+1 (e.g. "domain.co.uk" not just "co.uk").
	return len(etldPlusOne) > 0
}

// IsValidSSLHost validates a host for Let's Encrypt certificate requests:
// must be publicly resolvable (production-mode validation) and never a
// .onion address (Tor provides end-to-end encryption already).
func IsValidSSLHost(host string) bool {
	lower := strings.ToLower(host)

	if strings.HasSuffix(lower, ".onion") {
		return false
	}

	// SSL always requires a production-valid host (devMode=false);
	// project TLDs are dev-only, so projectName is irrelevant here.
	return IsValidHost(host, false, "")
}
