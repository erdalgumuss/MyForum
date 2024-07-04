package models

type Topic struct {
	ID       uint   `gorm:"primary_key"`
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	Category string `gorm:"not null"`
	UserID   uint
	User     User
}


