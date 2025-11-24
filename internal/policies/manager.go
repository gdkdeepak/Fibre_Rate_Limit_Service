package policies

import (
	"sync"

	"fibre_rate_limit_service/internal/limiters"
)

// Manager holds all limiters for policies
type Manager struct {
	mu       sync.RWMutex
	limiters map[string]limiters.Limiter // route name -> limiter
}

// NewManager creates a new policy manager
func NewManager() *Manager {
	return &Manager{
		limiters: make(map[string]limiters.Limiter),
	}
}

// GetLimiter returns the limiter for a given route
func (m *Manager) GetLimiter(route string) limiters.Limiter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.limiters[route]
}

// SetLimiter creates or updates a limiter for a route
func (m *Manager) SetLimiter(route string, l limiters.Limiter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.limiters[route] = l
}

// ListLimiters returns all route names
func (m *Manager) ListLimiters() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.limiters))
	for name := range m.limiters {
		names = append(names, name)
	}
	return names
}
