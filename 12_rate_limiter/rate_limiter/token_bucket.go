package rate_limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int
	tokens     float64
	refillRate float64 // tokens per minute
	LastRefill time.Time
	Mu         sync.Mutex
}

func GetNewTokenBucket(capacity int, refillRatePerMinute float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRatePerMinute,
		LastRefill: time.Now(),
	}
}

func (t *TokenBucket) Allow() bool {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	// Refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(t.LastRefill)

	// Calculate tokens to add based on elapsed time
	tokensToBeAdded := elapsed.Minutes() * t.refillRate

	if tokensToBeAdded > 0 {
		t.tokens = min(float64(t.capacity), t.tokens+tokensToBeAdded)
		t.LastRefill = now
	}

	if t.tokens >= 1.0 {
		t.tokens--
		return true
	}

	return false
}
