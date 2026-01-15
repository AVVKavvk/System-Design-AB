package redisClient

import (
	"encoding/json"

	"github.com/AVVKavvk/consistent-hashing/models"
)

// StoreUserDataWithHashToRedisWithNode
func StoreUserDataWithHashToRedisWithNode(nodeName string, hash string, user *models.User) error {
	client := GetRedisClient()

	key := nodeName + ":" + hash

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return client.Set(key, userJSON, 0).Err()

}

func GetUserDataWithHashFromRedisWithNode(nodeName string, hash string) (*models.User, error) {
	client := GetRedisClient()
	key := nodeName + ":" + hash
	var user models.User
	val, err := client.Get(key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserDataWithHashFromRedisWithNode(nodeName string, hash string) error {
	client := GetRedisClient()
	key := nodeName + ":" + hash
	err := client.Del(key).Err()
	if err != nil {
		return err
	}
	return nil
}
