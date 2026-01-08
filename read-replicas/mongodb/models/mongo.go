package models

type User struct {
	NAME  string `json:"name" bson:"name"`
	EMAIL string `json:"email" bson:"email"`
	AGE   int    `json:"age" bson:"age"`
}
