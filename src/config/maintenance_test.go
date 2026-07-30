package config

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testMaintenanceConfig(interval time.Duration, maxAttempts int) MaintenanceConfig {
	return MaintenanceConfig{
		SelfHealing: SelfHealingConfig{
			Enabled:       true,
			RetryInterval: Duration(interval),
			MaxAttempts:   maxAttempts,
		},
		Notify: NotifyConfig{OnEnter: true, OnExit: true},
	}
}

func TestMaintenanceStateString(t *testing.T) {
	tests := []struct {
		state MaintenanceState
		want  string
	}{
		{StateStarting, "starting"},
		{StateNormal, "normal"},
		{StateMaintenance, "maintenance"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("State %d String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestMaintenanceLifecycle(t *testing.T) {
	m := NewMaintenance(testMaintenanceConfig(time.Hour, 0), nil)
	if m.State() != StateStarting {
		t.Errorf("initial state = %v, want starting", m.State())
	}
	m.SetNormal()
	if m.State() != StateNormal || m.Active() {
		t.Errorf("after SetNormal state = %v", m.State())
	}
	m.Enter(ReasonDatabaseConnection, "db down", nil)
	if !m.Active() {
		t.Error("not active after Enter")
	}
	info := m.Info()
	if info.Reason != ReasonDatabaseConnection || info.Since.IsZero() || !info.SelfHealing {
		t.Errorf("Info = %+v", info)
	}
	// SetNormal must not override active maintenance
	m.SetNormal()
	if !m.Active() {
		t.Error("SetNormal cleared active maintenance mode")
	}
	m.Exit()
	if m.Active() {
		t.Error("still active after Exit")
	}
	if got := m.Info(); got.Reason != "" {
		t.Errorf("reason not cleared on Exit: %q", got.Reason)
	}
}

func TestMaintenanceNotify(t *testing.T) {
	var events []string
	m := NewMaintenance(testMaintenanceConfig(time.Hour, 0), func(event, reason, message string) {
		events = append(events, event+":"+reason)
	})
	m.SetNormal()
	m.Enter(ReasonFileWrite, "disk full", nil)
	m.Exit()
	want := []string{"enter:" + ReasonFileWrite, "exit:" + ReasonFileWrite}
	if len(events) != len(want) || events[0] != want[0] || events[1] != want[1] {
		t.Errorf("notify events = %v, want %v", events, want)
	}
}

func TestSelfHealingAutoExit(t *testing.T) {
	var calls atomic.Int32
	m := NewMaintenance(testMaintenanceConfig(10*time.Millisecond, 0), nil)
	m.SetNormal()
	// Heal fails twice, then succeeds; recovery must be automatic
	m.Enter(ReasonDatabaseConnection, "db down", func() error {
		if calls.Add(1) < 3 {
			return errors.New("still down")
		}
		return nil
	})
	deadline := time.Now().Add(2 * time.Second)
	for m.Active() && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	if m.Active() {
		t.Fatal("maintenance mode did not auto-exit after successful heal")
	}
	if got := calls.Load(); got < 3 {
		t.Errorf("heal attempts = %d, want >= 3", got)
	}
}

func TestSelfHealingMaxAttempts(t *testing.T) {
	var calls atomic.Int32
	m := NewMaintenance(testMaintenanceConfig(5*time.Millisecond, 2), nil)
	m.SetNormal()
	m.Enter(ReasonFileWrite, "disk full", func() error {
		calls.Add(1)
		return errors.New("still broken")
	})
	// Give the loop time to run past the limit if it were unbounded
	time.Sleep(100 * time.Millisecond)
	if got := calls.Load(); got != 2 {
		t.Errorf("heal attempts = %d, want exactly 2 (max_attempts)", got)
	}
	// Attempt limit reached: stays in maintenance until admin intervenes
	if !m.Active() {
		t.Error("maintenance mode exited without a successful heal")
	}
	m.Exit()
}

func TestMaintenanceMiddleware(t *testing.T) {
	m := NewMaintenance(testMaintenanceConfig(30*time.Second, 0), nil)
	m.SetNormal()
	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Normal mode: everything passes
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/things", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("normal-mode POST status = %d, want 200", rec.Code)
	}

	m.Enter(ReasonDatabaseConnection, "db down", nil)

	// Reads still pass (admin panel stays accessible)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/server/admin", nil))
	if rec.Code != http.StatusOK {
		t.Errorf("maintenance GET status = %d, want 200", rec.Code)
	}

	// Writes rejected with 503 plus canonical envelope and headers
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/things", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("maintenance POST status = %d, want 503", rec.Code)
	}
	if rec.Header().Get("Retry-After") != "30" {
		t.Errorf("Retry-After = %q, want 30", rec.Header().Get("Retry-After"))
	}
	if rec.Header().Get("X-Maintenance-Mode") != "true" {
		t.Errorf("X-Maintenance-Mode = %q, want true", rec.Header().Get("X-Maintenance-Mode"))
	}
	if rec.Header().Get("X-Maintenance-Reason") != ReasonDatabaseConnection {
		t.Errorf("X-Maintenance-Reason = %q", rec.Header().Get("X-Maintenance-Reason"))
	}
	body := rec.Body.String()
	if !strings.HasSuffix(body, "\n") {
		t.Error("503 body missing trailing newline")
	}
	var envelope struct {
		OK      bool   `json:"ok"`
		Error   string `json:"error"`
		Message string `json:"message"`
		Details struct {
			Reason      string `json:"reason"`
			SelfHealing bool   `json:"self_healing"`
		} `json:"details"`
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		t.Fatalf("503 body is not valid JSON: %v", err)
	}
	if envelope.OK || envelope.Error != "MAINTENANCE" {
		t.Errorf("envelope = %+v", envelope)
	}
	if envelope.Details.Reason != ReasonDatabaseConnection || !envelope.Details.SelfHealing {
		t.Errorf("details = %+v", envelope.Details)
	}
	// Retry-After and X-Maintenance-* metadata must never be body fields
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		t.Fatalf("raw body parse: %v", err)
	}
	for _, forbidden := range []string{"retry_after", "maintenance_mode", "since"} {
		if _, present := raw[forbidden]; present {
			t.Errorf("operational metadata %q leaked into the body", forbidden)
		}
	}
	m.Exit()
}

func TestIsWriteMethod(t *testing.T) {
	writes := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	reads := []string{http.MethodGet, http.MethodHead, http.MethodOptions}
	for _, method := range writes {
		if !isWriteMethod(method) {
			t.Errorf("isWriteMethod(%q) = false, want true", method)
		}
	}
	for _, method := range reads {
		if isWriteMethod(method) {
			t.Errorf("isWriteMethod(%q) = true, want false", method)
		}
	}
}
