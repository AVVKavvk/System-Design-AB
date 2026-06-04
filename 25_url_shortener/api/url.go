package api

import (
	"fmt"
	"net/http"

	"github.com/AVVKavvk/url-shortener/models"
	"github.com/AVVKavvk/url-shortener/pkg"
	"github.com/labstack/echo/v4"
)

func CreateShortUrl(ctx echo.Context) error {

	var body models.UrlShortenerRequest
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	up := pkg.UrlPkg{}
	resp, err := up.CreateLongUrl(&body)
	if err != nil {
		return err
	}
	return ctx.JSON(200, resp)
}

func GetLongUrl(ctx echo.Context) error {
	shortUrl := ctx.Param("short_url")
	up := pkg.UrlPkg{}
	longUrl, err := up.GetLongUrl(shortUrl)
	if err != nil {
		return err
	}
	fmt.Println("long url", longUrl)
	// 301 — permanent redirect (browser caches it)
	// 302 — temporary redirect (browser won't cache)
	return ctx.Redirect(http.StatusTemporaryRedirect, longUrl)
}
