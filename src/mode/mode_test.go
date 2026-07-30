package mode

import "testing"

// resetMode restores the package defaults after a test mutates global state.
func resetMode(t *testing.T) {
	t.Cleanup(func() {
		SetAppMode("production")
		SetDebugEnabled(false)
	})
}

func TestAppModeString(t *testing.T) {
	if got := Production.String(); got != "production" {
		t.Errorf("Production.String() = %q", got)
	}
	if got := Development.String(); got != "development" {
		t.Errorf("Development.String() = %q", got)
	}
}

func TestSetAppMode(t *testing.T) {
	resetMode(t)
	tests := []struct {
		input     string
		wantMode  AppMode
		wantDebug bool
	}{
		{"production", Production, false},
		{"prod", Production, false},
		{"dev", Development, false},
		{"devel", Development, false},
		{"development", Development, false},
		{"DEVELOPMENT", Development, false},
		{"debug", Development, true},
		{"garbage", Production, false},
		{"", Production, false},
	}
	for _, tt := range tests {
		SetDebugEnabled(false)
		SetAppMode(tt.input)
		if GetCurrentAppMode() != tt.wantMode {
			t.Errorf("SetAppMode(%q) mode = %v, want %v", tt.input, GetCurrentAppMode(), tt.wantMode)
		}
		if IsDebugEnabled() != tt.wantDebug {
			t.Errorf("SetAppMode(%q) debug = %v, want %v", tt.input, IsDebugEnabled(), tt.wantDebug)
		}
	}
}

func TestModePredicates(t *testing.T) {
	resetMode(t)
	SetAppMode("production")
	if !IsAppModeProd() || IsAppModeDev() {
		t.Error("production predicates wrong")
	}
	SetAppMode("development")
	if !IsAppModeDev() || IsAppModeProd() {
		t.Error("development predicates wrong")
	}
}

func TestGetAppModeString(t *testing.T) {
	resetMode(t)
	SetAppMode("production")
	SetDebugEnabled(false)
	if got := GetAppModeString(); got != "production" {
		t.Errorf("GetAppModeString() = %q", got)
	}
	SetDebugEnabled(true)
	if got := GetAppModeString(); got != "production [debugging]" {
		t.Errorf("GetAppModeString() with debug = %q", got)
	}
}

func TestFromEnv(t *testing.T) {
	resetMode(t)
	tests := []struct {
		name      string
		mode      string
		debug     string
		debugSet  bool
		wantMode  AppMode
		wantDebug bool
	}{
		{"defaults", "", "", false, Production, false},
		{"mode dev", "development", "", false, Development, false},
		{"debug env", "production", "true", true, Production, true},
		{"mode debug alias", "debug", "", false, Development, true},
		// Explicit DEBUG env always wins over the MODE=debug alias
		{"alias overridden", "debug", "false", true, Development, false},
		{"debug truthy word", "production", "yes", true, Production, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAppMode("production")
			SetDebugEnabled(false)
			t.Setenv("MODE", tt.mode)
			if tt.debugSet {
				t.Setenv("DEBUG", tt.debug)
			} else {
				t.Setenv("DEBUG", "")
			}
			FromEnv()
			if GetCurrentAppMode() != tt.wantMode {
				t.Errorf("mode = %v, want %v", GetCurrentAppMode(), tt.wantMode)
			}
			if IsDebugEnabled() != tt.wantDebug {
				t.Errorf("debug = %v, want %v", IsDebugEnabled(), tt.wantDebug)
			}
		})
	}
}
