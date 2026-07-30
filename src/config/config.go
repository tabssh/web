package config

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the full server.yml structure, per AI.md PART 5 Configuration.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Web    WebConfig    `yaml:"web"`
	Cache  *CacheData   `yaml:"_cache,omitempty"`
}

// ServerConfig holds the server section of server.yml.
type ServerConfig struct {
	Port        string            `yaml:"port"`
	FQDN        string            `yaml:"fqdn"`
	Address     string            `yaml:"address"`
	Mode        string            `yaml:"mode"`
	AdminPath   string            `yaml:"admin_path"`
	APIVersion  string            `yaml:"api_version"`
	Healthz     HealthzConfig     `yaml:"healthz"`
	Branding    BrandingConfig    `yaml:"branding"`
	SEO         SEOConfig         `yaml:"seo"`
	User        string            `yaml:"user,omitempty"`
	Group       string            `yaml:"group,omitempty"`
	PIDFile     Bool              `yaml:"pidfile"`
	Daemonize   Bool              `yaml:"daemonize"`
	Admin       AdminConfig       `yaml:"admin"`
	SSL         SSLConfig         `yaml:"ssl"`
	Scheduler   SchedulerConfig   `yaml:"scheduler"`
	RateLimit   RateLimitConfig   `yaml:"rate_limit"`
	Database    DatabaseConfig    `yaml:"database"`
	Maintenance MaintenanceConfig `yaml:"maintenance"`
}

// HealthzConfig controls the optional root /healthz compatibility alias.
type HealthzConfig struct {
	Root HealthzRootConfig `yaml:"root"`
}

// HealthzRootConfig gates mounting /healthz to the same handler as
// /server/healthz (never a redirect).
type HealthzRootConfig struct {
	Enabled Bool `yaml:"enabled"`
}

// BrandingConfig holds branding settings (see PART 16 for full details).
type BrandingConfig struct {
	Title       string `yaml:"title"`
	Tagline     string `yaml:"tagline"`
	Description string `yaml:"description"`
}

// SEOConfig holds SEO settings.
type SEOConfig struct {
	Keywords []string `yaml:"keywords"`
}

// AdminConfig holds admin panel contact settings. Username, password, and
// token live in the database (admins table), never in this config file.
type AdminConfig struct {
	Email string `yaml:"email"`
}

// SSLConfig holds SSL/TLS settings.
type SSLConfig struct {
	Enabled     Bool              `yaml:"enabled"`
	Cert        string            `yaml:"cert"`
	Key         string            `yaml:"key"`
	MinVersion  string            `yaml:"min_version"`
	LetsEncrypt LetsEncryptConfig `yaml:"letsencrypt"`
}

// LetsEncryptConfig holds Let's Encrypt settings.
type LetsEncryptConfig struct {
	Enabled   Bool   `yaml:"enabled"`
	Email     string `yaml:"email"`
	Challenge string `yaml:"challenge"`
	Staging   Bool   `yaml:"staging"`
}

// SchedulerConfig holds the built-in scheduler settings.
type SchedulerConfig struct {
	Enabled Bool                  `yaml:"enabled"`
	Tasks   map[string]TaskConfig `yaml:"tasks"`
}

// TaskConfig holds one scheduled task's settings. Optional fields apply only
// to tasks that use them (retention, renew_before, retry settings).
type TaskConfig struct {
	Enabled     Bool     `yaml:"enabled"`
	Schedule    string   `yaml:"schedule"`
	RetryOnFail Bool     `yaml:"retry_on_fail,omitempty"`
	RetryDelay  Duration `yaml:"retry_delay,omitempty"`
	Retention   int      `yaml:"retention,omitempty"`
	RenewBefore string   `yaml:"renew_before,omitempty"`
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	Enabled     Bool                `yaml:"enabled"`
	Read        RateLimitBucket     `yaml:"read"`
	Write       RateLimitBucket     `yaml:"write"`
	Health      RateLimitBucket     `yaml:"health"`
	GlobalBurst int                 `yaml:"global_burst"`
	Auth        RateLimitAuthConfig `yaml:"auth"`
}

