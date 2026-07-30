package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// MaintenanceState is the application health state, per AI.md PART 5 Mode
// States: Starting (checking systems), Normal (full functionality), or
// Maintenance (critical error, read-only plus admin guidance).
type MaintenanceState int

// The three PART 5 mode states.
const (
	StateStarting MaintenanceState = iota
	StateNormal
	StateMaintenance
)

// String returns the lowercase state name.
func (s MaintenanceState) String() string {
	switch s {
	case StateNormal:
		return "normal"
	case StateMaintenance:
		return "maintenance"
	default:
		return "starting"
	}
}

// Critical error reasons. Only database and file-write failures are truly
// critical; everything else self-heals without maintenance mode.
const (
	ReasonDatabaseConnection = "database_connection"
	ReasonDatabaseWrite      = "database_write"
	ReasonFileRead           = "file_read"
	ReasonFileWrite          = "file_write"
	ReasonDiskFull           = "disk_full"
)

// HealFunc attempts to fix the active critical error. A nil return means the
// issue is resolved and verified (health check plus read/write test).
type HealFunc func() error

// NotifyFunc receives maintenance enter/exit events ("enter" or "exit").
// Delivery (email etc.) is wired by the notification subsystem.
type NotifyFunc func(event, reason, message string)

// Maintenance tracks the application state and runs the background
// self-healing retry loop while a critical error is active.
type Maintenance struct {
	mu          sync.RWMutex
	cfg         MaintenanceConfig
	notify      NotifyFunc
	state       MaintenanceState
	reason      string
	message     string
	since       time.Time
	attempts    int
	lastAttempt time.Time
	nextAttempt time.Time
	cancel      chan struct{}
}

// MaintenanceInfo is a snapshot of the maintenance state for /server/healthz
// and the admin panel.
type MaintenanceInfo struct {
	State       MaintenanceState
	Reason      string
	Message     string
	Since       time.Time
	SelfHealing bool
	Attempts    int
	LastAttempt time.Time
	NextAttempt time.Time
}

// NewMaintenance creates a controller in the Starting state. notify may be
// nil when no notification channel is configured.
func NewMaintenance(cfg MaintenanceConfig, notify NotifyFunc) *Maintenance {
	return &Maintenance{cfg: cfg, notify: notify, state: StateStarting}
}

// SetNormal marks startup checks complete and the application healthy.
func (m *Maintenance) SetNormal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state != StateMaintenance {
		m.state = StateNormal
	}
}

// State returns the current mode state.
func (m *Maintenance) State() MaintenanceState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

// Active reports whether maintenance mode is active.
func (m *Maintenance) Active() bool {
	return m.State() == StateMaintenance
}

// Info returns a snapshot for health responses and the admin dashboard.
func (m *Maintenance) Info() MaintenanceInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return MaintenanceInfo{
		State:       m.state,
		Reason:      m.reason,
		Message:     m.message,
		Since:       m.since,
		SelfHealing: bool(m.cfg.SelfHealing.Enabled),
		Attempts:    m.attempts,
		LastAttempt: m.lastAttempt,
		NextAttempt: m.nextAttempt,
	}
}

// Enter switches to maintenance mode for a critical error and starts the
// self-healing retry loop. If maintenance is already active only the reason
// and message are updated; the running loop continues.
func (m *Maintenance) Enter(reason, message string, heal HealFunc) {
	m.mu.Lock()
	if m.state == StateMaintenance {
		m.reason = reason
		m.message = message
		m.mu.Unlock()
		return
	}
	m.state = StateMaintenance
	m.reason = reason
	m.message = message
	m.since = time.Now().UTC()
	m.attempts = 0
	m.lastAttempt = time.Time{}
	interval := m.retryIntervalLocked()
	m.nextAttempt = time.Now().UTC().Add(interval)
	cancel := make(chan struct{})
	m.cancel = cancel
	selfHealing := bool(m.cfg.SelfHealing.Enabled) && heal != nil
	notify := m.notify
	notifyEnter := bool(m.cfg.Notify.OnEnter)
	m.mu.Unlock()

	if notify != nil && notifyEnter {
		notify("enter", reason, message)
	}
	if selfHealing {
		go m.healLoop(heal, cancel, interval)
	}
}

