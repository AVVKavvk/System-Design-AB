package service

import (
	"context"
	"net/http"

	"github.com/AVVKavvk/mysql-replicas/database"
	"github.com/AVVKavvk/mysql-replicas/models"
	"github.com/labstack/echo/v4"
)

func GetUsersService(c context.Context, userToken string) ([]models.User, error) {
	var users []models.User

	db := database.GetReplicaByKey(userToken)

	if db.Error != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, db.Error.Error())
	}

	if err := db.Find(&users).Error; err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	database.CheckDatabaseConnection(db)

	return users, nil
}

func CreateUserService(c context.Context, user *models.User) error {
	// This automatically goes to the Primary
	db := database.PrimaryDB.Create(&user)

	if db.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, db.Error.Error())
	}

	// database.CheckDatabaseConnection(database.MysqlDB)

	return nil
}
