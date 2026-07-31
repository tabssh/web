package config

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
)

// DefaultSyncInterval is the periodic database-to-server.yml drift sync
// interval in cluster mode, per AI.md PART 5 Config Sync.
const DefaultSyncInterval = 5 * time.Minute

// Manager mode errors.
var (
	// ErrReadOnly is returned for config writes while the database is
	// unavailable and the node runs read-only from the cached config.
	ErrReadOnly = errors.New("configuration is read-only: database unavailable")
	// ErrNotCluster is returned when a database-backed setting write is
	// attempted in single-instance mode, where server.yml is the source of
	// truth and is edited via Config plus SaveLocal.
	ErrNotCluster = errors.New("not in cluster mode: edit server.yml via Config")
)

// ConfigDatabase is the remote configuration store contract that the PART 10
// database layer implements. In cluster mode the database is the source of
// truth; server.yml is a local cache and backup.
type ConfigDatabase interface {
	// Ping verifies database connectivity.
	Ping(ctx context.Context) error
	// GetAllConfig returns every configuration key-value pair.
	GetAllConfig(ctx context.Context) (map[string]interface{}, error)
	// SetConfig writes one configuration key-value pair.
	SetConfig(ctx context.Context, key string, value interface{}) error
}

// Manager implements the PART 5 configuration source-of-truth logic:
// single-instance mode (nil database) keeps server.yml authoritative;
// cluster mode keeps the database authoritative with server.yml as a synced
// cache/backup and falls back to read-only cached config when the database
// is unavailable.
type Manager struct {
	mu       sync.RWMutex
	cfg      *Config
	path     string
	db       ConfigDatabase
	readOnly bool
	settings map[string]interface{}
}

// NewManager creates a Manager. Pass a nil db for single-instance mode.
func NewManager(cfg *Config, path string, db ConfigDatabase) *Manager {
	m := &Manager{
		cfg:      cfg,
		path:     path,
		db:       db,
		settings: map[string]interface{}{},
	}
	// Preload the cached settings so a node that starts while the database
	// is down still has its last-synced configuration available.
	if cfg.Cache != nil {
		m.settings = flattenSettings(cfg.Cache.Settings)
	}
	return m
}

// IsCluster reports whether a remote configuration database is attached.
func (m *Manager) IsCluster() bool {
	return m.db != nil
}

// SourceOfTruth names the authoritative configuration store for this mode.
func (m *Manager) SourceOfTruth() string {
	if m.IsCluster() {
		return "database"
	}
	return "server.yml"
}

// ReadOnly reports whether the node is running read-only from cached config
// because the database is unavailable.
func (m *Manager) ReadOnly() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.readOnly
}

// Setting returns a database-synced setting by dotted key (cluster mode).
func (m *Manager) Setting(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.settings[key]
	return v, ok
}

// Config returns the underlying typed configuration.
func (m *Manager) Config() *Config {
	return m.cfg
}

// SaveLocal persists the typed configuration to server.yml. In
// single-instance mode this is how the admin panel writes settings.
func (m *Manager) SaveLocal() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cfg.Save(m.path)
}

// SyncFromDatabase pulls all configuration from the database into memory and
// the server.yml _cache section. On failure the node enters read-only mode
// using the cached config; a later successful sync clears read-only mode.
func (m *Manager) SyncFromDatabase(ctx context.Context) error {
	if !m.IsCluster() {
		return ErrNotCluster
	}
	settings, err := m.db.GetAllConfig(ctx)
	if err != nil {
		m.mu.Lock()
		m.readOnly = true
		m.mu.Unlock()
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.readOnly = false
	m.settings = settings
	return m.writeCacheLocked()
}

// SetSetting writes one configuration change: database first (source of
// truth), then the local server.yml cache, then the sync timestamp — the
// PART 5 onConfigChange flow.
func (m *Manager) SetSetting(ctx context.Context, key string, value interface{}) error {
	if !m.IsCluster() {
		return ErrNotCluster
	}
	if m.ReadOnly() {
		return ErrReadOnly
	}
	if err := m.db.SetConfig(ctx, key, value); err != nil {
		m.mu.Lock()
		m.readOnly = true
		m.mu.Unlock()
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.settings[key] = value
	return m.writeCacheLocked()
}

// StartPeriodicSync runs SyncFromDatabase every interval (default 5 minutes
// when interval <= 0) until ctx is done, catching any drift and recovering
// from read-only mode automatically. Cluster mode only; call in a goroutine
// context of the caller's choosing — this function blocks.
func (m *Manager) StartPeriodicSync(ctx context.Context, interval time.Duration) {
	if !m.IsCluster() {
		return
	}
	if interval <= 0 {
		interval = DefaultSyncInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Errors keep the node read-only; the next tick retries.
			m.SyncFromDatabase(ctx)
		}
	}
}

// writeCacheLocked rewrites the _cache section of server.yml from the
// in-memory settings and stamps last_sync. Secrets (passwords, tokens,
// secret keys) are never cached to disk. Callers must hold m.mu.
func (m *Manager) writeCacheLocked() error {
	m.cfg.Cache = &CacheData{
		LastSync: time.Now().UTC().Truncate(time.Second),
		Settings: nestSettings(m.settings),
	}
	return m.cfg.Save(m.path)
}

// sensitiveSetting reports whether a dotted config key must be excluded from
// the on-disk cache (credentials NOT cached, per PART 5). Matching is on
// substrings of the last segment so compound keys such as "secret_key",
// "api_key", and "private_key" are caught alongside the bare names.
func sensitiveSetting(key string) bool {
	parts := strings.Split(key, ".")
	last := strings.ToLower(parts[len(parts)-1])
	for _, marker := range []string{"password", "passwd", "secret", "token", "apikey", "api_key", "key", "credential", "passphrase"} {
		if strings.Contains(last, marker) {
			return true
		}
	}
	return false
}

// nestSettings converts flat dotted keys ("branding.title") into the nested
// map structure the _cache section uses, dropping sensitive keys.
func nestSettings(flat map[string]interface{}) map[string]interface{} {
	nested := map[string]interface{}{}
	for key, value := range flat {
		if sensitiveSetting(key) {
			continue
		}
		parts := strings.Split(key, ".")
		node := nested
		for _, part := range parts[:len(parts)-1] {
			child, ok := node[part].(map[string]interface{})
			if !ok {
				child = map[string]interface{}{}
				node[part] = child
			}
			node = child
		}
		node[parts[len(parts)-1]] = value
	}
	return nested
}

// flattenSettings converts the nested _cache structure back into flat
// dotted keys for in-memory lookup.
func flattenSettings(nested map[string]interface{}) map[string]interface{} {
	flat := map[string]interface{}{}
	flattenInto(flat, "", nested)
	return flat
}

// flattenInto recursively flattens nested maps under the given key prefix.
func flattenInto(flat map[string]interface{}, prefix string, node map[string]interface{}) {
	for key, value := range node {
		full := key
		if prefix != "" {
			full = prefix + "." + key
		}
		if child, ok := value.(map[string]interface{}); ok {
			flattenInto(flat, full, child)
			continue
		}
		flat[full] = value
	}
}
