package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Prometheus middleware (collects metrics)
	e.Use(echoprometheus.NewMiddleware("myapp"))

	// 2. Add the /metrics endpoint (Prometheus scrapes this)
	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Monitoring!")
	})

	e.GET("/slow", func(c echo.Context) error {
		// Generate a random delay between 0 and 2 seconds
		delay := rand.Intn(2000)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		return c.String(http.StatusOK, "This request was slow!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
