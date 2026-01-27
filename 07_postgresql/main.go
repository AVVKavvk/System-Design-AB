package main

import (
	"github.com/AVVKavvk/postgressql/api"
	"github.com/AVVKavvk/postgressql/config"
	_ "github.com/AVVKavvk/postgressql/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title PostgresSQL
// @version 1.0
// @description This is a sample server for managing users with postgresSQL
// @host localhost:8080
// @BasePath /
func main() {
	// fmt.Println("Hello")

	// Adding users table
	config.MyAppTables.Add("users")

	e := echo.New()

	// Standard Logger Middleware
	e.Use(middleware.Logger())

	// Optional: Recover middleware to prevent crashes from panics
	e.Use(middleware.Recover())

	users := e.Group("/users")
	{
		users.POST("", api.AddUser)
		users.GET("", api.GetAllUsers)
		users.GET("/:id", api.GetUserById)
		users.PUT("/:id", api.UpdateUser)
		users.DELETE("/:id", api.DeleteUser)
	}

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
