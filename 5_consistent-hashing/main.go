package main

import (
	"github.com/AVVKavvk/consistent-hashing/algo"
	"github.com/AVVKavvk/consistent-hashing/api"
	_ "github.com/AVVKavvk/consistent-hashing/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Consistent Hashing API
// @version 1.0
// @description This is a sample server for managing users with consistent hashing.
// @host localhost:8080
// @BasePath /
func main() {

	algo.InitHashRing()

	e := echo.New()

	// Standard Logger Middleware
	e.Use(middleware.Logger())

	// Optional: Recover middleware to prevent crashes from panics
	e.Use(middleware.Recover())

	e.POST("/users", api.AddUser)
	e.GET("/users/:id", api.GetUserById)

	serverApi := e.Group("/servers")
	{
		serverApi.POST("", api.AddServer)
		serverApi.GET("", api.GetAllServer)
		serverApi.DELETE("/:name", api.DeleServer)
		serverApi.GET("/:name", api.GetServerInfo)
	}

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
