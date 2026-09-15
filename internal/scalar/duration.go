package scalar

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration with text and YAML unmarshaling.
type Duration time.Duration

// Duration returns the underlying time.Duration.
func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

// String returns the formatted duration string.
func (d Duration) String() string {
	return time.Duration(d).String()
}

// ParseDuration parses a human-readable duration string.
func ParseDuration(s string) (Duration, error) {
	var d Duration
	err := d.UnmarshalText([]byte(s))
	return d, err
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (d *Duration) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "0" {
		*d = Duration(0)
		return nil
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}

	*d = Duration(parsed)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	return d.UnmarshalText([]byte(value.Value))
}
