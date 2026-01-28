package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/AVVKavvk/rate_limiter/rate_limiter"
	"github.com/labstack/echo/v4"
)

var (
	SlidingWindowCounterMiddleware echo.MiddlewareFunc
	slidingWindowCounterBuckets    = make(map[string]*rate_limiter.SlidingWindowCounter)
	slidingWindowCounterMutex      sync.Mutex
)

func addSlidingWindowCounterRateLimiter(limit int, windowSize time.Duration) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			// Get client identifier (IP address)
			clientIP := ctx.RealIP()

			slidingWindowCounterMutex.Lock()
			defer slidingWindowCounterMutex.Unlock()

			// Get or create bucket for this client
			bucket, exists := slidingWindowCounterBuckets[clientIP]

			if !exists {
				bucket = rate_limiter.GetSlidingWindowCounter(limit, windowSize)
				slidingWindowCounterBuckets[clientIP] = bucket
			}

			// Check if request is allowed
			if !bucket.Allow() {
				return ctx.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "Rate limit exceeded. Please try again later.",
				})
			}

			return next(ctx)
		}
	}
}

func init() {
	//  5 requests per 10 seconds
	SlidingWindowCounterMiddleware = addSlidingWindowCounterRateLimiter(5, 10*time.Second)
}
