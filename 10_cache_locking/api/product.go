package api

import (
	"github.com/AVVKavvk/cache_locking/service"
	"github.com/AVVKavvk/cache_locking/utils"
	"github.com/labstack/echo/v4"
)

func GetAllProductForDashboard(ctx echo.Context) error {
	c, _ := utils.GetRequestContextAndIdFromEchoContext(ctx)

	result, err := service.GetAllProductsForDashboardService(c)

	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(200, result)
}
