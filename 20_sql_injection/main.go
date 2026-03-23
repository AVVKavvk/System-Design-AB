package main

import (
	"github.com/AVVKavvk/sql-injection/service"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// VULNERABLE LOGIN ROUTE
	e.POST("/login", service.Login)
	e.POST("/login-safe", service.LoginSafe)
	e.Logger.Fatal(e.Start(":8080"))
}
