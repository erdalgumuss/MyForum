package models

import "github.com/jinzhu/gorm"

type Post struct {
	gorm.Model
	Title      string `gorm:"not null"`
	Content    string `gorm:"not null"`
	UserID     uint
	User       User
	Comments   []Comment
	Categories []Category `gorm:"many2many:post_categories;"`
	Likes      int
	Dislikes   int
}
