package main

import (
	"github.com/AVVKavvk/protobuf/api"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.POST("/users", api.CreateUser)

	e.GET("/users/:userId", api.GetUserById)
	e.Logger.Fatal(e.Start(":8080"))
}
