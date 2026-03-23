package models

type Book struct {
	ID       uint `gorm:"primaryKey"`
	Title    string
	AuthorID uint
}
