package main

import (
	"fmt"

	"github.com/AVVKavvk/system-design-ab/api"
	_ "github.com/AVVKavvk/system-design-ab/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title User API
// @version 1.0
// @description This is a sample server for managing users in MongoDB.
// @host localhost:8080
// @BasePath /
func main() {
	e := echo.New()

	fmt.Println("Server is running")

	e.GET("/users", api.GetAllDataFromDB)
	e.POST("/users", api.WriteDataToDB)

	// Swagger UI route
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())

	// Recovers from panics so the server doesn't crash
	e.Use(middleware.Recover())

	e.Logger.Fatal(e.Start(":8080"))
}
