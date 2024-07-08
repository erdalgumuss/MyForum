package controllers

import (
	"database/sql"
	"net/http"

	"MyForum/config"
	"MyForum/models"
	"MyForum/utils"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

// GetUserProfile retrieves user profile details
func GetUserProfile(c *gin.Context) {
	var user models.User
	err := config.DB.QueryRow("SELECT id, username, email FROM users WHERE id = ?", c.Param("id")).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

// UpdateProfile handles updating user profile.
func UpdateProfile(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user from database
	var user models.User
	err := config.DB.QueryRow("SELECT id, email, username FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Update user fields only if they are provided in the JSON input
	if input.Email != "" {
		user.Email = input.Email
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	// Update user in the database
	stmt, err := config.DB.Prepare("UPDATE users SET email = ?, username = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Username, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// ChangePassword handles changing the user's password.
func ChangePassword(c *gin.Context) {
	var input struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch user from database
	var user models.User
	err := config.DB.QueryRow("SELECT id, password FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user password"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update user password in the database
	stmt, err := config.DB.Prepare("UPDATE users SET password = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(hashedPassword, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

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

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
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
		user.Email = input.Email
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	// Execute the update statement
	_, err = config.DB.Exec("UPDATE users SET email=?, username=? WHERE id=?", user.Email, user.Username, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// RegisterProfileRoutes registers routes for user profile operations
func RegisterProfileRoutes(router *gin.Engine) {
	router.GET("/profile/:id", utils.AuthRequired(), GetUserProfile)
	router.PUT("/profile", utils.AuthRequired(), UpdateProfile)
	router.POST("/profile/change-password", utils.AuthRequired(), ChangePassword)
	router.GET("/profile", utils.AuthRequired(), ProfileView)
}
