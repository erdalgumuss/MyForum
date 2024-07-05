package handlers

import (
	"database/sql"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"github.com/gin-gonic/gin"
)

// ProfileView handles viewing user profile
func ProfileView(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Handle null values for user.Email and user.Username
	if !user.Email.Valid {
		user.Email.String = ""
	}
	if !user.Username.Valid {
		user.Username.String = ""
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": user,
	})
}

// ProfileUpdate handles updating user profile
func ProfileUpdate(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Update user fields only if they are provided in the JSON input
	if input.Email != "" {
		user.Email = sql.NullString{String: input.Email, Valid: true}
	}

	if input.Username != "" {
		user.Username = sql.NullString{String: input.Username, Valid: true}
	}

	// Execute the update statement
	_, err = config.DB.Exec("UPDATE users SET email=?, username=? WHERE id=?", user.Email.String, user.Username.String, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// RegisterProfileRoutes registers routes for user profile operations
func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile", utils.AuthRequired(), ProfileView)
	router.POST("/profile", utils.AuthRequired(), ProfileUpdate)
}
