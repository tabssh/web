package config

import "testing"

func TestParseBool(t *testing.T) {
	tests := []struct {
		input      string
		defaultVal bool
		want       bool
		wantErr    bool
	}{
		{"true", false, true, false},
		{"TRUE", false, true, false},
		{"yes", false, true, false},
		{"Yes", false, true, false},
		{"on", false, true, false},
		{"1", false, true, false},
		{"enable", false, true, false},
		{"enabled", false, true, false},
		{"y", false, true, false},
		{"false", true, false, false},
		{"FALSE", true, false, false},
		{"no", true, false, false},
		{"off", true, false, false},
		{"0", true, false, false},
		{"disable", true, false, false},
		{"disabled", true, false, false},
		{"nein", true, false, false},
		{"n", true, false, false},
		{"  yes  ", false, true, false},
		{"", true, true, false},
		{"", false, false, false},
		{"maybe", false, false, true},
		{"2", false, false, true},
	}
	for _, tt := range tests {
		got, err := ParseBool(tt.input, tt.defaultVal)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseBool(%q, %v) error = %v, wantErr %v", tt.input, tt.defaultVal, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseBool(%q, %v) = %v, want %v", tt.input, tt.defaultVal, got, tt.want)
		}
	}
}

func TestMustParseBool(t *testing.T) {
	if !MustParseBool("yes", false) {
		t.Error("MustParseBool(yes) = false")
	}
	if MustParseBool("no", true) {
		t.Error("MustParseBool(no) = true")
	}
	// Empty input returns the default
	if !MustParseBool("", true) {
		t.Error("MustParseBool(\"\", true) = false, want default true")
	}
	// Invalid input panics (startup-time hard failure)
	defer func() {
		if recover() == nil {
			t.Error("MustParseBool(garbage) did not panic")
		}
	}()
	MustParseBool("garbage", true)
}

func TestIsTruthyIsFalsy(t *testing.T) {
	truthy := []string{"true", "yes", "on", "1", "enable", "enabled", "Y", "TRUE"}
	falsy := []string{"false", "no", "off", "0", "disable", "disabled", "nein", "N"}
	for _, v := range truthy {
		if !IsTruthy(v) {
			t.Errorf("IsTruthy(%q) = false", v)
		}
		if IsFalsy(v) {
			t.Errorf("IsFalsy(%q) = true", v)
		}
	}
	for _, v := range falsy {
		if IsFalsy(v) == false {
			t.Errorf("IsFalsy(%q) = false", v)
		}
		if IsTruthy(v) {
			t.Errorf("IsTruthy(%q) = true", v)
		}
	}
	if IsTruthy("garbage") || IsFalsy("garbage") {
		t.Error("garbage input classified as truthy or falsy")
	}
}