// RateLimitBucket is one per-IP request budget over a window in seconds.
type RateLimitBucket struct {
	Requests int `yaml:"requests"`
	Window   int `yaml:"window"`
}

// RateLimitAuthConfig holds the stricter auth-endpoint limits, applied
// independently of the general limits.
type RateLimitAuthConfig struct {
	Login         RateLimitBucket `yaml:"login"`
	PasswordReset RateLimitBucket `yaml:"password_reset"`
	Registration  RateLimitBucket `yaml:"registration"`
}

// DatabaseConfig holds the database connection settings. In cluster mode
// this is the only non-cache content of server.yml (bootstrap settings).
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
	SSLMode  string `yaml:"sslmode,omitempty"`
}

// MaintenanceConfig holds maintenance mode and self-healing settings.
type MaintenanceConfig struct {
	SelfHealing SelfHealingConfig `yaml:"self_healing"`
	Cleanup     CleanupConfig     `yaml:"cleanup"`
	Notify      NotifyConfig      `yaml:"notify"`
}

// SelfHealingConfig controls the background recovery retry loop.
type SelfHealingConfig struct {
	Enabled       Bool     `yaml:"enabled"`
	RetryInterval Duration `yaml:"retry_interval"`
	MaxAttempts   int      `yaml:"max_attempts"`
}

// CleanupConfig controls auto-cleanup thresholds used during disk-pressure
// self-healing. The backup directory is resolved and cached at startup and
// never re-resolved at cleanup time.
type CleanupConfig struct {
	DiskThreshold    int `yaml:"disk_threshold"`
	LogRetentionDays int `yaml:"log_retention_days"`
	BackupKeepCount  int `yaml:"backup_keep_count"`
}

// NotifyConfig controls maintenance mode enter/exit notifications.
type NotifyConfig struct {
	OnEnter Bool `yaml:"on_enter"`
	OnExit  Bool `yaml:"on_exit"`
}

// WebConfig holds the frontend section of server.yml.
type WebConfig struct {
	UI   UIConfig `yaml:"ui"`
	CORS string   `yaml:"cors"`
}

// UIConfig holds UI settings.
type UIConfig struct {
	Theme string `yaml:"theme"`
}

// CacheData is the _cache section written in cluster mode: the local
// cache/backup of database-held configuration, used read-only when the
// database is unavailable.
type CacheData struct {
	LastSync time.Time              `yaml:"last_sync"`
	Settings map[string]interface{} `yaml:",inline"`
}

