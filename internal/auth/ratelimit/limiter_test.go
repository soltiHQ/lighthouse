package ratelimit

import (
	"errors"
	"testing"
	"time"

	"github.com/soltiHQ/control-plane/internal/auth"
)

func TestNew_DefaultConfig(t *testing.T) {
	l := New(Config{})

	if l.cfg.MaxAttempts != 5 {
		t.Fatalf("expected default MaxAttempts=5, got %d", l.cfg.MaxAttempts)
	}
	if l.cfg.BlockWindow != 10*time.Minute {
		t.Fatalf("expected default BlockWindow=10m, got %v", l.cfg.BlockWindow)
	}
}

func TestLimiter_NotBlockedInitially(t *testing.T) {
	l := New(Config{MaxAttempts: 3, BlockWindow: time.Minute})

	if l.Blocked("k", time.Now()) {
		t.Fatal("expected not blocked initially")
	}
}

func TestLimiter_BlockAfterThreshold(t *testing.T) {
	now := time.Now()
	l := New(Config{MaxAttempts: 3, BlockWindow: time.Minute})

	l.RecordFailure("k", now)
	l.RecordFailure("k", now)
	if l.Blocked("k", now) {
		t.Fatal("should not be blocked before threshold")
	}

	l.RecordFailure("k", now)
	if !l.Blocked("k", now) {
		t.Fatal("expected blocked after reaching threshold")
	}
}

func TestLimiter_CheckReturnsErrorWhenBlocked(t *testing.T) {
	now := time.Now()
	l := New(Config{MaxAttempts: 1, BlockWindow: time.Minute})

	l.RecordFailure("k", now)

	err := l.Check("k", now)
	if !errors.Is(err, auth.ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestLimiter_BlockExpires(t *testing.T) {
	start := time.Now()
	l := New(Config{MaxAttempts: 1, BlockWindow: time.Minute})

	l.RecordFailure("k", start)

	if !l.Blocked("k", start) {
		t.Fatal("expected blocked immediately")
	}

	after := start.Add(2 * time.Minute)

	if l.Blocked("k", after) {
		t.Fatal("expected block to expire")
	}

	// state should be cleared after expiry
	if _, ok := l.entries["k"]; ok {
		t.Fatal("expected entry to be removed after expiry")
	}
}

func TestLimiter_Reset(t *testing.T) {
	now := time.Now()
	l := New(Config{MaxAttempts: 1, BlockWindow: time.Minute})

	l.RecordFailure("k", now)

	if !l.Blocked("k", now) {
		t.Fatal("expected blocked")
	}

	l.Reset("k")

	if l.Blocked("k", now) {
		t.Fatal("expected not blocked after reset")
	}
}

func TestLimiter_IndependentKeys(t *testing.T) {
	now := time.Now()
	l := New(Config{MaxAttempts: 1, BlockWindow: time.Minute})

	l.RecordFailure("k1", now)

	if !l.Blocked("k1", now) {
		t.Fatal("k1 should be blocked")
	}
	if l.Blocked("k2", now) {
		t.Fatal("k2 should not be blocked")
	}
}

func TestLimiter_FailuresAccumulateUntilReset(t *testing.T) {
	now := time.Now()
	l := New(Config{MaxAttempts: 3, BlockWindow: time.Minute})

	l.RecordFailure("k", now)
	l.RecordFailure("k", now)

	if l.Blocked("k", now) {
		t.Fatal("should not be blocked yet")
	}

	l.RecordFailure("k", now)

	if !l.Blocked("k", now) {
		t.Fatal("expected blocked at exactly threshold")
	}
}

func TestLimiter_ExpiryResetsFailureState(t *testing.T) {
	start := time.Now()
	l := New(Config{MaxAttempts: 1, BlockWindow: time.Minute})

	l.RecordFailure("k", start)

	after := start.Add(2 * time.Minute)
	if l.Blocked("k", after) {
		t.Fatal("block should have expired")
	}

	l.RecordFailure("k", after)
	if !l.Blocked("k", after) {
		t.Fatal("should block again after new failure")
	}
}
