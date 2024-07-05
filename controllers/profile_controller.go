package controllers

import (
	"database/sql"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	// Update username if not empty
	if input.Username.Valid {
		user.Username = input.Username.String
	}

	// Update email if not empty
	if input.Email.Valid {
		user.Email = input.Email.String
	}

	// Prepare and execute update statement
	stmt, err := config.DB.Prepare("UPDATE users SET username = ?, email = ? WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare update statement"})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetUserProfile(c *gin.Context) {
	// Implement logic to fetch user profile
	// This function should retrieve and return user profile information
	c.JSON(http.StatusOK, gin.H{
		"message": "GetUserProfile function placeholder",
	})
}

func ChangePassword(c *gin.Context) {
	// Implement logic to change user password
	// This function should handle password change functionality
	c.JSON(http.StatusOK, gin.H{
		"message": "ChangePassword function placeholder",
	})
}
