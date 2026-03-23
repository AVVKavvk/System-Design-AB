package models

type Author struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Books []Book // One-to-Many relationship
}
