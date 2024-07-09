package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			log.Println("No session token found, redirecting to home")
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		session, err := getSessionByToken(sessionToken)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Invalid session token, no session found for token:", sessionToken)
			} else {
				log.Println("Error querying session:", err)
			}
			// Oturumu sonlandır
			logoutUser(c, sessionToken)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			log.Println("Session expired for token:", sessionToken)
			// Oturumu sonlandır
			logoutUser(c, sessionToken)
			return
		}

		c.Set("userID", session.UserID)
		c.Next()
	}
}

// Otomatik çıkış işlemi yapan fonksiyon
func logoutUser(c *gin.Context, sessionToken string) {
	// Veritabanından veya hafızadan oturum token'ını sil
	err := deleteSessionByToken(sessionToken)
	if err != nil {
		log.Println("Error deleting session:", err)
	}
	// Çerezi temizle
	c.SetCookie("session_token", "", -1, "/", "", false, true)
	// Kullanıcıyı anasayfaya yönlendir
	c.Redirect(http.StatusFound, "/")
	c.Abort()
}

// Veritabanından veya hafızadan oturum token'ını silen fonksiyon
func deleteSessionByToken(token string) error {
	// SQLite veritabanından oturum token'ını sil
	_, err := config.DB.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func getSessionByToken(token string) (*models.Session, error) {
	var session models.Session
	err := config.DB.QueryRow("SELECT id, user_id, created_at, expires_at FROM sessions WHERE token = ?", token).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
	return &session, err
}

func CreateSession(userID int) (string, error) {
	sessionToken := uuid.New().String()
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	updatedAt := createdAt
	expiresAt := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05") // 24 saat sonra

	insertSessionQuery := `
    INSERT INTO sessions (user_id, token, created_at, updated_at, expires_at)
    VALUES (?, ?, ?, ?, ?)
    `
	_, err := config.DB.Exec(insertSessionQuery, userID, sessionToken, createdAt, updatedAt, expiresAt)
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}

	return sessionToken, nil
}

// GetUserIDFromSession retrieves the user ID from the context set by the AuthMiddleware.
// Returns the user ID and a boolean indicating if the user ID was successfully retrieved.
func GetUserIDFromSession(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	// Type assertion with check
	id, ok := userID.(int)
	if !ok {
		log.Println("UserID found in context is not of type uint")
		return 0, false
	}

	return id, true
}

// AuthRequired middleware checks if the user is authenticated
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		var userID int
		err := config.DB.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

// CheckPasswordHash compares a password with its hashed value.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
