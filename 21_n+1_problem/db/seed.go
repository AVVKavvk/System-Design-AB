package db

import (
	"log"

	"github.com/AVVKavvk/n+1-problem/models"
	"gorm.io/gorm"
)

func Seed(database *gorm.DB) {
	// Clear existing data to start fresh
	database.Exec("TRUNCATE TABLE books, authors RESTART IDENTITY")

	// Define sample data
	authors := []models.Author{
		{Name: "J.K. Rowling", Books: []models.Book{{Title: "Stone"}, {Title: "Chamber"}}},
		{Name: "George R.R. Martin", Books: []models.Book{{Title: "Game of Thrones"}, {Title: "Clash of Kings"}}},
		{Name: "J.R.R. Tolkien", Books: []models.Book{{Title: "The Hobbit"}, {Title: "Fellowship"}}},
		{Name: "Stephen King", Books: []models.Book{{Title: "The Shining"}, {Title: "It"}}},
		{Name: "Agatha Christie", Books: []models.Book{{Title: "Murder on Orient Express"}}},
	}

	// Insert into Database
	for _, author := range authors {
		if err := database.Create(&author).Error; err != nil {
			log.Printf("Could not seed author %s: %v", author.Name, err)
		}
	}

	log.Println("Database seeded successfully!")
}
