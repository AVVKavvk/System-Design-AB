package service

import (
	"errors"
	"fmt"

	"github.com/AVVKavvk/partition/config"
	"github.com/AVVKavvk/partition/helper"
	"github.com/AVVKavvk/partition/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	db := config.GetMysqlClient()

	tableName := helper.GetTableName(user.YEAR)

	// If "users_2025" doesn't exist, GORM will create it based on the User struct
	if !db.Migrator().HasTable(tableName) {
		err := db.Table(tableName).AutoMigrate(&models.User{})
		if err != nil {
			return err
		} else {
			// YEAR in YYYY
			err := AddYearToMysql(user.YEAR)
			if err != nil {
				return err
			}
		}
	}
	return db.Table(tableName).Create(user).Error

}

func GetUserById(id string) (*models.User, error) {

	years, err := GetAllYear()

	if err != nil {
		return nil, err
	}

	for _, y := range years {
		db := config.GetMysqlClient()
		tableName := helper.GetTableName(y.YEAR)
		var user models.User
		err = db.Table(tableName).Where("id = ?", id).First(&user).Error
		fmt.Println(err)

		if err == nil {
			//  found user
			return &user, nil
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			//  not in this partition → continue
			continue
		}

		// real DB error
		return nil, err

	}
	return nil, errors.New("User not found")
}

func GetAllUserByYear(year string) (*[]models.User, error) {

	years, err := GetAllYear()
	var currentYear string
	if err != nil {
		return nil, err
	}

	for _, y := range years {
		if y.YEAR == year {
			currentYear = y.YEAR
			break
		}
	}
	if currentYear == "" {
		return nil, errors.New("Provided Year not found")
	}

	db := config.GetMysqlClient()
	tableName := helper.GetTableName(year)
	var users []models.User
	return &users, db.Table(tableName).Find(&users).Error

}
