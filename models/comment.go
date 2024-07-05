package models

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	PostID  uint
	Post    Post
	UserID  uint
	User    User
	Likes   int
	Dislikes int
}
