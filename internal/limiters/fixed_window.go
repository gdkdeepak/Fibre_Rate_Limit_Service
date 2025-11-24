package limiters

import (
	"time"

	"fibre_rate_limit_service/internal/storage"
)

// FixedWindowConfig defines configuration for a fixed window limiter
type FixedWindowConfig struct {
	Name   string
	Limit  int           // max requests per window
	Window time.Duration // window size
	TTL    time.Duration // optional TTL for storage
}

// FixedWindowLimiter implements Limiter interface
type FixedWindowLimiter struct {
	cfg   FixedWindowConfig
	store *storage.ShardedMap
}

// NewFixedWindowLimiter creates a new limiter
func NewFixedWindowLimiter(cfg FixedWindowConfig, store *storage.ShardedMap) Limiter {
	return &FixedWindowLimiter{
		cfg:   cfg,
		store: store,
	}
}

// Name returns the limiter name
func (fw *FixedWindowLimiter) Name() string {
	return fw.cfg.Name
}

// Check implements rate limiting logic
func (fw *FixedWindowLimiter) Check(key string) Result {
	now := time.Now()

	// Load current counter
	raw, _ := fw.store.Get(key)
	count := 0
	windowStart := now

	if raw != nil {
		entry := raw.(windowState)
		count = entry.Count
		windowStart = entry.Start
	}

	// Reset window if expired
	if now.Sub(windowStart) >= fw.cfg.Window {
		count = 0
		windowStart = now
	}

	allowed := false
	if count < fw.cfg.Limit {
		count++
		allowed = true
	}

	// Save updated state
	fw.store.Set(key, windowState{
		Start: windowStart,
		Count: count,
	}, fw.cfg.TTL)

	return Result{
		Allowed:   allowed,
		Remaining: fw.cfg.Limit - count,
		ResetAt:   windowStart.Add(fw.cfg.Window),
		Reason:    "",
	}
}

// UpdateConfig allows updating limiter settings
func (fw *FixedWindowLimiter) UpdateConfig(cfg Config) {
	fw.cfg.Limit = cfg.Limit
	fw.cfg.Window = cfg.Window
}

// windowState stores per-key counter and window start time
type windowState struct {
	Start time.Time
	Count int
}
