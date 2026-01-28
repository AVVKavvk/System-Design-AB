package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/AVVKavvk/rate_limiter/rate_limiter"
	"github.com/labstack/echo/v4"
)

var (
	FixedWindowMiddleware echo.MiddlewareFunc
	fixedWindows          = make(map[string]*rate_limiter.FixedWindow)
	fixedWindowsMutex     sync.Mutex
)

func addFixedWindowRateLimiter(limit int, window time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			// Get client identifier (IP address)
			clientIP := ctx.RealIP()

			fixedWindowsMutex.Lock()
			defer fixedWindowsMutex.Unlock()

			// Get or create bucket for this client
			bucket, exists := fixedWindows[clientIP]

			if !exists {
				bucket = rate_limiter.GetFixedWindow(limit, window)
				fixedWindows[clientIP] = bucket
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
	//For 2 requests per second:
	// Window changes every 1 second
	// currentWindow = time.Now().UnixNano() / 1000000000

	FixedWindowMiddleware = addFixedWindowRateLimiter(2, time.Second)

	// For 10 requests per 5 seconds:
	// Window changes every 5 seconds
	// currentWindow = time.Now().UnixNano() / 5000000000

	// FixedWindowMiddleware = addFixedWindowRateLimiter(10, 5*time.Second)

	// For 60 requests per minute:
	// Window changes every 60 seconds
	// currentWindow = time.Now().UnixNano() / 60000000000

	// FixedWindowMiddleware = addFixedWindowRateLimiter(60, time.Minute)

}
