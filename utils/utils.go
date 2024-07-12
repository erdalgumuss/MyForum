package utils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

var googleOauthConfig *oauth2.Config


func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	googleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

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

func GoogleCallback(c *gin.Context) {
	credential := struct {
		Credential string `json:"credential"`
	}{}
	if err := c.BindJSON(&credential); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := idtoken.Validate(context.Background(), credential.Credential, googleOauthConfig.ClientID)
	if err != nil {
		log.Println("Invalid token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := token.Subject

	var user models.User
	err = config.DB.QueryRow("SELECT id FROM users WHERE google_id = ?", userID).Scan(&user.ID)
	if err != nil {
		// User not found, create a new user
		_, err = config.DB.Exec("INSERT INTO users (google_id, email, username) VALUES (?, ?, ?)", userID, token.Claims["email"], token.Claims["name"])
		if err != nil {
			log.Println("Failed to create user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		err = config.DB.QueryRow("SELECT id FROM users WHERE google_id = ?", userID).Scan(&user.ID)
		if err != nil {
			log.Println("Failed to retrieve new user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	// Set the session token
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = config.DB.Exec("INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)", sessionToken, user.ID, time.Now(), expiresAt)
	if err != nil {
		log.Println("Failed to create session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.SetCookie("session_token", sessionToken, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
