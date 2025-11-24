package limiters

import "sync"

// Manager holds all active limiters.
type Manager struct {
	mu       sync.RWMutex
	limiters map[string]Limiter
}

// NewManager initializes an empty manager.
func NewManager() *Manager {
	return &Manager{
		limiters: make(map[string]Limiter),
	}
}

// AddLimiter adds or replaces a limiter by name.
func (m *Manager) AddLimiter(l Limiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[l.Name()] = l
}

// SetLimiter creates or updates a limiter by name
func (m *Manager) SetLimiter(name string, l Limiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[name] = l
}

// GetLimiter returns a limiter by name.
func (m *Manager) GetLimiter(name string) (Limiter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	l, ok := m.limiters[name]
	return l, ok
}

// ListLimiters returns a list of all limiter names.
func (m *Manager) ListLimiters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.limiters))
	for name := range m.limiters {
		names = append(names, name)
	}
	return names
}

// UpdateLimiter updates a limiter config if it exists.
func (m *Manager) UpdateLimiter(cfg Config) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l, ok := m.limiters[cfg.Name]; ok {
		l.UpdateConfig(cfg)
		return true
	}
	return false
}

// Limiters returns a copy of all limiters (for snapshot purposes)
func (m *Manager) Limiters() map[string]Limiter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string]Limiter, len(m.limiters))
	for k, v := range m.limiters {
		out[k] = v
	}
	return out
}
