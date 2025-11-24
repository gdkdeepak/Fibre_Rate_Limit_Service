package limiters

import "time"

// Result stores the outcome of a limiter check.
type Result struct {
	Allowed   bool
	Remaining int
	ResetAt   time.Time
	Reason    string
}

// Config defines settings for any limiter.
type Config struct {
	Name    string
	Limit   int
	Window  time.Duration
	Burst   int
	Enabled bool
}

// Limiter is the interface all rate limiters must implement.
type Limiter interface {
	Name() string
	Check(key string) Result
	UpdateConfig(cfg Config)
}
