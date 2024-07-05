package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	UserID    uint
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (session *Session) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	scope.SetColumn("ID", uuid)
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("ExpiresAt", time.Now().Add(24*time.Hour))
	return nil
}
