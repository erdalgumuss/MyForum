package handlers

import (
	"fmt"
	"net/http"

	"MyForum/models"

	"gorm.io/gorm"
)

var db *gorm.DB

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// POST isteğini al
	if r.Method != "POST" {
		http.Error(w, "Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
		return
	}

	// Form verilerini al
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// Yeni bir kullanıcı nesnesi oluştur
	newUser := models.User{
		Username: username,
		Password: password,
		Email:    email,
	}

	// Kullanıcıyı veritabanına kaydet
	if err := db.Create(&newUser).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Başarılı yanıtı gönder
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Kullanıcı başarıyla kaydedildi")
}

// LoginHandler kullanıcı giriş endpoint'i
func LoginHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
    // POST isteğini al
    if r.Method != "POST" {
        http.Error(w, "Sadece POST istekleri kabul edilir", http.StatusMethodNotAllowed)
        return
    }

    // Form verilerini al
    username := r.FormValue("username")
    password := r.FormValue("password")

    // Kullanıcıyı veritabanında ara
    var user models.User
    if err := db.Where("username = ?", username).First(&user).Error; err != nil {
        http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
        return
    }

    // Parolayı doğrula
    if user.Password != password {
        http.Error(w, "Geçersiz parola", http.StatusUnauthorized)
        return
    }

    // Başarılı yanıtı gönder
    fmt.Fprintf(w, "Kullanıcı girişi başarılı: %s", username)
}
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcıları işle
	// Örneğin, tüm kullanıcıları listeleme
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir kullanıcıyı işle
	// Örneğin, kullanıcı bilgilerini getirme veya silme
}

func TopicsHandler(w http.ResponseWriter, r *http.Request) {
	// Konuları işle
	// Örneğin, tüm konuları listeleme veya yeni konu oluşturma
}

func TopicHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir konuyu işle
	// Örneğin, konu bilgilerini getirme veya güncelleme
}

func CommentsHandler(w http.ResponseWriter, r *http.Request) {
	// Yorumları işle
	// Örneğin, belirli bir konuya yeni yorum ekleme veya tüm yorumları listeleme
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	// Belirli bir yorumu işle
	// Örneğin, yorum bilgilerini güncelleme veya silme
}

// main.go
