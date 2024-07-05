package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Posts []Post `gorm:"many2many:post_categories;"`
}
