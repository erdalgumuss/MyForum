package utils

import (
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

func CreateSession(userID uint) (string, error) {
	sessionID := uuid.New().String()
	_, err := config.DB.Exec("INSERT INTO sessions (id, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)", sessionID, userID, time.Now(), time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}
	return sessionID, nil
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
