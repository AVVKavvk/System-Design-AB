package service

import (
	"strconv"

	"github.com/AVVKavvk/consistent-hashing/algo"
	"github.com/AVVKavvk/consistent-hashing/models"
	"github.com/AVVKavvk/consistent-hashing/redisClient"
)

func AddUserDataService(user *models.User) (*models.ResponseModel, error) {
	id := user.ID
	hr := algo.GetHashRing()

	serverName, hash := hr.GetOwner(id)
	hashStr := strconv.Itoa(hash)

	hr.AddHashToNode(serverName, hash)
	hr.AddUserIdToNode(serverName, id)

	err := redisClient.StoreUserDataWithHashToRedisWithNode(serverName, hashStr, user)
	if err != nil {
		return nil, err
	}

	return &models.ResponseModel{ServerName: serverName, Hash: hashStr, Users: []models.User{*user}}, nil
}

func GetUserByIdService(userId string) (*models.ResponseModel, error) {
	hr := algo.GetHashRing()
	serverName, hash := hr.GetOwner(userId)
	hashStr := strconv.Itoa(hash)
	user, err := redisClient.GetUserDataWithHashFromRedisWithNode(serverName, hashStr)
	if err != nil {
		return nil, err
	}
	return &models.ResponseModel{ServerName: serverName, Hash: hashStr, Users: []models.User{*user}}, nil
}
