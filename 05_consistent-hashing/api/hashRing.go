package api

import (
	"fmt"

	"github.com/AVVKavvk/consistent-hashing/models"
	"github.com/AVVKavvk/consistent-hashing/service"
	"github.com/labstack/echo/v4"
)

// AddServer godoc
// @Summary Add a new physical/virtual server to the ring
// @Description Adds a server and re-distributes keys based on consistent hashing
// @Tags servers
// @Accept  json
// @Produce  json
// @Param server body models.CreateServer true "Server Details"
// @Success 201 {object} map[string]interface{}
// @Router /servers [post]
func AddServer(ctx echo.Context) error {
	var server models.CreateServer
	if err := ctx.Bind(&server); err != nil {
		return err
	}
	result := service.AddServerService(&server)

	return ctx.JSON(201, result)
}

// DeleServer godoc
// @Summary Remove a server from the ring
// @Description Removes a server by name and migrates its keys to the next available node
// @Tags servers
// @Param name path string true "Server Name"
// @Success 200 {string} string "ok"
// @Router /servers/{name} [delete]
func DeleServer(ctx echo.Context) error {
	name := ctx.Param("name")
	result := service.DeleteServerService(name)
	return ctx.JSON(200, result)
}

// GetAllServer godoc
// @Summary List all active servers
// @Description Returns a list of all servers currently in the consistent hashing ring
// @Tags servers
// @Produce  json
// @Success 200 {array} map[string]interface{}
// @Router /servers [get]
func GetAllServer(ctx echo.Context) error {
	result := service.GetAllServerInfoService()
	return ctx.JSON(200, result)
}

// GetServerInfo godoc
// @Summary Get information about a specific server
// @Description Retrieves information about a specific server based on its name
// @Tags servers
// @Produce  json
// @Param name path string true "Server Name"
// @Success 200 {array} map[string]interface{}
// @Router /servers/{name} [get]
func GetServerInfo(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return ctx.JSON(400, "name is required")
	}
	fmt.Println(name)
	result := service.GetInfoOfServerByName(name)
	return ctx.JSON(200, result)
}
