package server

import (
	"github.com/AVVKavvk/url-shortener/api"
	"github.com/labstack/echo/v4"
)

func AddRoutes(router *echo.Echo) {

	urls := router.Group("/api/urls")
	{
		urls.POST("/shorten", api.CreateShortUrl)
		urls.GET("/:short_url", api.GetLongUrl)
	}

}
