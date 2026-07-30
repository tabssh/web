package config

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestBoolUnmarshalYAML(t *testing.T) {
	tests := []struct {
		input   string
		want    bool
		wantErr bool
	}{
		{"true", true, false},
		{"yes", true, false},
		{"on", true, false},
		{"enabled", true, false},
		{"\"1\"", true, false},
		{"false", false, false},
		{"no", false, false},
		{"off", false, false},
		{"disabled", false, false},
		{"\"0\"", false, false},
		{"maybe", false, true},
	}
	for _, tt := range tests {
		var b Bool
		err := yaml.Unmarshal([]byte(tt.input), &b)
		if (err != nil) != tt.wantErr {
			t.Errorf("Bool unmarshal %q error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && bool(b) != tt.want {
			t.Errorf("Bool unmarshal %q = %v, want %v", tt.input, bool(b), tt.want)
		}
	}
}

func TestBoolMarshalYAML(t *testing.T) {
	out, err := yaml.Marshal(Bool(true))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != "true\n" {
		t.Errorf("Bool marshal = %q, want %q", out, "true\n")
	}
}

func TestDurationUnmarshalYAML(t *testing.T) {
	tests := []struct {
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"30s", 30 * time.Second, false},
		{"1h", time.Hour, false},
		{"1h30m", 90 * time.Minute, false},
		{"45", 45 * time.Second, false},
		{"bogus", 0, true},
	}
	for _, tt := range tests {
		var d Duration
		err := yaml.Unmarshal([]byte(tt.input), &d)
		if (err != nil) != tt.wantErr {
			t.Errorf("Duration unmarshal %q error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && time.Duration(d) != tt.want {
			t.Errorf("Duration unmarshal %q = %v, want %v", tt.input, time.Duration(d), tt.want)
		}
	}
}

func TestDurationRoundTrip(t *testing.T) {
	in := Duration(90 * time.Second)
	out, err := yaml.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var back Duration
	if err := yaml.Unmarshal(out, &back); err != nil {
		t.Fatalf("unmarshal %q: %v", out, err)
	}
	if back != in {
		t.Errorf("round trip = %v, want %v", time.Duration(back), time.Duration(in))
	}
}
