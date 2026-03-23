package api

import (
	"net/http"
	"strconv"

	"github.com/AVVKavvk/n+1-problem/db"
	"github.com/AVVKavvk/n+1-problem/models"
	"github.com/labstack/echo/v4"
)

func GetAuthorsNPlusOne(c echo.Context) error {
	var authors []models.Author
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid limit")
	}

	// The "1" Query: Fetch all authors
	database, err := db.GetDb()
	if err != nil {
		return err
	}
	database.Limit(int(limitInt)).Find(&authors)

	for i := range authors {
		// 2. The "N" Queries: Manually fetch books for EACH author inside the loop
		// This is the performance killer!
		database.Model(&authors[i]).Association("Books").Find(&authors[i].Books)
	}

	return c.JSON(http.StatusOK, authors)
}

func GetAuthors(c echo.Context) error {
	var authors []models.Author

	// .Preload("Books") handles the "Eager Loading"
	// GORM will execute:
	// 1. SELECT * FROM authors;
	// 2. SELECT * FROM books WHERE author_id IN (1, 2, 3...);
	database, err := db.GetDb()
	if err != nil {
		return err
	}
	if err := database.Preload("Books").Find(&authors).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, authors)
}
