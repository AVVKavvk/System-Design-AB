package db

import (
	"log"
	"os"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDb() (*gorm.DB, error) {
	// ... inside your DB connection function
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel: logger.Info, // This MUST be Info to see the SQL
			Colorful: true,        // Makes it pretty in your Ubuntu terminal
		},
	)
	dsn := "postgres://myuser:mypassword@localhost:5432/nplusone?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {

		log.Fatal(err)

		return nil, err
	}

	return db, nil
}
