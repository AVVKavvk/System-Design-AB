package api

import (
	"github.com/AVVKavvk/partition/service"
	"github.com/AVVKavvk/partition/utils"
	"github.com/labstack/echo"
)

func GetAllYear(ctx echo.Context) error {
	result, err := service.GetAllYear()
	if err != nil {
		return ctx.JSON(500, utils.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, utils.SuccessResponse(result))
}
