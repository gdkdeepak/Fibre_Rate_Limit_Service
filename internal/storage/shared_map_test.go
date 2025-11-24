package storage

import (
	"sync"
	"testing"
	"time"
)

func TestShardedMap_ConcurrencyAndTTL(t *testing.T) {
	s := NewShardedMap(8, 0, time.Millisecond*50)
	defer s.Close()

	wg := sync.WaitGroup{}
	wg.Add(3)

	// Writer goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			s.Set("key", i, time.Millisecond*100)
			time.Sleep(time.Millisecond)
		}
	}()

	// Reader goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 500; i++ {
			s.Get("key")
			time.Sleep(time.Millisecond)
		}
	}()

	// Snapshot goroutine
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			s.Snapshot()
			time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Wait()

	// Allow janitor to clean expired entries
	time.Sleep(time.Millisecond * 200)

	snap := s.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected map empty after TTL expiry, got %d entries", len(snap))
	}
}
