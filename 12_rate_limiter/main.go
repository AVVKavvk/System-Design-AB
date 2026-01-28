package main

import (
	"net/http"

	"github.com/AVVKavvk/rate_limiter/middlewares"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(echoMiddleware.RequestID())

	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Recover())

	// rate limiting middleware

	// token bucket
	e.GET("/users/token-bucket", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success",
			"algo":    "token-bucket",
			"data": []map[string]interface{}{
				{
					"id":   1,
					"name": "John Doe",
				},
				{
					"id":   2,
					"name": "Maria Jones",
				},
			},
		})
	}, middlewares.TokenBucketMiddleware)

	// leaky bucket
	e.GET("/users/leaky-bucket", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success",
			"algo":    "leaky-bucket",
			"data": []map[string]interface{}{
				{
					"id":   1,
					"name": "John Doe",
				},
				{
					"id":   2,
					"name": "Maria Jones",
				},
			},
		})
	}, middlewares.LeakyBucketMiddleware)

	// fixed window
	e.GET("/users/fixed-window", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success",
			"algo":    "fixed-window",
			"data": []map[string]interface{}{
				{
					"id":   1,
					"name": "John Doe",
				},
				{
					"id":   2,
					"name": "Maria Jones",
				},
			},
		})
	}, middlewares.FixedWindowMiddleware)

	// sliding window counter
	e.GET("/users/sliding-window-counter", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"message": "Success",
			"algo":    "sliding-window-counter",
			"data": []map[string]interface{}{
				{
					"id":   1,
					"name": "John Doe",
				},
				{
					"id":   2,
					"name": "Maria Jones",
				},
			},
		})
	}, middlewares.SlidingWindowCounterMiddleware)

	e.Logger.Fatal(e.Start(":8080"))
}
