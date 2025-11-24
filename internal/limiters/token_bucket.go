package limiters

import (
	"time"

	"fibre_rate_limit_service/internal/storage"
)

// TokenBucketConfig defines the input configuration.
type TokenBucketConfig struct {
	Name        string
	Capacity    int           // max tokens
	RefillRate  int           // tokens per interval
	RefillEvery time.Duration // interval for refilling tokens
	TTL         time.Duration // bucket TTL in storage
}

// TokenBucket implements the Limiter interface.
type TokenBucket struct { // ✅ uppercase
	cfg   TokenBucketConfig
	store *storage.ShardedMap
}

type bucketState struct {
	Tokens     int
	LastRefill time.Time
}

// NewTokenBucket creates a new token-bucket limiter.
func NewTokenBucket(cfg TokenBucketConfig, store *storage.ShardedMap) Limiter {
	return &TokenBucket{
		cfg:   cfg,
		store: store,
	}
}

// Name returns the limiter name.
func (tb *TokenBucket) Name() string {
	return tb.cfg.Name
}

// Check consumes 1 token if available and returns a Result.
func (tb *TokenBucket) Check(key string) Result {
	now := time.Now()

	raw, _ := tb.store.Get(key)
	var state bucketState
	if raw == nil {
		state = bucketState{
			Tokens:     tb.cfg.Capacity,
			LastRefill: now,
		}
	} else {
		state = raw.(bucketState)
	}

	// Refill tokens
	elapsed := now.Sub(state.LastRefill)
	if elapsed >= tb.cfg.RefillEvery {
		refills := int(elapsed / tb.cfg.RefillEvery)
		state.Tokens += refills * tb.cfg.RefillRate
		if state.Tokens > tb.cfg.Capacity {
			state.Tokens = tb.cfg.Capacity
		}
		state.LastRefill = now
	}

	// Check if allowed
	allowed := false
	if state.Tokens > 0 {
		state.Tokens--
		allowed = true
	}

	tb.store.Set(key, state, tb.cfg.TTL)

	return Result{
		Allowed:   allowed,
		Remaining: state.Tokens,
		ResetAt:   state.LastRefill.Add(tb.cfg.RefillEvery),
		Reason:    "",
	}
}

// UpdateConfig updates the limiter's config.
func (tb *TokenBucket) UpdateConfig(cfg Config) {
	tb.cfg.Capacity = cfg.Limit
	tb.cfg.RefillRate = cfg.Limit / int(cfg.Window.Seconds()) // approximate refill rate
	tb.cfg.RefillEvery = cfg.Window
}

// GetState returns the current state of the token bucket for a given key
func (tb *TokenBucket) GetState(key string) bucketState {
	raw, _ := tb.store.Get(key)
	if raw == nil {
		return bucketState{
			Tokens:     tb.cfg.Capacity,
			LastRefill: time.Now(),
		}
	}
	return raw.(bucketState)
}

func (tb *TokenBucket) GetConfig() TokenBucketConfig {
	return tb.cfg
}

func (tb *TokenBucket) StoreSnapshot() map[string]interface{} {
	result := make(map[string]interface{})
	data := tb.store.Snapshot()
	for k, v := range data {
		result[k] = v
	}
	return result
}
