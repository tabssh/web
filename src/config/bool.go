// Package config implements TabSSH Web's server.yml configuration handling,
// per AI.md PART 5 (Configuration). This file implements the mandatory
// flexible boolean parser used for every boolean value in the application
// (env vars, config file values, CLI flags, API/form input) — never parse
// booleans any other way.
package config

import (
	"fmt"
	"strings"
)

// truthyValues are case-insensitive strings that parse to true.
var truthyValues = map[string]bool{
	"1": true, "y": true, "t": true,
	"yes": true, "true": true, "on": true, "ok": true,
	"enable": true, "enabled": true,
	"yep": true, "yup": true, "yeah": true,
	"aye": true, "si": true, "oui": true, "da": true, "hai": true,
	"affirmative": true, "accept": true, "allow": true, "grant": true,
	"sure": true, "totally": true,
}

// falsyValues are case-insensitive strings that parse to false.
var falsyValues = map[string]bool{
	"0": true, "n": true, "f": true,
	"no": true, "false": true, "off": true,
	"disable": true, "disabled": true,
	"nope": true, "nah": true, "nay": true,
	"nein": true, "non": true, "niet": true, "iie": true, "lie": true,
	"negative": true, "reject": true, "block": true, "revoke": true,
	"deny": true, "never": true, "noway": true,
}

// ParseBool parses a string into a boolean using the truthy/falsy value
// tables. An empty string returns defaultVal. Any other unrecognized value
// is an error — invalid values are never silently defaulted.
func ParseBool(s string, defaultVal bool) (bool, error) {
	s = strings.TrimSpace(strings.ToLower(s))

	if s == "" {
		return defaultVal, nil
	}

	if truthyValues[s] {
		return true, nil
	}

	if falsyValues[s] {
		return false, nil
	}

	return false, fmt.Errorf("invalid boolean value: %q", s)
}

// MustParseBool parses a string into a boolean, panicking on an invalid
// value. Use only during startup initialization where a bad config value
// should halt the process immediately.
func MustParseBool(s string, defaultVal bool) bool {
	val, err := ParseBool(s, defaultVal)
	if err != nil {
		panic(err)
	}
	return val
}

// IsTruthy reports whether s is a recognized truthy value. Unlike ParseBool
// it never errors: empty, invalid, or falsy input all return false.
func IsTruthy(s string) bool {
	return truthyValues[strings.TrimSpace(strings.ToLower(s))]
}

// IsFalsy reports whether s is a recognized falsy value. Unlike ParseBool
// it never errors: empty, invalid, or truthy input all return false.
func IsFalsy(s string) bool {
	return falsyValues[strings.TrimSpace(strings.ToLower(s))]
}
