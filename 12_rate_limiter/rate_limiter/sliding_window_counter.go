package rate_limiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	Limit      int
	Window     time.Duration
	PrevCount  int
	CurrCount  int
	LastWindow time.Time
	Mu         sync.Mutex
}

func GetSlidingWindowCounter(limit int, window time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		Limit:      limit,
		Window:     window,
		PrevCount:  0,
		CurrCount:  0,
		LastWindow: time.Now(),
	}
}

func (swc *SlidingWindowCounter) Allow() bool {
	swc.Mu.Lock()
	defer swc.Mu.Unlock()

	// Example:
	// prevCount = 4 (had 4 requests in previous window 0-10s)
	// currCount = 1 (had 1 request so far in current window 10-20s)
	// lastWindow = 10s (current window started at 10s)
	// now = 12s

	now := time.Now()

	//now.Sub(swc.lastWindow) = 12s - 10s = 2s (time elapsed in current window)

	elapsed := now.Sub(swc.LastWindow)

	if elapsed >= swc.Window {
		swc.PrevCount = swc.CurrCount
		swc.CurrCount = 0
		swc.LastWindow = now
	}

	// Calculate weighted request count
	// formula: prevCount * (remaining time in prev window / Window) + currCount

	// weight = (Window - elapsed) / Window
	//  (10s - 2s) / 10s
	//  8s / 10s
	//  0.8

	weight := float64(swc.Window-elapsed) / float64(swc.Window)

	// weightedCount = prevCount * weight + currCount
	//               = 4 * 0.8 + 1
	//               = 3.2 + 1
	//               = 4 (after int conversion)

	weightCount := int(swc.PrevCount*int(weight)) + swc.CurrCount

	if weightCount < swc.Limit {
		swc.CurrCount++
		return true
	}

	return false
}
