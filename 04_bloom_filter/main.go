package main

import (
	"github.com/AVVKavvk/bloom_filter/api"
	"github.com/AVVKavvk/bloom_filter/bloomFilter"

	_ "github.com/AVVKavvk/bloom_filter/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Bloom Filter API
// @version         1.0
// @description     This is a bloom filter service.
// @host            localhost:8080
// @BasePath        /
func main() {
	AlgoDryRun(1)

	// Init bloomFiler with 1 Kb size and 64 columns
	bloomFilter.InitBloomFilter(1, 64)

	e := echo.New()

	e.POST("/words", api.AddWord)
	e.POST("/words/check", api.CheckWeatherWordIsExist)

	// Route to serve the Swagger UI
	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
