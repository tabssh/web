package config

import (
	"fmt"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Bool is a YAML boolean that accepts the full truthy/falsy word set via
// ParseBool, per AI.md PART 5 Boolean Handling. Config file fields must use
// this type instead of plain bool so "yes", "on", "enable", etc. all work.
type Bool bool

// UnmarshalYAML parses the scalar through ParseBool. An empty value keeps the
// pre-populated default; an unrecognized value is an error, never silently
// defaulted.
func (b *Bool) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("boolean value must be a scalar")
	}
	v, err := ParseBool(n.Value, bool(*b))
	if err != nil {
		return err
	}
	*b = Bool(v)
	return nil
}

// MarshalYAML always emits a canonical true/false value.
func (b Bool) MarshalYAML() (interface{}, error) {
	return bool(b), nil
}

// Duration is a YAML duration accepting Go duration strings ("30s", "1h") or
// a bare integer meaning seconds. An empty value keeps the field's default.
type Duration time.Duration

// UnmarshalYAML parses a duration scalar; bare integers are seconds.
func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("duration value must be a scalar")
	}
	s := n.Value
	if s == "" {
		return nil
	}
	if v, err := strconv.Atoi(s); err == nil {
		*d = Duration(time.Duration(v) * time.Second)
		return nil
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration value: %q", s)
	}
	*d = Duration(v)
	return nil
}

// MarshalYAML emits the canonical Go duration string form.
func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}
