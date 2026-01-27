package circuitbreaker

import (
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var (
	cb *gobreaker.CircuitBreaker
	mu sync.Mutex
)

func GetCircuitBreaker() *gobreaker.CircuitBreaker {
	mu.Lock()
	defer mu.Unlock()

	return cb
}
func InitCB() {
	settings := gobreaker.Settings{
		Name:        "CB Demo",
		MaxRequests: 3,                // Max requests allowed in half-open state
		Interval:    time.Second * 10, // Period to clear failure counts
		Timeout:     time.Second * 30, // Time to stay open before half-open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests) // counts.TotalFailures / counts.Requests
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	cb = gobreaker.NewCircuitBreaker(settings)
}