// detectFQDN resolves the default FQDN from the host, falling back to
// localhost. Reverse-proxy header resolution happens per-request (PART 8).
func detectFQDN() string {
	if d := os.Getenv("DOMAIN"); d != "" {
		return d
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	if h := os.Getenv("HOSTNAME"); h != "" {
		return h
	}
	return "localhost"
}

// Defaults returns the built-in sane defaults for every setting. The port is
// left empty: a random unused 64xxx port is selected and persisted on first
// run only.
func Defaults() *Config {
	fqdn := detectFQDN()
	adminEmail := "admin@" + fqdn
	return &Config{
		Server: ServerConfig{
			Port:       "",
			FQDN:       fqdn,
			Address:    "[::]",
			Mode:       "production",
			AdminPath:  "admin",
			APIVersion: "v1",
			Healthz:    HealthzConfig{Root: HealthzRootConfig{Enabled: false}},
			Branding: BrandingConfig{
				Title:       "tabssh",
				Tagline:     "",
				Description: "",
			},
			SEO:       SEOConfig{Keywords: []string{}},
			PIDFile:   true,
			Daemonize: false,
			Admin:     AdminConfig{Email: adminEmail},
			SSL: SSLConfig{
				Enabled:    false,
				Cert:       "",
				Key:        "",
				MinVersion: "TLS1.2",
				LetsEncrypt: LetsEncryptConfig{
					Enabled:   false,
					Email:     adminEmail,
					Challenge: "http-01",
					Staging:   false,
				},
			},
			Scheduler: SchedulerConfig{
				Enabled: true,
				Tasks: map[string]TaskConfig{
					"geoip_update":     {Enabled: true, Schedule: "0 3 * * 0", RetryOnFail: true, RetryDelay: Duration(time.Hour)},
					"blocklist_update": {Enabled: true, Schedule: "0 4 * * *", RetryOnFail: true, RetryDelay: Duration(time.Hour)},
					"cve_update":       {Enabled: true, Schedule: "0 5 * * *", RetryOnFail: true, RetryDelay: Duration(time.Hour)},
					"log_rotation":     {Enabled: true, Schedule: "0 0 * * *"},
					"session_cleanup":  {Enabled: true, Schedule: "@hourly"},
					"backup":           {Enabled: true, Schedule: "0 2 * * *", Retention: 4},
					"ssl_renewal":      {Enabled: true, Schedule: "0 3 * * *", RenewBefore: "7d"},
					"health_check":     {Enabled: true, Schedule: "*/5 * * * *"},
					"tor_health":       {Enabled: true, Schedule: "*/10 * * * *"},
				},
			},
			RateLimit: RateLimitConfig{
				Enabled:     true,
				Read:        RateLimitBucket{Requests: 120, Window: 60},
				Write:       RateLimitBucket{Requests: 10, Window: 60},
				Health:      RateLimitBucket{Requests: 120, Window: 60},
				GlobalBurst: 240,
				Auth: RateLimitAuthConfig{
					Login:         RateLimitBucket{Requests: 5, Window: 900},
					PasswordReset: RateLimitBucket{Requests: 3, Window: 3600},
					Registration:  RateLimitBucket{Requests: 5, Window: 3600},
				},
			},
			Database: DatabaseConfig{Driver: "sqlite"},
			Maintenance: MaintenanceConfig{
				SelfHealing: SelfHealingConfig{
					Enabled:       true,
					RetryInterval: Duration(30 * time.Second),
					MaxAttempts:   0,
				},
				Cleanup: CleanupConfig{
					DiskThreshold:    90,
					LogRetentionDays: 7,
					BackupKeepCount:  5,
				},
				Notify: NotifyConfig{OnEnter: true, OnExit: true},
			},
		},
		Web: WebConfig{
			UI:   UIConfig{Theme: "dark"},
			CORS: "*",
		},
	}
}

// applyInitEnv applies the PART 5 init-only environment variables that map
// into server.yml. They are honored ONCE, during first-run config creation,
// then ignored on every later start. Directory init-only variables
// (CONFIG_DIR, DATA_DIR, ...) are resolved by the caller via src/paths.
func applyInitEnv(c *Config) {
	if v := os.Getenv("PORT"); v != "" {
		c.Server.Port = v
	}
	if v := os.Getenv("LISTEN"); v != "" {
		c.Server.Address = v
	}
	if v := os.Getenv("APPLICATION_NAME"); v != "" {
		c.Server.Branding.Title = v
	}
	if v := os.Getenv("APPLICATION_TAGLINE"); v != "" {
		c.Server.Branding.Tagline = v
	}
}

// pickRandomPort selects a random unused TCP port in the 64000-64999 default
// range, per AI.md PART 5 Port Rules.
func pickRandomPort() (int, error) {
	for i := 0; i < 200; i++ {
		p := 64000 + rand.Intn(1000)
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err != nil {
			continue
		}
		l.Close()
		return p, nil
	}
	return 0, errors.New("no unused port available in 64000-64999")
}

// SetAdminPath validates and sets the admin panel path via SafePath.
func (c *Config) SetAdminPath(input string) error {
	safe, err := SafePath(input)
	if err != nil {
		return fmt.Errorf("invalid admin_path: %w", err)
	}
	c.Server.AdminPath = safe
	return nil
}

// Validate checks path-typed configuration values, per the PART 5 rule that
// every configuration path value is normalized and validated.
func (c *Config) Validate() error {
	if _, err := SafePath(c.Server.AdminPath); err != nil {
		return fmt.Errorf("invalid admin_path: %w", err)
	}
	if err := validatePathSegment(c.Server.APIVersion); err != nil {
		return fmt.Errorf("invalid api_version: %w", err)
	}
	return nil
}

// LoadOrCreate loads server.yml from path, creating it with defaults on
// first run. If a legacy server.yaml exists next to a missing server.yml it
// is auto-migrated (renamed) first. The returned bool is true when the file
// was created.
func LoadOrCreate(path string) (*Config, bool, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if migrated, merr := migrateLegacyConfig(path); merr != nil {
			return nil, false, merr
		} else if !migrated {
			cfg, cerr := createDefaultConfig(path)
			if cerr != nil {
				return nil, false, cerr
			}
			return cfg, true, nil
		}
	} else if err != nil {
		return nil, false, err
	}
	cfg, err := Load(path)
	if err != nil {
		return nil, false, err
	}
	return cfg, false, nil
}

