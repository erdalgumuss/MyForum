package utils

import (
	"fmt"
	"net/http"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session models.Session
		err = config.DB.QueryRow("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?", sessionToken).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		c.Set("userID", session.UserID)
		c.Next()
	}
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

func GetUserIDFromSession(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// AuthRequired middleware checks if the user is authenticated

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session models.Session
		err = config.DB.QueryRow("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?", sessionToken).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		c.Set("userID", session.UserID)
		c.Next()
	}
}
