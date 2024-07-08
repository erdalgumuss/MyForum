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

// AuthMiddleware is a middleware function for checking user authentication via session token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve session token from the cookies
		sessionToken, err := c.Cookie("session_token")
		if err != nil {
			log.Println("No session token found, redirecting to home")
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		log.Println("Session token retrieved from cookie:", sessionToken)

		// Query the session from the database using the session token
		var session models.Session
		err = config.DB.QueryRow("SELECT id, user_id, created_at, expires_at FROM sessions WHERE id = ?", sessionToken).Scan(&session.ID, &session.UserID, &session.CreatedAt, &session.ExpiresAt)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Invalid session token, no session found for token:", sessionToken)
			} else {
				log.Println("Error querying session:", err)
			}
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		// Check if the session is expired
		if session.ExpiresAt.Before(time.Now()) {
			log.Println("Session expired for token:", sessionToken)
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		// Set the userID in the context for further handlers
		log.Println("User authenticated with ID:", session.UserID)
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
