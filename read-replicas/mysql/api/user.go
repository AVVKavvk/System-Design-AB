package api

import (
	"net/http"

	"github.com/AVVKavvk/mysql-replicas/models"
	"github.com/AVVKavvk/mysql-replicas/service"
	"github.com/labstack/echo/v4"
)

// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of users. This query is automatically routed to Read Replicas (Port 3307, 3308, 3309).
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_token  header    string  false  "User Token for identifying specific user shard/replica"
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func GetUsers(c echo.Context) error {

	userToken := c.Request().Header.Get("user_token")
	if userToken == "" {
		userToken = "defaultUserToken"
	}
	result, err := service.GetUsersService(c.Request().Context(), userToken)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user. This query is automatically routed to the Primary DB (Port 3306).
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Data"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users [post]
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	err := service.CreateUserService(c.Request().Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, user)
}
