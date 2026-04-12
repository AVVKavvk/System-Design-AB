package main

import (
	"context"
	"log"

	"github.com/AVVKavvk/s3-multiparts/s3_pkg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // Must be a slice
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		panic(err)
	}

	svc := &s3_pkg.S3Service{
		Client:        s3.NewFromConfig(cfg),
		PresignClient: s3.NewPresignClient(s3.NewFromConfig(cfg)),
	}

	// 1. Initialize Upload
	e.POST("/initiate-upload", svc.HandleInit)
	// 2. Get Presigned URL for a part
	e.POST("/get-presigned-url", svc.HandlePresign)
	// 3. Complete Upload
	e.POST("/complete-upload", svc.HandleComplete)

	e.Logger.Fatal(e.Start(":8080"))
}
