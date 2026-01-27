package api

import (
	"github.com/AVVKavvk/consistent-hashing/models"
	"github.com/AVVKavvk/consistent-hashing/service"
	"github.com/labstack/echo/v4"
)

// AddUser godoc
// @Summary Add a new user
// @Description Takes a user object and stores it using consistent hashing
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func AddUser(ctx echo.Context) error {
	var user models.User

	if err := ctx.Bind(&user); err != nil {
		return err
	}
	result, err := service.AddUserDataService(&user)
	if err != nil {
		return err
	}
	return ctx.JSON(201, result)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Retrieves a user's data based on their unique ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserById(ctx echo.Context) error {
	id := ctx.Param("id")
	result, err := service.GetUserByIdService(id)
	if err != nil {
		return err
	}
	return ctx.JSON(200, result)
}
