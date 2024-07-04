package main

import (
	"fmt"
	"log"
	"net/http"

	"MyForum/controllers"
	"MyForum/handlers"
	"MyForum/models"
	"MyForum/routes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func main() {
	// Veritabanı bağlantısını açma
	var err error
	db, err = gorm.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı başarısız: %v", err)
	}
	defer db.Close()

	// Tabloları oluşturma (Migration)
	db.AutoMigrate(&models.User{}, &models.Topic{}, &models.Comment{})

	fmt.Println("Veritabanı tabloları oluşturuldu.")
	// API endpoint'lerini tanımla
	http.HandleFunc("/api/users", handlers.UsersHandler)
	http.HandleFunc("/api/users/{id}", handlers.UserHandler)
	http.HandleFunc("/api/users/register", handlers.RegisterHandler)
	http.HandleFunc("/api/users/login", handlers.LoginHandler(db))
	http.HandleFunc("/api/topics", handlers.TopicsHandler)
	http.HandleFunc("/api/topics/{id}", handlers.TopicHandler)
	http.HandleFunc("/api/topics/{topic_id}/comments", handlers.CommentsHandler)
	http.HandleFunc("/api/topics/{topic_id}/comments/{comment_id}", handlers.CommentHandler)
	http.HandleFunc("/api/users/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegisterUser(w, r, db)
	})
	http.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginUser(w, r, db)
	})
	// Forum route'larını ayarla
	routes.SetupForumRoutes()

	// Web sunucusunu başlat
	fmt.Println("API server başlatıldı: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Kullanıcı kayıt endpoint'i
