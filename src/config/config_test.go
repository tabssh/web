package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestDefaults(t *testing.T) {
	c := Defaults()
	if c.Server.Mode != "production" {
		t.Errorf("default mode = %q, want %q", c.Server.Mode, "production")
	}
	if c.Server.AdminPath != "admin" {
		t.Errorf("default admin_path = %q, want %q", c.Server.AdminPath, "admin")
	}
	if c.Server.APIVersion != "v1" {
		t.Errorf("default api_version = %q, want %q", c.Server.APIVersion, "v1")
	}
	if c.Server.Port != "" {
		t.Errorf("default port = %q, want empty (random 64xxx chosen on first run only)", c.Server.Port)
	}
	if c.Server.Database.Driver != "sqlite" {
		t.Errorf("default database driver = %q, want %q", c.Server.Database.Driver, "sqlite")
	}
	if !bool(c.Server.Maintenance.SelfHealing.Enabled) {
		t.Error("self-healing must be enabled by default")
	}
	if c.Web.UI.Theme != "dark" {
		t.Errorf("default theme = %q, want %q", c.Web.UI.Theme, "dark")
	}
	if err := c.Validate(); err != nil {
		t.Errorf("Defaults().Validate() = %v", err)
	}
}

func TestPickRandomPort(t *testing.T) {
	p, err := pickRandomPort()
	if err != nil {
		t.Fatalf("pickRandomPort: %v", err)
	}
	if p < 64000 || p > 64999 {
		t.Errorf("port %d outside 64000-64999", p)
	}
}

func TestSetAdminPath(t *testing.T) {
	c := Defaults()
	if err := c.SetAdminPath("/my-admin/"); err != nil {
		t.Fatalf("SetAdminPath valid: %v", err)
	}
	if c.Server.AdminPath != "my-admin" {
		t.Errorf("admin_path = %q, want %q", c.Server.AdminPath, "my-admin")
	}
	if err := c.SetAdminPath("../etc"); !errors.Is(err, ErrPathTraversal) {
		t.Errorf("SetAdminPath traversal error = %v, want ErrPathTraversal", err)
	}
	if c.Server.AdminPath != "my-admin" {
		t.Errorf("admin_path changed on invalid input: %q", c.Server.AdminPath)
	}
}

func TestValidateRejectsBadPaths(t *testing.T) {
	c := Defaults()
	c.Server.AdminPath = "../evil"
	if err := c.Validate(); err == nil {
		t.Error("Validate accepted traversal admin_path")
	}
	c = Defaults()
	c.Server.APIVersion = "V1!"
	if err := c.Validate(); err == nil {
		t.Error("Validate accepted invalid api_version")
	}
}

func TestLoadOrCreateFirstRun(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tabssh", "server.yml")

	cfg, created, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("LoadOrCreate first run: %v", err)
	}
	if !created {
		t.Fatal("created = false on first run, want true")
	}
	port, err := strconv.Atoi(cfg.Server.Port)
	if err != nil || port < 64000 || port > 64999 {
		t.Errorf("first-run port = %q, want random 64000-64999", cfg.Server.Port)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat created config: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("config permissions = %o, want 0600", info.Mode().Perm())
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read created config: %v", err)
	}
	if !strings.HasSuffix(string(data), "\n") || strings.HasSuffix(string(data), "\n\n") {
		t.Error("created config must end with exactly one trailing newline")
	}

	// Second start: file exists, port must persist, created = false
	cfg2, created2, err := LoadOrCreate(path)
	if err != nil {
		t.Fatalf("LoadOrCreate second run: %v", err)
	}
	if created2 {
		t.Error("created = true on second run, want false")
	}
	if cfg2.Server.Port != cfg.Server.Port {
		t.Errorf("port changed across restarts: %q -> %q", cfg.Server.Port, cfg2.Server.Port)
	}
}

func TestLoadOrCreateMigratesLegacyYAML(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, "server.yaml")
	target := filepath.Join(dir, "server.yml")
	content := "server:\n  port: \"64123\"\n  fqdn: legacy.example\n"
	if err := os.WriteFile(legacy, []byte(content), 0o600); err != nil {
		t.Fatalf("write legacy: %v", err)
	}

	cfg, created, err := LoadOrCreate(target)
	if err != nil {
		t.Fatalf("LoadOrCreate with legacy file: %v", err)
	}
	if created {
		t.Error("created = true for migrated config, want false")
	}
	if cfg.Server.Port != "64123" {
		t.Errorf("migrated port = %q, want %q", cfg.Server.Port, "64123")
	}
	if cfg.Server.FQDN != "legacy.example" {
		t.Errorf("migrated fqdn = %q, want %q", cfg.Server.FQDN, "legacy.example")
	}
	if _, err := os.Stat(legacy); !errors.Is(err, os.ErrNotExist) {
		t.Error("legacy server.yaml still exists after migration")
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("server.yml missing after migration: %v", err)
	}
}

func TestLoadLayersOverDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.yml")
	content := "server:\n  port: \"64500\"\n  branding:\n    title: Custom\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.Port != "64500" {
		t.Errorf("port = %q, want %q", cfg.Server.Port, "64500")
	}
	if cfg.Server.Branding.Title != "Custom" {
		t.Errorf("title = %q, want %q", cfg.Server.Branding.Title, "Custom")
	}
	// Omitted keys keep defaults
	if cfg.Server.AdminPath != "admin" {
		t.Errorf("omitted admin_path = %q, want default %q", cfg.Server.AdminPath, "admin")
	}
	if !bool(cfg.Server.RateLimit.Enabled) {
		t.Error("omitted rate_limit.enabled lost its default true")
	}
}

func TestSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.yml")
	cfg := Defaults()
	cfg.Server.Port = "64777"
	cfg.Server.Branding.Title = "RoundTrip"
	if err := cfg.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	back, err := Load(path)
	if err != nil {
		t.Fatalf("Load after Save: %v", err)
	}
	if back.Server.Port != "64777" || back.Server.Branding.Title != "RoundTrip" {
		t.Errorf("round trip lost values: port=%q title=%q", back.Server.Port, back.Server.Branding.Title)
	}
}

func TestRenderDefaultConfigParsesAndMatchesDefaults(t *testing.T) {
	cfg := Defaults()
	cfg.Server.Port = "64888"
	dir := t.TempDir()
	path := filepath.Join(dir, "server.yml")
	if err := writeFileAtomic(path, []byte(renderDefaultConfig(cfg))); err != nil {
		t.Fatalf("write rendered config: %v", err)
	}
	back, err := Load(path)
	if err != nil {
		t.Fatalf("rendered default config does not parse: %v", err)
	}
	if back.Server.Port != cfg.Server.Port {
		t.Errorf("rendered port = %q, want %q", back.Server.Port, cfg.Server.Port)
	}
	if back.Server.AdminPath != cfg.Server.AdminPath {
		t.Errorf("rendered admin_path = %q, want %q", back.Server.AdminPath, cfg.Server.AdminPath)
	}
	if back.Server.Mode != cfg.Server.Mode {
		t.Errorf("rendered mode = %q, want %q", back.Server.Mode, cfg.Server.Mode)
	}
	if bool(back.Server.Scheduler.Enabled) != bool(cfg.Server.Scheduler.Enabled) {
		t.Error("rendered scheduler.enabled drifted from Defaults()")
	}
	if back.Server.Maintenance.SelfHealing.MaxAttempts != cfg.Server.Maintenance.SelfHealing.MaxAttempts {
		t.Error("rendered self_healing.max_attempts drifted from Defaults()")
	}
	// Inline comments are forbidden in YAML: every comment line starts with #
	for _, line := range strings.Split(renderDefaultConfig(cfg), "\n") {
		if i := strings.Index(line, "#"); i > 0 && strings.TrimSpace(line[:i]) != "" {
			t.Errorf("inline YAML comment found: %q", line)
		}
	}
}

func TestApplyInitEnv(t *testing.T) {
	t.Setenv("PORT", "64321")
	t.Setenv("LISTEN", "127.0.0.1")
	t.Setenv("APPLICATION_NAME", "EnvName")
	t.Setenv("APPLICATION_TAGLINE", "EnvTag")
	c := Defaults()
	applyInitEnv(c)
	if c.Server.Port != "64321" {
		t.Errorf("PORT not applied: %q", c.Server.Port)
	}
	if c.Server.Address != "127.0.0.1" {
		t.Errorf("LISTEN not applied: %q", c.Server.Address)
	}
	if c.Server.Branding.Title != "EnvName" {
		t.Errorf("APPLICATION_NAME not applied: %q", c.Server.Branding.Title)
	}
	if c.Server.Branding.Tagline != "EnvTag" {
		t.Errorf("APPLICATION_TAGLINE not applied: %q", c.Server.Branding.Tagline)
	}
}

func TestWriteFileAtomicCreatesParents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a", "b", "c.yml")
	if err := writeFileAtomic(path, []byte("x: 1\n")); err != nil {
		t.Fatalf("writeFileAtomic: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if string(data) != "x: 1\n" {
		t.Errorf("content = %q", data)
	}
	if _, err := os.Stat(path + ".tmp"); !errors.Is(err, os.ErrNotExist) {
		t.Error("temp file left behind")
	}
}
