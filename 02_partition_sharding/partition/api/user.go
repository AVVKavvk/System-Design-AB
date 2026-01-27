package api

import (
	"github.com/AVVKavvk/partition/models"
	"github.com/AVVKavvk/partition/service"
	"github.com/AVVKavvk/partition/utils"
	"github.com/labstack/echo"
)

func CreateUser(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}

	err := service.CreateUser(&user)
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))

	}
	return ctx.JSON(200, utils.SuccessResponse(user))
}

func GetAllUserByYear(ctx echo.Context) error {
	year := ctx.Param("year")
	if year == "" {
		return ctx.JSON(400, utils.ErrorResponse("Year not provided"))
	}
	result, err := service.GetAllUserByYear(year)
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, utils.SuccessResponse(result))
}
func GetUserById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(400, utils.ErrorResponse("Id not provided"))
	}
	result, err := service.GetUserById(id)
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, utils.SuccessResponse(result))
}
