// Package ratelimit provides a simple token-bucket rate limiter for Vault API calls.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter controls the rate of operations using a token bucket.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter allowing burst up to max tokens, refilling at rate tokens/sec.
func New(ratePerSec float64, burst int) *Limiter {
	now := time.Now()
	return &Limiter{
		tokens:   float64(burst),
		max:      float64(burst),
		rate:     ratePerSec,
		lastTick: now,
		clock:    time.Now,
	}
}

// newWithClock is used for testing with an injectable clock.
func newWithClock(ratePerSec float64, burst int, clock func() time.Time) *Limiter {
	l := New(ratePerSec, burst)
	l.clock = clock
	l.lastTick = clock()
	l.tokens = float64(burst)
	return l
}

// Allow returns nil if the operation is permitted, or an error if rate-limited.
func (l *Limiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return fmt.Errorf("rate limit exceeded: retry after %.2fs", (1-l.tokens)/l.rate)
	}

	l.tokens--
	return nil
}

// Wait blocks until a token is available or the context deadline is reached.
func (l *Limiter) Wait() {
	for {
		if err := l.Allow(); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
