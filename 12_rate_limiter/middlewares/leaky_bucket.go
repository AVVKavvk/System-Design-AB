package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/AVVKavvk/rate_limiter/rate_limiter"
	"github.com/labstack/echo/v4"
)

var (
	LeakyBucketMiddleware echo.MiddlewareFunc
	leakyBuckets          = make(map[string]*rate_limiter.LeakyBucket)
	leakyBucketsMutex     sync.Mutex
)

func addLeakyBucketRateLimiter(capacity int, processRate time.Duration) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			// Get client identifier (IP address)
			clientIP := ctx.RealIP()

			leakyBucketsMutex.Lock()
			defer leakyBucketsMutex.Unlock()

			// Get or create bucket for this client
			bucket, exists := leakyBuckets[clientIP]

			if !exists {
				bucket = rate_limiter.GetLeakyBucket(capacity, processRate)
				leakyBuckets[clientIP] = bucket
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

func init() {
	// Initialize middleware with 100 capacity and 1 request/second => 60 requests/minute
	LeakyBucketMiddleware = addLeakyBucketRateLimiter(100, time.Second)
}
