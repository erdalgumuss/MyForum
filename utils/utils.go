package utils

import (
	"MyForum/config"
	"MyForum/models"
	"net/http"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		if err := config.DB.Where("id = ?", uuid.MustParse(sessionID)).First(&session).Error; err != nil {
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
