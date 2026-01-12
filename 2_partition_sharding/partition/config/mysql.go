package config

import (
	"log"
	"sync"

	"github.com/AVVKavvk/partition/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDSN    = "root:secret@tcp(localhost:3306)/partition_db?charset=utf8mb4&parseTime=True&loc=Local"
	once        sync.Once
	MysqlClient *gorm.DB = nil
)

func InitMysqlDB() {
	once.Do(func() {
		var err error
		MysqlClient, err = gorm.Open(mysql.Open(MysqlDSN), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		err = MysqlClient.AutoMigrate(&models.User{})
		if err != nil {
			panic(err)
		}
		err = MysqlClient.AutoMigrate(&models.YearDetails{})
		if err != nil {
			panic(err)
		}

	})
}

func GetMysqlClient() *gorm.DB {
	if MysqlClient == nil {
		log.Fatalln("mysql client not initialized")
		panic("mysql client not initialized")

	}
	return MysqlClient
}
