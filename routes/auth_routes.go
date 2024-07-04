package routes

import (
	"MyForum/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func SetupAuthRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(w, r, db)
	}).Methods("POST")

	// Diğer auth endpoint'lerini burada tanımlayabilirsiniz
}
