package main

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	// Connect to PgBouncer instead of Postgres directly
	// Note: sslmode=disable for local dev; target port 6432
	db, err := sql.Open("pgx", os.Getenv("PG_BOUNCER_DATABASE_URL"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	e.GET("/ping", func(c echo.Context) error {
		var greeting string
		var sleep string
		err := db.QueryRow("SELECT pg_sleep(2), 'Hello'").Scan(&sleep, &greeting)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": greeting})
	})

	e.GET("/test", Test)

	e.Logger.Fatal(e.Start(":" + port))
}
