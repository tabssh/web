package config

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

// fakeDB is a test double for the PART 10 ConfigDatabase contract.
type fakeDB struct {
	settings map[string]interface{}
	pingErr  error
	getErr   error
	setErr   error
}

func newFakeDB() *fakeDB {
	return &fakeDB{settings: map[string]interface{}{}}
}

func (f *fakeDB) Ping(ctx context.Context) error {
	return f.pingErr
}

func (f *fakeDB) GetAllConfig(ctx context.Context) (map[string]interface{}, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	out := map[string]interface{}{}
	for k, v := range f.settings {
		out[k] = v
	}
	return out, nil
}

func (f *fakeDB) SetConfig(ctx context.Context, key string, value interface{}) error {
	if f.setErr != nil {
		return f.setErr
	}
	f.settings[key] = value
	return nil
}

func testManager(t *testing.T, db ConfigDatabase) *Manager {
	t.Helper()
	path := filepath.Join(t.TempDir(), "server.yml")
	cfg := Defaults()
	cfg.Server.Port = "64100"
	return NewManager(cfg, path, db)
}

func TestSourceOfTruth(t *testing.T) {
	single := testManager(t, nil)
	if single.IsCluster() {
		t.Error("nil db must be single-instance mode")
	}
	if got := single.SourceOfTruth(); got != "server.yml" {
		t.Errorf("single-instance source of truth = %q, want server.yml", got)
	}

	cluster := testManager(t, newFakeDB())
	if !cluster.IsCluster() {
		t.Error("non-nil db must be cluster mode")
	}
	if got := cluster.SourceOfTruth(); got != "database" {
		t.Errorf("cluster source of truth = %q, want database", got)
	}
}

func TestSingleInstanceRejectsDatabaseWrites(t *testing.T) {
	m := testManager(t, nil)
	if err := m.SetSetting(context.Background(), "branding.title", "x"); !errors.Is(err, ErrNotCluster) {
		t.Errorf("SetSetting single-instance = %v, want ErrNotCluster", err)
	}
	if err := m.SyncFromDatabase(context.Background()); !errors.Is(err, ErrNotCluster) {
		t.Errorf("SyncFromDatabase single-instance = %v, want ErrNotCluster", err)
	}
	if err := m.SaveLocal(); err != nil {
		t.Errorf("SaveLocal single-instance: %v", err)
	}
}

func TestSetSettingWritesDatabaseThenCache(t *testing.T) {
	db := newFakeDB()
	m := testManager(t, db)
	if err := m.SetSetting(context.Background(), "branding.title", "Clustered"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	// Database first (source of truth)
	if db.settings["branding.title"] != "Clustered" {
		t.Error("value not written to database")
	}
	// Then in-memory
	if v, ok := m.Setting("branding.title"); !ok || v != "Clustered" {
		t.Errorf("in-memory setting = %v, %v", v, ok)
	}
	// Then the server.yml _cache section with a sync timestamp
	back, err := Load(m.path)
	if err != nil {
		t.Fatalf("load cache file: %v", err)
	}
	if back.Cache == nil {
		t.Fatal("_cache section missing from server.yml")
	}
	if back.Cache.LastSync.IsZero() {
		t.Error("last_sync not stamped")
	}
	nested, ok := back.Cache.Settings["branding"].(map[string]interface{})
	if !ok || nested["title"] != "Clustered" {
		t.Errorf("cached settings = %#v", back.Cache.Settings)
	}
}

func TestSensitiveSettingsNeverCached(t *testing.T) {
	db := newFakeDB()
	m := testManager(t, db)
	if err := m.SetSetting(context.Background(), "smtp.password", "hunter2"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	// Database holds it (source of truth) ...
	if db.settings["smtp.password"] != "hunter2" {
		t.Error("password missing from database")
	}
	// ... but the on-disk cache must not
	back, err := Load(m.path)
	if err != nil {
		t.Fatalf("load cache file: %v", err)
	}
	if back.Cache != nil {
		if smtp, ok := back.Cache.Settings["smtp"].(map[string]interface{}); ok {
			if _, leaked := smtp["password"]; leaked {
				t.Error("password cached to disk")
			}
		}
	}
}

func TestDatabaseFailureEntersReadOnly(t *testing.T) {
	db := newFakeDB()
	db.settings["branding.title"] = "FromDB"
	m := testManager(t, db)

	if err := m.SyncFromDatabase(context.Background()); err != nil {
		t.Fatalf("initial sync: %v", err)
	}
	if m.ReadOnly() {
		t.Error("read-only after successful sync")
	}

	// Database goes down: node falls back to read-only cached config
	db.getErr = errors.New("connection refused")
	if err := m.SyncFromDatabase(context.Background()); err == nil {
		t.Fatal("sync succeeded with dead database")
	}
	if !m.ReadOnly() {
		t.Error("not read-only after database failure")
	}
	// Cached settings still served
	if v, ok := m.Setting("branding.title"); !ok || v != "FromDB" {
		t.Errorf("cached setting lost in read-only mode: %v, %v", v, ok)
	}
	// Writes rejected while read-only
	if err := m.SetSetting(context.Background(), "branding.title", "x"); !errors.Is(err, ErrReadOnly) {
		t.Errorf("SetSetting read-only = %v, want ErrReadOnly", err)
	}

	// Database recovers: sync clears read-only mode
	db.getErr = nil
	if err := m.SyncFromDatabase(context.Background()); err != nil {
		t.Fatalf("recovery sync: %v", err)
	}
	if m.ReadOnly() {
		t.Error("still read-only after recovery")
	}
}

func TestManagerPreloadsCachedSettings(t *testing.T) {
	path := filepath.Join(t.TempDir(), "server.yml")
	cfg := Defaults()
	cfg.Server.Port = "64100"
	cfg.Cache = &CacheData{Settings: map[string]interface{}{
		"branding": map[string]interface{}{"title": "Cached"},
	}}
	m := NewManager(cfg, path, newFakeDB())
	if v, ok := m.Setting("branding.title"); !ok || v != "Cached" {
		t.Errorf("preloaded setting = %v, %v", v, ok)
	}
}

func TestNestFlattenRoundTrip(t *testing.T) {
	flat := map[string]interface{}{
		"a.b.c": 1,
		"a.b.d": "x",
		"e":     true,
	}
	back := flattenSettings(nestSettings(flat))
	if len(back) != len(flat) {
		t.Fatalf("round trip size = %d, want %d", len(back), len(flat))
	}
	for k, v := range flat {
		if back[k] != v {
			t.Errorf("round trip %q = %v, want %v", k, back[k], v)
		}
	}
}
