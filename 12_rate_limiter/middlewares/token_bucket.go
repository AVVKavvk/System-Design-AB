package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AVVKavvk/rate_limiter/rate_limiter"
	"github.com/labstack/echo/v4"
)

var (
	TokenBucketMiddleware echo.MiddlewareFunc
	tokenBuckets          = make(map[string]*rate_limiter.TokenBucket)
	tokenBucketsMutex     sync.Mutex
)

func addTokenBucketRateLimiter(capacity int, refillRatePerMinute float64) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get client identifier (IP address)
			clientIP := ctx.RealIP()

			// Get or create bucket for this client
			tokenBucketsMutex.Lock()
			defer tokenBucketsMutex.Unlock()

			bucket, exists := tokenBuckets[clientIP]

			if !exists {
				bucket = rate_limiter.GetNewTokenBucket(capacity, refillRatePerMinute)
				tokenBuckets[clientIP] = bucket

			}

			// Check if request is allowed
			if !bucket.Allow() {
				return ctx.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded. Please try again later.",
				})
			}

			// Request is allowed, proceed to next handler
			return next(ctx)

		}
	}

}

func cleanupInactiveTokenBuckets(interval time.Duration) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("Recovered from panic for cleanupInactiveTokenBuckets: %v\n", r)

			// Recover from panic and again start the cleanup process
			go cleanupInactiveTokenBuckets(interval)
		}
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		tokenBucketsMutex.Lock()
		now := time.Now()
		for ip, bucket := range tokenBuckets {
			bucket.Mu.Lock()
			// Remove buckets inactive for more than 1 hour
			if now.Sub(bucket.LastRefill) > time.Hour {
				delete(tokenBuckets, ip)
			}
			bucket.Mu.Unlock()
		}
		tokenBucketsMutex.Unlock()
	}
}

func init() {
	// Initialize middleware with 100 capacity and 60 requests/minute
	TokenBucketMiddleware = addTokenBucketRateLimiter(100, 60.0)

	// cleanup
	go cleanupInactiveTokenBuckets(10 * time.Minute)
}
