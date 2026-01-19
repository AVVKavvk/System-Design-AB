package api

import (
	"errors"

	"github.com/AVVKavvk/postgressql/models"
	"github.com/AVVKavvk/postgressql/service"
	"github.com/labstack/echo/v4"
)

// AddUser godoc
// @Summary Add a new user
// @Description Create a new user in database
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func AddUser(ctx echo.Context) error {
	var user models.User

	if err := ctx.Bind(&user); err != nil {
		return err
	}

	result, err := service.AddUserService(&user)
	if err != nil {
		return err
	}

	return ctx.JSON(201, result)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return errors.New("id is required")
	}

	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}

	result, err := service.UpdateUserService(id, &user)
	if err != nil {
		return err
	}

	return ctx.JSON(200, result)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Router /users/{id} [delete]
func DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return errors.New("id is required")
	}

	if err := service.DeleteUserService(id); err != nil {
		return err
	}

	return ctx.JSON(200, nil)
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Fetch a user using ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func GetUserById(ctx echo.Context) error {
	id := ctx.Param("id")

	if id == "" {
		return errors.New("id is required")
	}

	result, err := service.GetUserByIdService(id)
	if err != nil {
		return err
	}

	return ctx.JSON(200, result)
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Fetch all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User
// @Router /users [get]
func GetAllUsers(ctx echo.Context) error {
	result, err := service.GetAllUsersService()
	if err != nil {
		return err
	}

	return ctx.JSON(200, result)
}
