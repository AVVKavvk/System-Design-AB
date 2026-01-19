package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}