// migrateLegacyConfig renames server.yaml to server.yml when the .yml file
// is missing but the legacy .yaml exists. Returns true when migrated.
func migrateLegacyConfig(path string) (bool, error) {
	if !strings.HasSuffix(path, ".yml") {
		return false, nil
	}
	legacy := strings.TrimSuffix(path, ".yml") + ".yaml"
	if _, err := os.Stat(legacy); err != nil {
		return false, nil
	}
	if err := os.Rename(legacy, path); err != nil {
		return false, fmt.Errorf("migrating %s to %s: %w", legacy, path, err)
	}
	return true, nil
}

// createDefaultConfig builds the first-run configuration: defaults plus
// init-only env vars plus a persisted random 64xxx port, written as a fully
// commented server.yml.
func createDefaultConfig(path string) (*Config, error) {
	cfg := Defaults()
	applyInitEnv(cfg)
	if cfg.Server.Port == "" {
		port, err := pickRandomPort()
		if err != nil {
			return nil, err
		}
		cfg.Server.Port = fmt.Sprintf("%d", port)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if err := writeFileAtomic(path, []byte(renderDefaultConfig(cfg))); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Load reads and parses an existing server.yml, layering the file's values
// over the built-in defaults so omitted keys keep sane defaults.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := Defaults()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the configuration back to disk atomically with 0600
// permissions (the file may contain database credentials).
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return writeFileAtomic(path, data)
}

// writeFileAtomic writes data to a temp file in the target directory and
// renames it into place, creating parent directories first.
func writeFileAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

// renderDefaultConfig renders the commented first-run server.yml. The
// content must stay in sync with Defaults(); the round-trip unit test
// enforces this.
func renderDefaultConfig(c *Config) string {
	return fmt.Sprintf(`# =============================================================================
# SERVER CONFIGURATION
# =============================================================================

server:
  # Default: random unused port in 64xxx range, persisted on first run
  port: "%s"
  # Auto-detected from host
  fqdn: "%s"
  # [::] = all interfaces IPv4/IPv6
  address: "%s"
  # production or development
  mode: %s
  # Admin panel path (default: admin) - see PART 17
  admin_path: %s
  # API version prefix (default: v1) - used in /api/{api_version}/ routes
  api_version: %s
  healthz:
    root:
      # Optional root compatibility alias for /server/healthz (never redirect)
      enabled: false

  # Branding & SEO
  branding:
    title: "%s"
    tagline: "%s"
    description: ""
  seo:
    keywords: []

  # PID file
  pidfile: true

  # Daemonize on start (modern service managers prefer foreground)
  daemonize: false

  # Admin Panel
  admin:
    # Username, password, and token are stored in the database, not here
    email: %s

  # SSL/TLS
  ssl:
    enabled: false
    # Manual cert path (optional, empty = auto-detection)
    cert: ""
    # Manual key path (optional, empty = auto-detection)
    key: ""
    # TLS1.2, TLS1.3
    min_version: "TLS1.2"
    letsencrypt:
      enabled: false
      email: %s
      # http-01, tls-alpn-01, dns-01
      challenge: http-01
      # Use staging server for testing
      staging: false

  # Scheduler - manages all background tasks
  scheduler:
    enabled: true
    tasks:
      # Security database updates
      geoip_update:
        enabled: true
        # Weekly: Sunday 3:00 AM
        schedule: "0 3 * * 0"
        retry_on_fail: true
        retry_delay: 1h
      blocklist_update:
        enabled: true
        # Daily: 4:00 AM
        schedule: "0 4 * * *"
        retry_on_fail: true
        retry_delay: 1h
      cve_update:
        enabled: true
        # Daily: 5:00 AM
        schedule: "0 5 * * *"
        retry_on_fail: true
        retry_delay: 1h

      # Maintenance tasks
      log_rotation:
        enabled: true
        # Daily: midnight
        schedule: "0 0 * * *"
      session_cleanup:
        enabled: true
        # Hourly
        schedule: "@hourly"
      backup:
        enabled: true
        # Daily: 2:00 AM
        schedule: "0 2 * * *"
        # Keep max 4 backups (storage management)
        retention: 4

      # SSL certificate management (only when app manages certs)
      ssl_renewal:
        enabled: true
        # Daily: 3:00 AM (after backup at 2:00)
        schedule: "0 3 * * *"
        # Renew 7 days before expiry
        renew_before: 7d

      # Health checks
      health_check:
        enabled: true
        # Every 5 minutes
        schedule: "*/5 * * * *"

      # Tor maintenance
      tor_health:
        enabled: true
        # Every 10 minutes
        schedule: "*/10 * * * *"

  rate_limit:
    enabled: true
    read:
      # per minute per IP
      requests: 120
      window: 60
    write:
      # per minute per IP
      requests: 10
      window: 60
    health:
      # per minute per IP (health/status endpoints)
      requests: 120
      window: 60
    # per minute per IP (absolute ceiling across all endpoint types)
    global_burst: 240
    # Auth endpoints - stricter limits, applied independently
    auth:
      login:
        requests: 5
        # 15 minutes
        window: 900
      password_reset:
        requests: 3
        # 1 hour
        window: 3600
      registration:
        requests: 5
        # 1 hour
        window: 3600

  # Database
  database:
    driver: sqlite

  # Maintenance mode & self-healing
  maintenance:
    self_healing:
      enabled: true
      # Seconds between retry attempts
      retry_interval: 30s
      # 0 = unlimited (keep trying forever)
      max_attempts: 0
    cleanup:
      # Start cleanup when disk > 90%% full
      disk_threshold: 90
      # Delete logs older than 7 days during cleanup
      log_retention_days: 7
      # Keep last 5 backups during cleanup
      backup_keep_count: 5
    notify:
      # Notify when entering maintenance mode
      on_enter: true
      # Notify when exiting maintenance mode
      on_exit: true

# =============================================================================
# FRONTEND CONFIGURATION
# =============================================================================

web:
  ui:
    theme: dark
  cors: "*"
`,
		c.Server.Port,
		c.Server.FQDN,
		c.Server.Address,
		c.Server.Mode,
		c.Server.AdminPath,
		c.Server.APIVersion,
		c.Server.Branding.Title,
		c.Server.Branding.Tagline,
		c.Server.Admin.Email,
		c.Server.SSL.LetsEncrypt.Email,
	)
}
