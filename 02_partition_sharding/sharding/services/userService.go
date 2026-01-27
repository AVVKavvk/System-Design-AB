package services

import (
	"github.com/AVVKavvk/sharding/config"
	"github.com/AVVKavvk/sharding/helper"
	"github.com/AVVKavvk/sharding/models"
)

func CreateUserService(clientId string, user *models.User) (*models.ResponseWithShard, error) {

	index := helper.GetShardIndexByClientId(clientId)
	mysqlClient, err := config.GetMysqlClient(index)
	if err != nil {
		return nil, err
	}
	err = mysqlClient.Create(user).Error
	if err != nil {
		return nil, err
	}
	return &models.ResponseWithShard{Shard: index, Data: []models.User{*user}}, nil
}

func GetAllUserService(clientId string) (*models.ResponseWithShard, error) {
	index := helper.GetShardIndexByClientId(clientId)
	mysqlClient, err := config.GetMysqlClient(index)
	if err != nil {
		return nil, err
	}
	var users []models.User
	err = mysqlClient.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &models.ResponseWithShard{Shard: index, Data: users}, nil
}
