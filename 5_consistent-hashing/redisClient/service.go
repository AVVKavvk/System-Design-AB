package redisClient

import (
	"encoding/json"
	"log"

	"github.com/AVVKavvk/consistent-hashing/models"
)

// AppendDataToNode adds a user to a List named after the node
func AppendDataToNode(nodeName string, user models.User) error {
	client := GetRedisClient()

	// Serialize struct to JSON string
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// RPush to add to a list (creates the list if it doesn't exist)
	return client.RPush(nodeName, data).Err()
}

// AppendBulkData adds multiple users to a List named after the node
func AppendBulkData(nodeName string, users []*models.User) error {
	client := GetRedisClient()
	for _, user := range users {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = client.RPush(nodeName, data).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAllTheDataForNode retrieves all users stored in the node's list
func GetAllTheDataForNode(nodeName string) ([]*models.User, error) {
	client := GetRedisClient()

	// Get all elements from the list (0 to -1 means all)
	records, err := client.LRange(nodeName, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var users []*models.User
	for _, record := range records {
		var user models.User
		// Unmarshal JSON back into the struct
		if err := json.Unmarshal([]byte(record), &user); err != nil {
			log.Printf("failed to unmarshal user: %v", err)
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}
