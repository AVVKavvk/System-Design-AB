package api

import (
	"errors"

	"github.com/AVVKavvk/sharding/models"
	"github.com/AVVKavvk/sharding/services"
	"github.com/AVVKavvk/sharding/utils"
	"github.com/labstack/echo"
)

func CreateUser(ctx echo.Context) error {
	clientId := ctx.Request().Header.Get("client-x-id")

	if clientId == "" {
		return errors.New("Client-X-Id is missing in header")
	}
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(400, utils.ErrorResponse(err.Error()))
	}
	result, err := services.CreateUserService(clientId, &user)
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))
	}
	return ctx.JSON(201, utils.SuccessResponse(result))
}

func GetAllUsers(ctx echo.Context) error {
	clientId := ctx.Request().Header.Get("client-x-id")
	if clientId == "" {
		return errors.New("Client-X-Id is missing in header")
	}
	result, err := services.GetAllUserService(clientId)
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, utils.SuccessResponse(result))
}
