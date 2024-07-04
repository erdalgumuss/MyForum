package models

type Comment struct {
	ID      uint   `gorm:"primary_key"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	TopicID uint
	Topic   Topic
}