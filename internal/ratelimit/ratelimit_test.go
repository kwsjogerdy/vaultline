package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_BurstPermitted(t *testing.T) {
	l := New(1, 3)
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("expected allow on token %d, got: %v", i, err)
		}
	}
}

func TestAllow_ExceedsBurst(t *testing.T) {
	l := New(1, 2)
	l.Allow()
	l.Allow()
	if err := l.Allow(); err == nil {
		t.Fatal("expected rate limit error after burst exhausted")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	now := time.Unix(1000, 0)
	clock := func() time.Time { return now }

	l := newWithClock(2, 1, clock)
	l.Allow() // consume the 1 burst token

	if err := l.Allow(); err == nil {
		t.Fatal("expected rate limit before refill")
	}

	// advance time by 1 second → +2 tokens, capped at max=1
	now = now.Add(1 * time.Second)
	if err := l.Allow(); err != nil {
		t.Fatalf("expected allow after refill, got: %v", err)
	}
}

func TestAllow_TokensCapAtMax(t *testing.T) {
	now := time.Unix(1000, 0)
	clock := func() time.Time { return now }

	l := newWithClock(5, 3, clock)
	l.Allow()
	l.Allow()
	l.Allow()

	// advance 100 seconds — tokens should cap at max=3, not overflow
	now = now.Add(100 * time.Second)
	for i := 0; i < 3; i++ {
		if err := l.Allow(); err != nil {
			t.Fatalf("expected allow %d after large refill, got: %v", i, err)
		}
	}
	if err := l.Allow(); err == nil {
		t.Fatal("expected rate limit after consuming all capped tokens")
	}
}
