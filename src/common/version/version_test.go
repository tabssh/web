package version

import (
	"strings"
	"testing"
)

func TestUserAgent(t *testing.T) {
	tests := []struct {
		name    string
		project string
		version string
		want    string
	}{
		{"defaults", "tabssh", "dev", "tabssh/dev"},
		{"release", "tabssh", "1.0.0", "tabssh/1.0.0"},
		{"beta", "tabssh", "20260101120000-beta", "tabssh/20260101120000-beta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origP, origV := ProjectName, Version
			defer func() { ProjectName, Version = origP, origV }()
			ProjectName, Version = tt.project, tt.version
			if got := UserAgent(); got != tt.want {
				t.Errorf("UserAgent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBinaryName(t *testing.T) {
	got := BinaryName()
	if got == "" {
		t.Fatal("BinaryName() returned empty string")
	}
	if strings.ContainsRune(got, '/') {
		t.Errorf("BinaryName() = %q, must not contain path separators", got)
	}
}

func TestGoVersion(t *testing.T) {
	if !strings.HasPrefix(GoVersion(), "go") {
		t.Errorf("GoVersion() = %q, want go-prefixed runtime version", GoVersion())
	}
}
