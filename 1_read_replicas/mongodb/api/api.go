package api

import (
	"github.com/AVVKavvk/system-design-ab/models"
	"github.com/AVVKavvk/system-design-ab/service"
	"github.com/labstack/echo/v4"
)

// WriteDataToDB godoc
// @Summary Create a new user
// @Description Save user data to MongoDB
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body      models.User  true  "User Data"
// @Success 201   {object}  models.User
// @Failure 400   {object}  map[string]string
// @Router /users [post]
func WriteDataToDB(ctx echo.Context) error {
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	result, err := service.WriteDataToDBService(user)
	if err != nil {
		return err
	}
	return ctx.JSON(201, result)

}

// GetAllDataFromDB godoc
// @Summary Get all users
// @Description Retrieve all user records from MongoDB
// @Tags users
// @Produce  json
// @Success 200 {array}   models.User
// @Router /users [get]
func GetAllDataFromDB(ctx echo.Context) error {
	users, err := service.GetAllDataFromDBService()
	if err != nil {
		return err
	}
	return ctx.JSON(200, users)
}
