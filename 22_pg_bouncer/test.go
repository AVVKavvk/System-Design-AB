package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func Test(ctx echo.Context) error {
	var wg sync.WaitGroup
	requests := 20 // More than your 10-15 target

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Calling the /ping endpoint we created earlier
			resp, err := http.Get("http://localhost:8080/ping")
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", id, err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Request %d: %s\n", id, resp.Status)
		}(i)
	}
	wg.Wait()
	return nil
}
