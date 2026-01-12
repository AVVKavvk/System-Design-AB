package models

type User struct {
	ID   string `json:"id"`
	NAME string `json:"name"`
	AGE  int    `json:"age"`
	YEAR string `json:"year"`
}
