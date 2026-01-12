package main

import (
	"github.com/AVVKavvk/mysql-replicas/api"
	"github.com/AVVKavvk/mysql-replicas/database"

	_ "github.com/AVVKavvk/mysql-replicas/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           MySQL Read Replica API
// @version         1.0
// @description     This is a sample server verifying Read/Write splitting with GORM.
// @host            localhost:8080
// @BasePath        /
func main() {
	database.InitDB()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger UI route
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// Routes
	e.GET("/users", api.GetUsers)    // READ -> Goes to Port 3307, 3308, or 3309
	e.POST("/users", api.CreateUser) // WRITE -> Goes to Port 3306

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
