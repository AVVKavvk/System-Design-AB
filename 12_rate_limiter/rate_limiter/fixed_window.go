package rate_limiter

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindow struct {
	Limit  int
	Window time.Duration
	Counts map[int64]int
	Mu     sync.Mutex
}

func GetFixedWindow(limit int, window time.Duration) *FixedWindow {
	fw := &FixedWindow{
		Limit:  limit,
		Window: window,
		Counts: make(map[int64]int),
	}
	// Clean up old windows every window duration
	go fw.cleanup()

	return fw
}

func (fw *FixedWindow) Allow() bool {
	fw.Mu.Lock()
	defer fw.Mu.Unlock()

	currentWindow := time.Now().UnixNano() / int64(fw.Window)

	fmt.Println("currentWindow", currentWindow)

	if fw.Counts[currentWindow] < fw.Limit {
		fw.Counts[currentWindow]++
		return true
	}
	return false
}

func (fw *FixedWindow) cleanup() {
	ticker := time.NewTicker(fw.Window)
	defer ticker.Stop()

	for range ticker.C {
		fw.Mu.Lock()
		currentWindow := time.Now().UnixNano() / int64(fw.Window)

		// Remove windows older than current window
		for window := range fw.Counts {
			if window < currentWindow {
				delete(fw.Counts, window)
			}
		}
		fw.Mu.Unlock()
	}
}
