package service

import (
	"fmt"
	"net/http"

	"github.com/AVVKavvk/sql-injection/database"
	"github.com/labstack/echo/v4"
)

func Login(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	db, err := database.GetDb()
	if err != nil {
		return err
	}
	// DANGER: Concatenating strings directly into the query
	query := fmt.Sprintf("SELECT id FROM users WHERE email=%s AND password=%s", email, password)

	fmt.Println("Query:", query)

	var id int
	err = db.QueryRow(ctx.Request().Context(), query).Scan(&id)

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid login"})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Logged in!", "id": fmt.Sprint(id)})
}

func LoginSafe(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	db, err := database.GetDb()
	if err != nil {
		return err
	}

	query := "SELECT id FROM users WHERE email=$1 AND password=$2"
	fmt.Println("Query:", query)

	var id int
	err = db.QueryRow(ctx.Request().Context(), query, email, password).Scan(&id)

	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid login"})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Logged in!", "id": fmt.Sprint(id)})
}
