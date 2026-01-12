package service

import (
	"context"
	"fmt"

	"github.com/AVVKavvk/system-design-ab/config"
	"github.com/AVVKavvk/system-design-ab/constant"
	"github.com/AVVKavvk/system-design-ab/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func WriteDataToDBService(user models.User) (any, error) {
	client, err := config.GetMongoClient()
	if err != nil {
		fmt.Println("Error while getting mongo client")
		return nil, err
	}
	db := constant.TEST_DB
	coll := constant.TEST_USERS_COLLECTION

	userCollection := client.Database(db).Collection(coll)

	res, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error while inserting data")
		return nil, err
	}
	// println(res.InsertedID)
	return res, nil
}

func GetAllDataFromDBService() ([]models.User, error) {
	var users []models.User

	client, err := config.GetMongoClient()
	if err != nil {
		fmt.Println("Error while getting mongo client")
		return nil, err
	}
	db := constant.TEST_DB
	coll := constant.TEST_USERS_COLLECTION

	userCollection := client.Database(db).Collection(coll)
	cursor, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		fmt.Println("Error while finding data")
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("Error while decoding data")
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
