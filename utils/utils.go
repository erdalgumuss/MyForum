package utils

import (
	"net/http"
	"strconv"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		var session models.Session
		// Örneğin, session tablosunu sorgulamak için
		stmt, err := config.DB.Prepare("SELECT id, user_id, expires_at FROM sessions WHERE id = ?")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}
		defer stmt.Close()

		// Sorguyu çalıştır
		err = stmt.QueryRow(sessionID).Scan(&session.ID, &session.UserID, &session.ExpiresAt)
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

		c.Set("user_id", session.UserID)
		c.Next()
	}
}

func StringToUint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0
	}
	return uint(num)
}
