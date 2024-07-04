// controllers/auth_controller.go

package controllers

import (
	"net/http"

	"MyForum/handlers"

	"github.com/jinzhu/gorm"
)

// RegisterUser kullanıcı kayıt controller'ı
func RegisterUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	handlers.RegisterHandler(w, r, db)
}
func LoginUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    handlers.LoginHandler(w, r, db)
}