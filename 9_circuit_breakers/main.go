package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AVVKavvk/circuit-breakers/circuitbreaker"
	"github.com/labstack/echo/v4"
	"github.com/sony/gobreaker"
)

func main() {
	fmt.Println("Running..........")

	circuitbreaker.InitCB()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.GET("/users", func(ctx echo.Context) error {
		return GetUser(ctx)
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func GetUser(ctx echo.Context) error {

	cb := circuitbreaker.GetCircuitBreaker()

	users, err := cb.Execute(func() (interface{}, error) {
		return handleExternalServiceCall()
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		switch err {
		case gobreaker.ErrOpenState:
			// Circuit is open - failing fast
			fmt.Println("Service unavailable, circuit open")
			return echo.ErrServiceUnavailable
		case gobreaker.ErrTooManyRequests:
			// Too many requests in half-open state
			fmt.Println("Too many requests")
			return echo.ErrTooManyRequests
		default:
			// Actual error from the service
			fmt.Printf("Service error: %v\n", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.JSON(200, map[string]interface{}{
		"message": "Success",
		"users":   users,
	})
}
func handleExternalServiceCall() ([]string, error) {

	return nil, errors.New("Created Error for testing CB")
}
