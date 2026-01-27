package main

import (
	"github.com/AVVKavvk/cache_locking/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.RequestID())
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/products", api.GetAllProductForDashboard)

	e.Logger.Fatal(e.Start(":8080"))
}
