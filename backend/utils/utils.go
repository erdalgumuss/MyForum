package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			// End session
			logoutUser(c, sessionToken)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			log.Println("Session expired for token:", sessionToken)
			// End session
			logoutUser(c, sessionToken)
			return
		}

		var username string
		err = config.DB.QueryRow("SELECT username FROM users WHERE id = ?", session.UserID).Scan(&username)
		if err != nil {
			log.Println("Failed to fetch username:", err)
			logoutUser(c, sessionToken)
			return
		}

		c.Set("userID", session.UserID)
		c.Set("username", username)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		Role := c.GetString("role")
		if Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if role == nil || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ModeratorOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		role := session.Get("role")
		if role == nil || role != "moderator" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}
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
	sessionToken := uuid.NewString()
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
		log.Println("UserID found in context is not of type int")
		return 0, false
	}
	return id, true
}

var forbiddenWords = []string{
	"küfür1", "küfür2", "illegal1", "illegal2", // Eklemek istediğiniz diğer kelimeler
}

func ContainsForbiddenContent(content string) bool {
	content = strings.ToLower(content)
	for _, word := range forbiddenWords {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}
