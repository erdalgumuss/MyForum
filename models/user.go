package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `json:"password"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New())
}
