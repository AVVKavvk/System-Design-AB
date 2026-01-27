package service

import (
	"github.com/AVVKavvk/partition/config"
	"github.com/AVVKavvk/partition/models"
)

func AddYearToMysql(year string) error {
	db := config.GetMysqlClient()
	yearDetails := models.YearDetails{YEAR: year}
	return db.Create(yearDetails).Error
}

func GetAllYear() ([]models.YearDetails, error) {
	db := config.GetMysqlClient()
	var yearDetails []models.YearDetails
	return yearDetails, db.Find(&yearDetails).Error
}