// Exit leaves maintenance mode, stops the retry loop, and re-enables writes.
func (m *Maintenance) Exit() {
	m.mu.Lock()
	if m.state != StateMaintenance {
		m.mu.Unlock()
		return
	}
	m.state = StateNormal
	reason := m.reason
	message := m.message
	m.reason = ""
	m.message = ""
	if m.cancel != nil {
		close(m.cancel)
		m.cancel = nil
	}
	notify := m.notify
	notifyExit := bool(m.cfg.Notify.OnExit)
	m.mu.Unlock()

	if notify != nil && notifyExit {
		notify("exit", reason, message)
	}
}

// retryIntervalLocked returns the configured retry interval, defaulting to
// 30 seconds. Callers must hold m.mu.
func (m *Maintenance) retryIntervalLocked() time.Duration {
	if d := time.Duration(m.cfg.SelfHealing.RetryInterval); d > 0 {
		return d
	}
	return 30 * time.Second
}

// healLoop retries heal every interval until it succeeds, the attempt limit
// is reached (0 = unlimited), or the loop is cancelled. Recovery is
// automatic: a successful heal exits maintenance mode.
func (m *Maintenance) healLoop(heal HealFunc, cancel chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-cancel:
			return
		case <-ticker.C:
			m.mu.Lock()
			m.attempts++
			attempts := m.attempts
			m.lastAttempt = time.Now().UTC()
			m.nextAttempt = m.lastAttempt.Add(interval)
			maxAttempts := m.cfg.SelfHealing.MaxAttempts
			m.mu.Unlock()

			if heal() == nil {
				m.Exit()
				return
			}
			if maxAttempts > 0 && attempts >= maxAttempts {
				// Attempt limit reached: stay in maintenance mode without
				// further retries until an admin intervenes.
				return
			}
		}
	}
}

// maintenanceErrorBody is the canonical PART 14 error envelope with the
// maintenance-specific details PART 5 prescribes.
type maintenanceErrorBody struct {
	OK      bool                    `json:"ok"`
	Error   string                  `json:"error"`
	Message string                  `json:"message"`
	Details maintenanceErrorDetails `json:"details"`
}

// maintenanceErrorDetails carries the machine-readable maintenance context.
type maintenanceErrorDetails struct {
	Reason      string `json:"reason"`
	SelfHealing bool   `json:"self_healing"`
}

// Middleware rejects write operations with 503 and the canonical error body
// while maintenance mode is active. Operational metadata goes in the
// Retry-After and X-Maintenance-* headers, never in ad-hoc body fields.
// Read-only requests pass through so the admin panel stays accessible.
func (m *Maintenance) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Active() || !isWriteMethod(r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		m.mu.RLock()
		reason := m.reason
		message := m.message
		selfHealing := bool(m.cfg.SelfHealing.Enabled)
		interval := m.retryIntervalLocked()
		m.mu.RUnlock()

		body, err := json.Marshal(maintenanceErrorBody{
			OK:      false,
			Error:   "MAINTENANCE",
			Message: "Server is in maintenance mode due to: " + message,
			Details: maintenanceErrorDetails{Reason: reason, SelfHealing: selfHealing},
		})
		if err != nil {
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Retry-After", fmt.Sprintf("%d", int(interval.Seconds())))
		w.Header().Set("X-Maintenance-Mode", "true")
		w.Header().Set("X-Maintenance-Reason", reason)
		w.WriteHeader(http.StatusServiceUnavailable)
		// Every non-HTML response ends with a single trailing newline
		w.Write(append(body, '\n'))
	})
}

// isWriteMethod reports whether the HTTP method is a write operation that
// maintenance mode must reject.
func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
