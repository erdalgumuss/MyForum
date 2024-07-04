package main

import (
	"fmt"
	"log"
	"net/http"

	"MyForum/models"

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
	http.HandleFunc("/api/users", usersHandler)
	http.HandleFunc("/api/users/{id}", userHandler)
	http.HandleFunc("/api/topics", topicsHandler)
	http.HandleFunc("/api/topics/{id}", topicHandler)
	http.HandleFunc("/api/topics/{topic_id}/comments", commentsHandler)
	http.HandleFunc("/api/topics/{topic_id}/comments/{comment_id}", commentHandler)

	// Web sunucusunu başlat
	fmt.Println("API server başlatıldı: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcıları işle
	// Örneğin, tüm kullanıcıları listeleme
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir kullanıcıyı işle
	// Örneğin, kullanıcı bilgilerini getirme veya silme
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	// Konuları işle
	// Örneğin, tüm konuları listeleme veya yeni konu oluşturma
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir konuyu işle
	// Örneğin, konu bilgilerini getirme veya güncelleme
}

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	// Yorumları işle
	// Örneğin, belirli bir konuya yeni yorum ekleme veya tüm yorumları listeleme
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir yorumu işle
	// Örneğin, yorum bilgilerini güncelleme veya silme
}
