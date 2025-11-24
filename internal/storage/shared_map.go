package storage

import (
	"hash/fnv"
	"sync"
	"time"
)

func fnv32(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// getShard selects a shard index based on the key's hash.
func (s *ShardedMap) getShard(key string) *shard {
	h := fnv32(key)
	return s.shards[h%uint32(len(s.shards))]
}

// Entry wraps a stored value and its optional expiration time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

type ShardedItem struct {
	Value      interface{}
	Expiration time.Time
}

func (s *ShardedMap) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.m[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (e Entry) isExpired(now time.Time) bool {
	if e.ExpiresAt.IsZero() {
		return false
	}
	return now.After(e.ExpiresAt)
}

// Each shard contains its own map + mutex.
type shard struct {
	mu sync.RWMutex
	m  map[string]Entry
}

// ShardedMap splits keys across N shards for concurrency.
type ShardedMap struct {
	shards     []*shard
	nShards    uint32
	defaultTTL time.Duration

	janitorStop chan struct{}
	stopOnce    sync.Once
}

// NewShardedMap initializes shards and starts janitor goroutine.
func NewShardedMap(nShards int, defaultTTL time.Duration, cleanupInterval time.Duration) *ShardedMap {
	if nShards <= 0 {
		nShards = 16
	}

	s := &ShardedMap{
		shards:      make([]*shard, nShards),
		nShards:     uint32(nShards),
		defaultTTL:  defaultTTL,
		janitorStop: make(chan struct{}),
	}

	for i := 0; i < nShards; i++ {
		s.shards[i] = &shard{
			m: make(map[string]Entry),
		}
	}

	go s.janitor(cleanupInterval)

	return s
}

func (s *ShardedMap) shardFor(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	i := h.Sum32() % s.nShards
	return s.shards[i]
}

// Set writes a value with optional TTL (pass -1 for NO expiry)
func (s *ShardedMap) Set(key string, val interface{}, ttl time.Duration) {
	var exp time.Time

	switch {
	case ttl == -1:
		// No expiry
	case ttl > 0:
		exp = time.Now().Add(ttl)
	case ttl == 0 && s.defaultTTL > 0:
		exp = time.Now().Add(s.defaultTTL)
	}

	sh := s.shardFor(key)
	sh.mu.Lock()
	sh.m[key] = Entry{Value: val, ExpiresAt: exp}
	sh.mu.Unlock()
}

// Get reads a value and purges if expired.
func (s *ShardedMap) Get(key string) (interface{}, bool) {
	sh := s.shardFor(key)
	now := time.Now()

	sh.mu.RLock()
	ent, found := sh.m[key]
	sh.mu.RUnlock()

	if !found {
		return nil, false
	}

	if ent.isExpired(now) {
		sh.mu.Lock()
		delete(sh.m, key)
		sh.mu.Unlock()
		return nil, false
	}

	return ent.Value, true
}

// Delete removes a key.
func (s *ShardedMap) Delete(key string) {
	sh := s.shardFor(key)
	sh.mu.Lock()
	delete(sh.m, key)
	sh.mu.Unlock()
}

// Snapshot returns a merged map copy of non-expired keys.
func (s *ShardedMap) Snapshot() map[string]Entry {
	out := make(map[string]Entry)
	now := time.Now()

	for _, sh := range s.shards {
		sh.mu.RLock()
		for k, v := range sh.m {
			if !v.isExpired(now) {
				out[k] = v
			}
		}
		sh.mu.RUnlock()
	}

	return out
}

// janitor periodically removes expired keys.
func (s *ShardedMap) janitor(interval time.Duration) {
	if interval <= 0 {
		interval = time.Second * 10
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for _, sh := range s.shards {
				sh.mu.Lock()
				for k, v := range sh.m {
					if v.isExpired(now) {
						delete(sh.m, k)
					}
				}
				sh.mu.Unlock()
			}

		case <-s.janitorStop:
			return
		}
	}
}

// Close stops the janitor goroutine.
func (s *ShardedMap) Close() {
	s.stopOnce.Do(func() {
		close(s.janitorStop)
	})
}
