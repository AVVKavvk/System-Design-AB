package main

import (
	"github.com/AVVKavvk/sharding/api"
	"github.com/AVVKavvk/sharding/config"
	"github.com/labstack/echo"
)

func main() {
	// Init Mysql
	config.InitMysqlDB()
	e := echo.New()

	e.GET("/users", api.GetAllUsers)
	e.POST("/users", api.CreateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
