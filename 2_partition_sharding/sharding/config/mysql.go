package config

import (
	"errors"
	"log"
	"sync"

	"github.com/AVVKavvk/sharding/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDSNShard_0     = "root:secret@tcp(localhost:3306)/sharding_db?charset=utf8mb4&parseTime=True&loc=Local"
	MysqlDSNShard_1     = "root:secret@tcp(localhost:3307)/sharding_db?charset=utf8mb4&parseTime=True&loc=Local"
	MysqlDSNShard_2     = "root:secret@tcp(localhost:3308)/sharding_db?charset=utf8mb4&parseTime=True&loc=Local"
	once                sync.Once
	MysqlClient_Shard_0 *gorm.DB = nil
	MysqlClient_Shard_1 *gorm.DB = nil
	MysqlClient_Shard_2 *gorm.DB = nil
)

func InitMysqlDB() {
	once.Do(func() {
		var err error
		MysqlClient_Shard_0, err = gorm.Open(mysql.Open(MysqlDSNShard_0), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		MysqlClient_Shard_1, err = gorm.Open(mysql.Open(MysqlDSNShard_1), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		MysqlClient_Shard_2, err = gorm.Open(mysql.Open(MysqlDSNShard_2), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		err = MysqlClient_Shard_0.AutoMigrate(&models.User{})
		if err != nil {
			panic(err)
		}

		err = MysqlClient_Shard_1.AutoMigrate(&models.User{})
		if err != nil {
			panic(err)
		}
		err = MysqlClient_Shard_2.AutoMigrate(&models.User{})
		if err != nil {
			panic(err)
		}

	})
}

func GetMysqlClient(index int) (*gorm.DB, error) {
	switch index {
	case 0:
		if MysqlClient_Shard_0 == nil {
			log.Fatal("Mysql Not Init")
		}
		return MysqlClient_Shard_0, nil
	case 1:
		if MysqlClient_Shard_1 == nil {
			log.Fatal("Mysql Not Init")
		}
		return MysqlClient_Shard_1, nil
	case 2:
		if MysqlClient_Shard_2 == nil {
			log.Fatal("Mysql Not Init")
		}
		return MysqlClient_Shard_2, nil

	default:
		return nil, errors.New("Index is not matched any shard")
	}

}
