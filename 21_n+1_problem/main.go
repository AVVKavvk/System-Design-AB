package main

import (
	"os"
	"time"

	"github.com/AVVKavvk/n+1-problem/api"
	"github.com/AVVKavvk/n+1-problem/db"
	"github.com/AVVKavvk/n+1-problem/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	e := echo.New()

	zerolog.TimeFieldFormat = time.RFC3339
	logger := log.Output(os.Stdout)

	database, err := db.GetDb()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
		return
	}
	// Auto-migrate creates the tables based on your structs
	database.AutoMigrate(&models.Author{}, &models.Book{})

	// Run the seed
	db.Seed(database)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// ── What to capture ────────────────────────────────────────────────
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: true,
		LogResponseSize:  true,

		// Capture specific headers (canonical form required)
		LogHeaders: []string{
			"Accept",
			"Accept-Encoding",
			"Content-Type",
			"Authorization", // ⚠ redact in production — see LogValuesFunc below
		},

		// Capture specific query params and form values
		LogQueryParams: []string{"page", "limit", "sort"},
		LogFormValues:  []string{"username"},

		// ── Skip health-check / readiness probes from logs ─────────────────
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},

		// ── Called BEFORE the handler — useful for injecting a request ID ──
		BeforeNextFunc: func(c echo.Context) {
			c.Set("start_time", time.Now())
		},

		// ── Global error handler integration ───────────────────────────────
		// Setting true means echo's global error handler runs first, so
		// LogStatus reflects the final HTTP status after error mapping.
		HandleError: true,

		// ── Write the log entry ────────────────────────────────────────────
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			event := logger.Info()

			// Downgrade 4xx/5xx to the right log level
			switch {
			case v.Status >= 500:
				event = logger.Error()
			case v.Status >= 400:
				event = logger.Warn()
			}

			// Core request fields
			event = event.
				Str("request_id", v.RequestID).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("path", v.URIPath).
				Str("route", v.RoutePath).
				Str("host", v.Host).
				Str("remote_ip", v.RemoteIP).
				Str("protocol", v.Protocol).
				Str("referer", v.Referer).
				Str("user_agent", v.UserAgent).
				Str("content_length", v.ContentLength).
				Int("status", v.Status).
				Int64("response_bytes", v.ResponseSize).
				Dur("latency_ms", v.Latency)

			// Log error message (if any)
			if v.Error != nil {
				event = event.Err(v.Error)
			}

			// Captured headers — redact sensitive values
			if len(v.Headers) > 0 {
				headers := event.Dict("headers", zerolog.Dict())
				for key, vals := range v.Headers {
					if key == "Authorization" {
						headers = headers.Strs(key, []string{"[REDACTED]"})
						continue
					}
					headers = headers.Strs(key, vals)
				}
				_ = headers
			}

			// Captured query params
			if len(v.QueryParams) > 0 {
				qp := zerolog.Dict()
				for k, vals := range v.QueryParams {
					qp = qp.Strs(k, vals)
				}
				event = event.Dict("query", qp)
			}

			// Captured form values
			if len(v.FormValues) > 0 {
				fv := zerolog.Dict()
				for k, vals := range v.FormValues {
					fv = fv.Strs(k, vals)
				}
				event = event.Dict("form", fv)
			}

			event.Msg("request")
			return nil
		},
	}))

	e.GET("/nplusone", api.GetAuthorsNPlusOne)
	e.GET("/safe", api.GetAuthors)
	e.Logger.Fatal(e.Start(":8080"))
}
