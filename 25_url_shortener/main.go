package main

import (
	"github.com/AVVKavvk/url-shortener/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	router := echo.New()

	router.Use(middleware.Recover())

	server.AddRoutes(router)

	router.Logger.Fatal(router.Start(":" + "8888"))
}
