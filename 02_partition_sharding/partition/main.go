package main

import (
	"github.com/AVVKavvk/partition/api"
	"github.com/AVVKavvk/partition/config"
	"github.com/labstack/echo"
)

func main() {

	// Mysql Client
	config.InitMysqlDB()

	e := echo.New()

	e.POST("/users", api.CreateUser)
	e.GET("/users/:year", api.GetAllUserByYear)
	e.GET("/users/id/:id", api.GetUserById)
	e.GET("/years", api.GetAllYear)

	e.Logger.Fatal(e.Start(":8080"))
}
