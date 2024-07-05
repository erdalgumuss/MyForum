package handlers

import (
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register function for registering new users
func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email or username already exists
	var user models.User
	if config.DB.Where("email = ?", input.Email).First(&user).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already taken"})
		return
	}
	if config.DB.Where("username = ?", input.Username).First(&user).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user = models.User{
		Email:    input.Email,
		Username: input.Username,
		Password: string(hashedPassword),
	}

	config.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// Login function for authenticating users
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if config.DB.Where("email = ?", input.Email).First(&user).RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := models.Session{
		UserID: user.ID,
	}
	config.DB.Create(&session)

	c.SetCookie("session_id", session.ID.String(), 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout function for logging out users
func Logout(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	var session models.Session
	if config.DB.Where("id = ?", sessionID).First(&session).RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	config.DB.Delete(&session)

	c.SetCookie("session_id", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
