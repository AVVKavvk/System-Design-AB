package api

import (
	"net/http"

	"github.com/AVVKavvk/mysql-replicas/database"
	"github.com/AVVKavvk/mysql-replicas/models"
	"github.com/labstack/echo/v4"
)

// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieve a list of users. This query is automatically routed to Read Replicas (Port 3307, 3308, 3309).
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func GetUsers(c echo.Context) error {
	var users []models.User

	// This looks like a normal query, but dbresolver sends it to a Replica
	db := database.MysqlDB.Find(&users)

	if db.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, db.Error.Error())
	}

	database.CheckDatabaseConnection(database.MysqlDB)

	return c.JSON(http.StatusOK, users)
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

	// This automatically goes to the Primary
	db := database.MysqlDB.Create(&user)

	if db.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, db.Error.Error())
	}

	database.CheckDatabaseConnection(database.MysqlDB)

	return c.JSON(http.StatusCreated, user)
}
